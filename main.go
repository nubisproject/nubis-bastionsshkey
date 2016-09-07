package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	configFilePath   string
	execType         string
	testDestEmail    string
	noop             bool
	useDynamo        bool
	region           string
	key              string
	environment      string
	accountName      string
	service          string
	unicredsPath     string
	consulPort       string
	consulDomain     string
	useLambda        bool
	testUserName     string
	userCreationPath string
)

func parseFlags() {
	flag.StringVar(&testUserName, "testUserName", "", "Test UserName for creating a user. Will be removed, for debugging only.")
	flag.StringVar(&userCreationPath, "userCreationPath", "", "Test userCreationPath for creating a user. Will be removed, for debugging only.")
	flag.StringVar(&configFilePath, "c", "", "Configuration file to use")
	flag.StringVar(&execType, "execType", "consul", "consul|IAM\nUse consul to sync LDAP to consul, use IAM to sync IAM users from LDAP")
	flag.StringVar(&testDestEmail, "testDestEmail", "", "Email Address for testing email")

	// dynamoDB flags
	flag.BoolVar(&useDynamo, "useDynamo", false, "Bool to use dynamodb for config file")
	flag.StringVar(&region, "region", "us-west-2", "dynamoDB Region")
	flag.StringVar(&key, "key", "", "dynamoDB key")
	flag.StringVar(&environment, "environment", "", "dynamoDB Region")
	flag.StringVar(&service, "service", "", "dynamoDB Region")
	flag.StringVar(&accountName, "accountName", "", "accountName")
	flag.StringVar(&consulPort, "consulPort", "8500", "Consul port to connect to")
	flag.StringVar(&consulDomain, "consulDomain", "localhost", "Domain of the consul server")
	// end dynamoDB flags
	flag.BoolVar(&noop, "noop", false, "noop - providing noop makes functionality displayed without taking any action")
	flag.BoolVar(&useLambda, "lambda", false, "Use lambda flag")
	flag.Parse()
}

func main() {
	// @TODO: Temporary variable to be removed after confident user creation is being handled correctly
	testUserName := ""
	userCreationPath := ""
	parseFlags()
	if configFilePath != "" && useDynamo != false {
		log.Fatal("Incorrect flags. dynamoDBPath and configFilePath cannot both be provided.")
	}

	d := ConfigOptions{}
	log.Println(os.Args[1:])
	if useDynamo == true {
		if region == "" {
			log.Fatal("-region is required when using dynamoDBPath")
		}
		if key == "" {
			log.Fatal("-key is required when using dynamoDBPath")
		}
		if accountName == "" {
			log.Fatal("-accountName is required when using dynamoDBPath")
		}
		if environment == "" {
			log.Fatal("-environment is required when using dynamoDBPath")
		}
		if service == "" {
			log.Fatal("-service is required when using dynamoDBPath")
		}
		if unicredsPath == "" {
			unicredsPath = "./unicreds"
		}
		d.Region = region
		d.Environment = environment
		d.Service = service
		d.AccountName = accountName
		d.Key = key
		d.UseDynamo = true
		d.UnicredsPath = "./unicreds"
		d.ConsulDomain = consulDomain

		if useLambda {
			// FIXME: If you are using dynamodb and lambda
			// it means you need to export proxy info
			// since unicreds need to be able to do this
			http_proxy := fmt.Sprintf("http://proxy.%s.%s.%s.%s:3128/", d.Environment, d.Region, d.AccountName, d.ConsulDomain)
			https_proxy := fmt.Sprintf("https://proxy.%s.%s.%s.%s:3128/", d.Environment, d.Region, d.AccountName, d.ConsulDomain)
			os.Setenv("HTTP_PROXY", http_proxy)
			os.Setenv("HTTPS_PROXY", https_proxy)
		}

	}
	if useDynamo == false && configFilePath == "" {
		d.ConfigFilePath = "config.yml"
		d.UseDynamo = false
	}
	configuration, err := getConfig(d)
	if err != nil {
		log.Fatal("Unable to read configuration")
	}
	configValid, configError := validateConfig(configuration)
	if configValid == false {
		fmt.Println(configError)
		os.Exit(2)
	}
	if useDynamo == true && useLambda == true {
		configuration.AWS.Region = d.Region
		d.ConsulDomain = consulDomain
		d.ConsulPort = consulPort
		configuration.Consul.Server = d.DeriveConsulServer()
		configuration.Consul.Token = d.getConsulACLToken()
	}

	c := GetConsulClient(configuration)
	var usersSet []string
	var allLDAPGroupUserObjects []LDAPUserObject
	var allEntries []ConsulEntries
	for _, x := range configuration.LdapServer.IAMGroupMapping {
		tmpGroupMembers := getGroupMembers(configuration, x)
		for _, user := range tmpGroupMembers {
			usersSet = append(usersSet, user.Uid)
			allLDAPGroupUserObjects = append(allLDAPGroupUserObjects, user)
		}
		tmp := ConsulEntries{tmpGroupMembers, x}
		allEntries = append(allEntries, tmp)
	}
	usersSet = SortUsers(usersSet)

	if execType == "consul" {
		for _, g_entry := range allEntries {
			for _, entry := range g_entry.Users {
				if noop == false {
					c.Put(entry, configuration, g_entry.Group.ConsulPath)
				} else {
					fmt.Println(entry.Uid)
				}
			}
			SyncLDAPToConsul(g_entry.Group.ConsulPath, g_entry.Users, noop, c, configuration)
		}
	} else if execType == "IAM" {
		IAMUsers, IAMUsersErr := GetAllIAMUsers(configuration)
		if IAMUsersErr != nil {
			log.Println(IAMUsersErr)
		}
		if IAMUsersErr == nil {
			iamUsers := []string{}
			for _, user := range IAMUsers.Users {
				username := *user.UserName
				path := *user.Path
				ignoreUser := IgnoreUser(configuration.AWS.AWSIgnoreUserList, username)
				ignorePath := IgnoreUser(configuration.AWS.AWSIgnorePathList, path)
				if ignoreUser == false && ignorePath == false {
					iamUsers = append(iamUsers, username)
					if noop {
						log.Printf("Acting upon user: %s", username)
					}
				} else {
					if noop {
						log.Printf("Ignoring user: %s", username)
					}
				}
			}
			for _, user := range iamUsers {
				if IgnoreUser(usersSet, user) == false {
					if noop {
						log.Printf("User: %s doesn't exist in usersSet", user)
					} else {
						log.Printf("Removing: %s from iamUsers", user)

					}
				}
			}
			for _, user := range usersSet {
				if IgnoreUser(iamUsers, user) == false {
					if noop {
						log.Printf("User: %s doesn't exist in iamUsers", user)
					} else {
						path := userCreationPath
						if user == testUserName {
							userRet, err := CreateIAMUser(configuration, user, path)
							if err != nil {
								log.Fatal(err)
							}
							userLDAPObj, found := GetLDAPUserObjectFromGroup(user, allLDAPGroupUserObjects)
							if found == false {
								fmt.Println("Here we need to log/track that people don't have a PGPPublicKey in LDAP")
							}
							emailBody := []byte(fmt.Sprintf("AccessKey: %s\nSecretKey: %s", userRet.AccessKey, userRet.SecretKey))
							testEncrypted, encryptErr := EncryptMailBody(emailBody, userLDAPObj.PGPPublicKey, testDestEmail)
							if encryptErr != nil {
								log.Fatal(encryptErr)
							}
							SendWelcomeMail(configuration, testDestEmail, testEncrypted)

						}
						log.Printf("Adding: %s to iamUsers", user)
					}
				}
			}
		}
	}
}
