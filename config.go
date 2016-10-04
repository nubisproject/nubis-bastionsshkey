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
		Token     string `yaml:"Token"`
	} `yaml:"Consul"`
	AWS struct {
		Region            string   `yaml:"Region"`
		AccessKey         string   `yaml:"AccessKey,omitempty"`
		SecretKey         string   `yaml:"SecretKey,omitempty"`
		AWSIgnoreUserList []string `yaml:"AWSIgnoreUserList"`
		AWSIgnorePathList []string `yaml:"AWSIgnorePathList"`
		SMTPUsername      string   `yaml:"SMTPUsername"`
		SMTPPassword      string   `yaml:"SMTPPassword"`
		SMTPHostname      string   `yaml:"SMTPHostname"`
		SMTPPort          string   `yaml:"SMTPPort"`
		SMTPFromAddress   string   `yaml:"SMTPFromAddress"`
	} `yaml:"AWS"`
}

func ConfigFromYaml(yamlData []byte) (Configuration, error) {
	configuration := Configuration{}
	yamlErr := yaml.Unmarshal([]byte(yamlData), &configuration)
	return configuration, yamlErr
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
		cmdArgs := []string{
			"--region", c.Region,
			"get", c.Key,
			"-E", fmt.Sprintf("environment:%s", c.Environment),
			"-E", fmt.Sprintf("service:%s", c.Service),
			"-E", fmt.Sprintf("region:%s", c.Region),
		}
		log.Printf("%s --region %s get %s -E environment:%s -E service:%s -E region:%s", c.UnicredsPath, c.Region, c.Key, c.Environment, c.Service, c.Region)
		cmd := exec.Command(unicreds, cmdArgs...)
		cmd.Stdout = &out
		cmd.Stderr = &stdErr
		err := cmd.Run()
		if err != nil {
			log.Print(err)
		}
		yamlData = []byte(out.String())
	}
	configuration, err := ConfigFromYaml(yamlData)
	return configuration, err

}
