package main

import (
	"github.com/aws/aws-sdk-go/service/iam"
)

type RemoveIAMUsers struct {
	ListUsersOutput *iam.ListUsersOutput
	LDAPUserList    []string
	IgnorePath      []string
	IgnoreUser      []string
}

func UsernameInLDAP(s []string, e string) bool {
	return StringInSlice(s, e)
}
func (r *RemoveIAMUsers) UserInIgnoreUser(username string) bool {
	return StringInSlice(r.IgnoreUser, username)
}

func (r *RemoveIAMUsers) PathInIgnorePath(path string) bool {
	return StringInSlice(r.IgnorePath, path)
}

func (r *RemoveIAMUsers) getUsersToRemove() []string {
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
