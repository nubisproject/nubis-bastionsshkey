package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "c", "config.yml", "Configuration file to use")
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
	for _, entry := range globalAdmins {
		c.Put(entry, configuration, "global-admins")
		fmt.Println(entry.Uid)
	}
	for _, entry := range sudoUsers {
		c.Put(entry, configuration, "sudo-users")
		fmt.Println(entry.Uid)
	}
}
