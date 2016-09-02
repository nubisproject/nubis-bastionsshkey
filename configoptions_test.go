package main

import (
	"testing"
)

var defaultConfig = ConfigOptions{}

func TestOverrideStringField(t *testing.T) {
	defaultConfig.OverrideField("Region", "Overridden")
	if defaultConfig.Region != "Overridden" {
		t.Fatal("Unable to override field. Value set to:", defaultConfig.Region)
	} else {
		t.Log(defaultConfig.Region)
	}
}

func TestOverrideBoolFieldTrue(t *testing.T) {
	defaultConfig.UseDynamo = false
	defaultConfig.OverrideField("UseDynamo", "true")
	if defaultConfig.UseDynamo != true {
		t.Fatal("Unable to override field. Value set to: %s", defaultConfig.UseDynamo)
	} else {
		t.Log(defaultConfig.Region)
	}
}

func TestOverrideBoolFieldFalse(t *testing.T) {
	defaultConfig.UseDynamo = true
	defaultConfig.OverrideField("UseDynamo", "false")
	if defaultConfig.UseDynamo != false {
		t.Fatal("Unable to override field. Value set to: %s", defaultConfig.UseDynamo)
	} else {
		t.Log(defaultConfig.Region)
	}
}
func TestConfigurationOverrideConsulServer(t *testing.T) {
	newConsulServer := "new.consul.server"
	defaultConfig.OverrideConsulServer(newConsulServer)
	if defaultConfig.ConsulServer != newConsulServer {
		t.Fatal("Unable to override field. Value set to:", defaultConfig.ConsulServer)
	}
}

func TestConfigurationDeriveConsulServer(t *testing.T) {
	environment := "stage"
	region := "us-doesnt-exist"
	accountName := "theAccountName"
	newConsulServer := "ui.consul.stage.us-doesnt-exist.theAccountName.provided.domain.name:9900"
	domain := "provided.domain.name"
	port := 9900
	defaultConfig.DeriveConsulServer()
	if defaultConfig.ConsulServer != newConsulServer {
		t.Fatal("Unable to override field. Value set to:", defaultConfig.ConsulServer)
	}
}
func TestConfigurationShouldDeriveConsulServerEmpty(t *testing.T) {
	args := []string{}
	shouldDerive := defaultConfig.ShouldDeriveConsulServer(args)
	if shouldDerive != false {
		t.Fatal("Derive incorrectly set to %b", shouldDerive)
	}
}
