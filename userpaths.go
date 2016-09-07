package main

import ()

type UserPath struct {
	username string
	path     string
}

type UserPathList struct {
	list []UserPath
}

func (u *UserPathList) getPathByUsername(username string) string {
	for _, x := range u.list {
		if x.username == username {
			return x.path
		}
	}
	return ""
}

func (u *UserPathList) contains(username string) bool {
	for _, x := range u.list {
		if x.username == username {
			return true
		}
	}
	return false
}

func (u *UserPathList) add(userPath UserPath) {
	if u.contains(userPath.username) == false {
		u.list = append(u.list, userPath)
	}
}
