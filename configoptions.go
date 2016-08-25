package main

type ConfigOptions struct {
	Region         string
	Key            string
	Environment    string
	Service        string
	ConfigFilePath string
	UseDynamo      bool
	UnicredsPath   string
}
