package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	configFilePath   string
	execType         string
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
	userCreationPath string
	showVersion      bool
	userPathList     UserPathList
	Version          string
)

func parseFlags() {
	flag.StringVar(&userCreationPath, "userCreationPath", "", "Test userCreationPath for creating a user. Will be removed, for debugging only.")
	flag.StringVar(&configFilePath, "c", "", "Configuration file to use")
	flag.StringVar(&execType, "execType", "consul", "consul|IAM\nUse consul to sync LDAP to consul, use IAM to sync IAM users from LDAP")

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
	flag.BoolVar(&showVersion, "version", false, "Show version and exit")
	flag.Parse()
}

func main() {
	parseFlags()

	if showVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	if configFilePath != "" && useDynamo != false {
		log.Fatal("Incorrect flags. dynamoDBPath and configFilePath cannot both be provided.")
	}

	d := ConfigOptions{}
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

	log.Printf("Running version %s of nubis-bastionsshkey", Version)
	c := GetConsulClient(configuration)
	var usersSet []string
	var allLDAPGroupUserObjects []LDAPUserObject
	var allEntries []ConsulEntries
	for _, x := range configuration.LdapServer.IAMGroupMapping {
		tmpGroupMembers := getGroupMembers(configuration, x)
		for _, user := range tmpGroupMembers {
			usersSet = append(usersSet, user.Uid)
			allLDAPGroupUserObjects = append(allLDAPGroupUserObjects, user)
			tmp := UserPath{user.Uid, x.IAMPath}
			userPathList.add(tmp)
		}
		tmp := ConsulEntries{tmpGroupMembers, x}
		allEntries = append(allEntries, tmp)
	}
	usersSet = SortUsers(usersSet)

	if execType == "consul" {
		for _, g_entry := range allEntries {
			for _, entry := range g_entry.Users {
				if noop == false {
					c.Put(entry, configuration, g_entry.Group.LDAPGroup)
				} else {
					fmt.Println(entry.Uid)
				}
			}
			SyncLDAPToConsul(g_entry.Group.LDAPGroup, g_entry.Users, noop, c, configuration)
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
				}
			}
			iamUsersDiff := IAMUsersDiff{
				IAMUsers,
				usersSet,
				configuration.AWS.AWSIgnorePathList,
				configuration.AWS.AWSIgnoreUserList,
			}
			usersToAdd := iamUsersDiff.getUsersToAdd()
			for _, user := range usersToAdd {
				if IgnoreUser(iamUsers, user) == false {
					if noop == true {
						log.Printf("NOOP: Adding: %s to iamUsers", user)
					} else {
						path := userPathList.getPathByUsername(user)
						log.Printf("Creating user: %s at path: %s", user, path)
						userRet, err := CreateIAMUser(configuration, user, path)
						// Reason this needs to be here is because sometimes a user gets created
						// by AWS but it doesn't show up immediately
						time.Sleep(5 * time.Second)
						if err != nil {
							log.Fatal(err)
						}
						ApplyRoles(configuration, user, path)
						userLDAPObj, found := GetLDAPUserObjectFromGroup(user, allLDAPGroupUserObjects)
						if found == false || string(userLDAPObj.PGPPublicKey) == "" {
							fmt.Println("Here we need to log/track that people don't have a PGPPublicKey in LDAP")
							continue
						}
						userArn, _ := GetUserArn(configuration, user)
						roleArn, _ := GetRoleArn(configuration, user)
						emailBody := []byte(fmt.Sprintf("AccessKey: %s\nSecretKey: %s\nUserArn: %s\nRoleArn: %s\n", userRet.AccessKey, userRet.SecretKey, userArn, roleArn))
						testEncrypted, encryptErr := EncryptMailBody(emailBody, userLDAPObj.PGPPublicKey, userLDAPObj.Mail)
						if encryptErr != nil {
							log.Print("Unable to encrypt message to: ", userLDAPObj.Mail, " with error: ", encryptErr)
						}
						SendWelcomeMail(configuration, userLDAPObj.Mail, testEncrypted)
						log.Printf("Adding: %s to iamUsers", user)
					}
				}
			}
			usersToRemove := iamUsersDiff.getUsersToRemove()
			for _, user := range usersToRemove {
				if noop == true {
					log.Printf("NOOP: Removing: %s from iamUsers", user)
				} else {
					DetachGroup(configuration, user)
					_, deletedErr := DeleteIAMUser(configuration, user)
					if deletedErr == nil {
						log.Printf("Removing: %s from iamUsers", user)
						DeleteRoles(configuration, user)
					}
				}
			}
		}
	}
}
