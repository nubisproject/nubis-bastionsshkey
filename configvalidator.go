package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
)

type IAMGroupMapping struct {
	LDAPGroup  string `yaml:"LDAPGroup"`
	IAMPath    string `yaml:"IAMPath"`
	ConsulPath string `yaml:"ConsulPath"`
}

type Configuration struct {
	LdapServer struct {
		LDAPHost         string            `yaml:"LDAPHost"`
		LDAPServer       string            `yaml:"LDAPServer"`
		LDAPBindPassword string            `yaml:"LDAPBindPassword"`
		LDAPBindUser     string            `yaml:"LDAPBindUser"`
		LDAPInsecure     bool              `yaml:"LDAPInsecure"`
		StartTLS         bool              `yaml:"StartTLS"`
		IAMGroupMapping  []IAMGroupMapping `yaml:"IAMGroupMapping"`
		TLSCrt           string            `yaml:"TLSCrt"`
		TLSKey           string            `yaml:"TLSKey"`
		CACrt            string            `yaml:"CACrt"`
	} `yaml:"LdapServer"`
	Consul struct {
		Server    string `yaml:"Server"`
		Namespace string `yaml:"Namespace"`
	} `yaml:"Consul"`
	AWS struct {
		AccessKey         string   `yaml:"AccessKey"`
		SecretKey         string   `yaml:"SecretKey"`
		AWSIgnoreUserList []string `yaml:"AWSIgnoreUserList"`
		AWSIgnorePathList []string `yaml:"AWSIgnorePathList"`
		SMTPUsername      string   `yaml:"SMTPUsername"`
		SMTPPassword      string   `yaml:"SMTPPassword"`
		SMTPHostname      string   `yaml:"SMTPHostname"`
		SMTPPort          string   `yaml:"SMTPPort"`
		SMTPFromAddress   string   `yaml:"SMTPFromAddress"`
	} `yaml:"AWS"`
}

func getConfig(c ConfigOptions) (Configuration, error) {
	var yamlData []byte
	var err error
	if c.UseDynamo == false {
		filename, _ := filepath.Abs(c.ConfigFilePath)
		yamlData, err = ioutil.ReadFile(filename)
	} else {
		unicreds := c.UnicredsPath
		var out bytes.Buffer
		var stdErr bytes.Buffer
		cmdArgs := []string{"--region", c.Region, "get", c.Key, "-E", fmt.Sprintf("environment:%s", c.Environment), "-E", fmt.Sprintf("service:%s", c.Service)}
		cmd := exec.Command(unicreds, cmdArgs...)
		cmd.Stdout = &out
		cmd.Stderr = &stdErr
		err := cmd.Run()
		if err != nil {
			log.Print(err)
		}
		log.Print(stdErr.String())
		log.Print(out.String())
		//cmdString := fmt.Sprintf("%s --region %s get %s -E environment:%s -E service:%s", unicreds, c.Region, c.Key, c.Environment, c.Service)
		yamlData = []byte(out.String())
		fmt.Println(yamlData)
		// Here we connect to dynamoDB and return the yamlData
	}
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
