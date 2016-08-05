package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type Configuration struct {
	LdapServer struct {
		LDAPHost         string   `yaml:"LDAPHost"`
		LDAPServer       string   `yaml:"LDAPServer"`
		LDAPBindPassword string   `yaml:"LDAPBindPassword"`
		LDAPBindUser     string   `yaml:"LDAPBindUser"`
		LDAPInsecure     bool     `yaml:"LDAPInsecure"`
		StartTLS         bool     `yaml:"StartTLS"`
		GlobalAdmins     []string `yaml:"GlobalAdmins"`
		SudoUsers        []string `yaml:"SudoUsers"`
		TLSCrt           string   `yaml:"TLSCrt"`
		TLSKey           string   `yaml:"TLSKey"`
		CACrt            string   `yaml:"CACrt"`
	} `yaml:"LdapServer"`
	Consul struct {
		Server                string `yaml:"Server"`
		Namespace             string `yaml:"Namespace"`
		Token                 string `yaml:"Token"`
		SSHPublicKeyDelimeter string `yaml:"SSHPublicKeyDelimeter"`
	} `yaml:"Consul"`
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
	if config.LdapServer.LDAPHost == "" {
		returnString = "LDAPHost required"
		returnValid = false
	} else if config.LdapServer.LDAPBindUser == "" {
		returnString = "LDAPBindUser required"
		returnValid = false
	} else if config.LdapServer.LDAPBindPassword == "" {
		returnString = "LDAPBindPassword required"
		returnValid = false
	}

	return returnValid, returnString

}
