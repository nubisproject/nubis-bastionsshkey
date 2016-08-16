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

func main() {
	var configFilePath string
	var execType string
	flag.StringVar(&configFilePath, "c", "config.yml", "Configuration file to use")
	flag.StringVar(&execType, "execType", "consul", "consul|IAM\nUse consul to sync LDAP to consul, use IAM to sync IAM users from LDAP")
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
	if execType == "consul" {
		for _, entry := range globalAdmins {
			c.Put(entry, configuration, "global-admins")
			fmt.Println(entry.Uid)
		}
		for _, entry := range sudoUsers {
			c.Put(entry, configuration, "sudo-users")
			fmt.Println(entry.Uid)
		}
	} else if execType == "IAM" {

		usersSet := SortUsers(globalAdmins, sudoUsers)
		for _, entry := range usersSet {
			log.Printf(entry)

		}

	}
}
