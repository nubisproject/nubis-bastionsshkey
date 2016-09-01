package main

import (
	"testing"
)

var user1 = LDAPUserObject{
	"mail=testing@mozilla.com,o=com,dc=mozilla",
	"testing",
	"testing@mozilla.com",
	[]string{},
	[]byte{},
}

var group1 = []LDAPUserObject{user1}

func TestUserInGroup(t *testing.T) {
	username := "testing"
	if UserInGroup("testing", group1) == false {
		t.Errorf("user %s not found in group", username)
	}
	if UserInGroup("unknown", group1) == true {
		t.Errorf("user %s found in group", username)
	}
}

func TestGetLDAPUserObjectFromGroupUserFound(t *testing.T) {
	username := "testing"
	retUser, found := GetLDAPUserObjectFromGroup(username, group1)
	if found == false {
		t.Errorf("found incorrectly set to false")
	}
	if retUser.Uid != "testing" {
		t.Errorf("user %s not found in group", username)
	}
}

func TestGetLDAPUserObjectFromGroupUserNotFound(t *testing.T) {
	username := "unknown"
	retUser, found := GetLDAPUserObjectFromGroup(username, group1)
	if found == true {
		t.Errorf("found incorrectly set to false")
	}
	if retUser.Uid != "" {
		t.Errorf("user %s incorrectly found in group", username)
	}
}
