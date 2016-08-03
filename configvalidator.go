package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type Configuration struct {
	Server struct {
		LDAPBINDPASSWORD string   `yaml:"LDAPBINDPASSWORD"`
		LDAPBINDUSER     string   `yaml:"LDAPBINDUSER"`
		LDAPHOST         string   `yaml:"LDAPHOST"`
		SEARCHGROUPS     []string `yaml:"SEARCHGROUPS"`
	} `yaml:"Server"`
}

func getConfig(configPath string) (Configuration, error) {
	filename, _ := filepath.Abs("./config.yml")
	yamlData, err := ioutil.ReadFile(filename)
	configuration := Configuration{}
	yamlErr := yaml.Unmarshal([]byte(yamlData), &configuration)
	if yamlErr != nil {
		fmt.Println("error:", yamlErr)
	}
	return configuration, err

}

func validateConfig(config Configuration) (bool, string) {
	returnValid := true
	returnString := ""
	if config.Server.LDAPHOST == "" {
		returnString = "LDAPHOST required"
		returnValid = false
	} else if config.Server.LDAPBINDUSER == "" {
		returnString = "LDAPBINDUSER required"
		returnValid = false
	} else if config.Server.LDAPBINDPASSWORD == "" {
		returnString = "LDAPBINDPASSWORD required"
		returnValid = false
	} else if len(config.Server.SEARCHGROUPS) == 0 {
		returnString = "At minimum 1 SEARCHGROUPS required"
		returnValid = false
	}

	return returnValid, returnString

}
