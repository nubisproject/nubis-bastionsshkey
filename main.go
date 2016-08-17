package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
)

func RemoveDuplicates(xs *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[x] {
			found[x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}

func SortUsers(globalAdmins []LDAPUserObject, sudoUsers []LDAPUserObject) []string {
	returnSlice := []string{}
	for _, g_entry := range globalAdmins {
		returnSlice = append(returnSlice, g_entry.Uid)
	}
	RemoveDuplicates(&returnSlice)
	sort.Strings(returnSlice)
	return returnSlice

}
func IgnoreUser(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func main() {
	var configFilePath string
	var execType string
	var noop bool
	flag.StringVar(&configFilePath, "c", "config.yml", "Configuration file to use")
	flag.StringVar(&execType, "execType", "consul", "consul|IAM\nUse consul to sync LDAP to consul, use IAM to sync IAM users from LDAP")
	flag.BoolVar(&noop, "noop", false, "noop - providing noop makes functionality displayed without taking any action")
	flag.Parse()
	configuration, err := getConfig(configFilePath)
	if err != nil {
		log.Fatal("Unable to read configuration")
	}
	configValid, configError := validateConfig(configuration)
	if configValid == false {
		fmt.Println(configError)
		os.Exit(2)
	}
	c := GetConsulClient(configuration)
	globalAdmins, sudoUsers := getGroupMembers(configuration)
	usersSet := SortUsers(globalAdmins, sudoUsers)
	if execType == "consul" {
		for _, entry := range globalAdmins {
			if noop == false {
				c.Put(entry, configuration, "global-admins")
			} else {
				fmt.Println(entry.Uid)
			}
		}
		for _, entry := range sudoUsers {
			if noop == false {
				c.Put(entry, configuration, "sudo-users")
			} else {
				fmt.Println(entry.Uid)
			}
		}
	} else if execType == "IAM" {

		for _, entry := range usersSet {
			log.Printf(entry)
		}
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
						log.Printf("Adding: %s to iamUsers", user)
					}
				}
			}
		}
	}
}
