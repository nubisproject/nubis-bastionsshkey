package main

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
