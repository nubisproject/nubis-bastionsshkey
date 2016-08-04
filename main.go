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
	ldapUsers := getGroupMembers(configuration)
	for _, entry := range ldapUsers {
		fmt.Println(entry.Uid)
	}
}
