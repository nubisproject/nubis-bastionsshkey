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
		t.Fatalf("Unable to override field. Value set to: %t", defaultConfig.UseDynamo)
	} else {
		t.Log(defaultConfig.Region)
	}
}

func TestOverrideBoolFieldFalse(t *testing.T) {
	defaultConfig.UseDynamo = true
	defaultConfig.OverrideField("UseDynamo", "false")
	if defaultConfig.UseDynamo != false {
		t.Fatalf("Unable to override field. Value set to: %t", defaultConfig.UseDynamo)
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

func TestConfigurationDeriveConsulServerPortProvided(t *testing.T) {
	newConsulServer := "ui.consul.stage.us-doesnt-exist.theAccountName.provided.domain.name:9900"
	defaultConfig.Arena = "stage"
	defaultConfig.Region = "us-doesnt-exist"
	defaultConfig.AccountName = "theAccountName"
	defaultConfig.ConsulPort = "9900"
	defaultConfig.ConsulDomain = "provided.domain.name"
	defaultConfig.ConsulServer = defaultConfig.DeriveConsulServer()
	if defaultConfig.ConsulServer != newConsulServer {
		t.Fatal("Unable to override field. Value set to:", defaultConfig.ConsulServer)
	}
}
func TestConfigurationDeriveConsulServerPortDefault(t *testing.T) {
	newConsulServer := "ui.consul.stage.us-doesnt-exist.theAccountName.provided.domain.name:8500"
	defaultConfig.Arena = "stage"
	defaultConfig.Region = "us-doesnt-exist"
	defaultConfig.AccountName = "theAccountName"
	defaultConfig.ConsulDomain = "provided.domain.name"
	defaultConfig.ConsulPort = ""
	defaultConfig.ConsulServer = defaultConfig.DeriveConsulServer()
	if defaultConfig.ConsulServer != newConsulServer {
		t.Fatal("Unable to override field. Value set to:", defaultConfig.ConsulServer)
	}
}
func TestConfigurationShouldDeriveConsulServerEmpty(t *testing.T) {
	args := []string{}
	shouldDerive := defaultConfig.ShouldDeriveConsulServer(args)
	if shouldDerive != false {
		t.Fatalf("Derive incorrectly set to %t", shouldDerive)
	}
}
