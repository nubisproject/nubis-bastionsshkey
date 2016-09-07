package main

import (
	"testing"
)

func TestAddingUserPathToUserPathListEmpty(t *testing.T) {
	username := "username"
	path := "global-admins"
	tmp := UserPath{username, path}
	list := UserPathList{}
	list.add(tmp)
	if len(list.list) != 1 {
		t.Errorf("UserPath not added to UserPathList")
	}
	if list.list[0].username != username {
		t.Errorf("username not set correctly")
	}
	if list.list[0].path != path {
		t.Errorf("path not set correctly")
	}
}

func TestAddingExistingUser(t *testing.T) {
	username := "username"
	path := "global-admins"
	tmp := UserPath{username, path}
	tmp2 := UserPath{username, path}
	list := UserPathList{}
	list.add(tmp)
	list.add(tmp2)
	if len(list.list) != 1 {
		t.Errorf("UserPath.List incorrectly allowed a duplicate user")
	}
}

func TestGetPathByUsername(t *testing.T) {
	username := "username"
	path := "global-admins"
	tmp := UserPath{username, path}
	list := UserPathList{}
	list.add(tmp)
	retPath := list.getPathByUsername(username)
	if retPath != path {
		t.Errorf("Unable to getPathByUsername")
	}
}
