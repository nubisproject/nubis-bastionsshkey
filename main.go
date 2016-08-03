package main

import (
	"fmt"
	"os"
)

func main() {
	filePath := "config.yaml"
	configuration, err := getConfig(filePath)
	if err != nil {
		fmt.Println("Unable to read configuration")
	}
	configValid, configError := validateConfig(configuration)
	if configValid == false {
		fmt.Println(configError)
		os.Exit(2)
	}
	fmt.Println(configuration.Server.LDAPHOST)
}
