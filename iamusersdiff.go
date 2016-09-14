package main

import (
	"github.com/aws/aws-sdk-go/service/iam"
)

type IAMUsersDiff struct {
	ListUsersOutput *iam.ListUsersOutput
	LDAPUserList    []string
	IgnorePath      []string
	IgnoreUser      []string
}

func UsernameInLDAP(s []string, e string) bool {
	return StringInSlice(s, e)
}
func UsernameInIAM(iamListUsersOutput *iam.ListUsersOutput, e string) bool {
	iamUsernames := []string{}
	for _, iamUser := range iamListUsersOutput.Users {
		iamUsernames = append(iamUsernames, *iamUser.UserName)
	}
	if len(iamUsernames) > 0 {
		return StringInSlice(iamUsernames, e)
	}
	return false
}
func (r *IAMUsersDiff) UserInIgnoreUser(username string) bool {
	return StringInSlice(r.IgnoreUser, username)
}

func (r *IAMUsersDiff) PathInIgnorePath(path string) bool {
	return StringInSlice(r.IgnorePath, path)
}

func (r *IAMUsersDiff) getUsersToAdd() []string {
	returnList := []string{}
	for _, LDAPUser := range r.LDAPUserList {
		if r.UserInIgnoreUser(LDAPUser) == true {
			continue
		}
		if UsernameInIAM(r.ListUsersOutput, LDAPUser) == false {
			returnList = append(returnList, LDAPUser)
		}
	}
	return returnList
}

func (r *IAMUsersDiff) getUsersToRemove() []string {
	returnList := []string{}
	for _, IAMUser := range r.ListUsersOutput.Users {
		if r.UserInIgnoreUser(*IAMUser.UserName) == true {
			continue
		}
		if r.PathInIgnorePath(*IAMUser.Path) == true {
			continue
		}
		if UsernameInLDAP(r.LDAPUserList, *IAMUser.UserName) == false {
			returnList = append(returnList, *IAMUser.UserName)
		}
	}
	return returnList
}
