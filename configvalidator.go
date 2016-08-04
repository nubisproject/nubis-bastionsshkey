package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type Configuration struct {
	Server struct {
		LDAPHost         string   `yaml:"LDAPHost"`
		LDAPServer       string   `yaml:"LDAPServer"`
		LDAPBindPassword string   `yaml:"LDAPBindPassword"`
		LDAPBindUser     string   `yaml:"LDAPBindUser"`
		LDAPInsecure     bool     `yaml:"LDAPInsecure"`
		StartTLS         bool     `yaml:"StartTLS"`
		SearchGroups     []string `yaml:"SearchGroups"`
		TLSCrt           string   `yaml:"TLSCrt"`
		TLSKey           string   `yaml:"TLSKey"`
		CACrt            string   `yaml:"CACrt"`
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
	if config.Server.LDAPHost == "" {
		returnString = "LDAPHost required"
		returnValid = false
	} else if config.Server.LDAPBindUser == "" {
		returnString = "LDAPBindUser required"
		returnValid = false
	} else if config.Server.LDAPBindPassword == "" {
		returnString = "LDAPBindPassword required"
		returnValid = false
	} else if len(config.Server.SearchGroups) == 0 {
		returnString = "At minimum 1 SearchGroups required"
		returnValid = false
	}

	return returnValid, returnString

}
