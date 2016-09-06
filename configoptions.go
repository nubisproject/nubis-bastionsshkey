package main

import (
	"fmt"
	"gopkg.in/oleiade/reflections.v1"
	"io/ioutil"
	"path/filepath"
)

type ConfigOptions struct {
	Region         string
	Key            string
	Environment    string
	Service        string
	ConfigFilePath string
	AccountName    string
	UseDynamo      bool
	UnicredsPath   string
	ConsulServer   string
	ConsulPort     string
	ConsulDomain   string
	ConsulToken    string
}

func getDefaultConfig() Configuration {
	filePath := "config.yml-dist"
	filename, _ := filepath.Abs(filePath)
	yamlData, _ := ioutil.ReadFile(filename)
	c, _ := ConfigFromYaml(yamlData)
	return c
}

func (c *ConfigOptions) OverrideField(field string, value string) {
	switch value {
	case "true":
		reflections.SetField(c, field, true)
	case "false":
		reflections.SetField(c, field, false)
	default:
		reflections.SetField(c, field, value)
	}
}
func (c *ConfigOptions) DeriveConsulServer() string {
	derivedConsulHostname := fmt.Sprintf(
		"ui.consul.%s.%s.%s.%s", c.Environment, c.Region, c.AccountName, c.ConsulDomain,
	)
	return derivedConsulHostname
}

func (c *ConfigOptions) OverrideConsulServer(server string) {
	c.ConsulServer = server
}

func (c *ConfigOptions) ShouldDeriveConsulServer(args []string) bool {
	if len(args) == 0 {
		return false
	}
	fmt.Print(args)
	return true
}
