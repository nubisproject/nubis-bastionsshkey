package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	filePath := "config.yaml"
	configuration, err := getConfig(filePath)
	if err != nil {
		log.Fatal("Unable to read configuration")
	}
	configValid, configError := validateConfig(configuration)
	if configValid == false {
		fmt.Println(configError)
		os.Exit(2)
	}
	globalAdmins, sudoUsers := getGroupMembers(configuration)
	c := GetConsulClient(configuration)
	for _, entry := range globalAdmins {
		c.Put(entry, configuration, "global-admins")
		fmt.Println(entry.Uid)
	}
	for _, entry := range sudoUsers {
		c.Put(entry, configuration, "sudo-users")
		fmt.Println(entry.Uid)
	}
}
