package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/iam"
	"testing"
)

func generateRandomListUsersOutput(count int) iam.ListUsersOutput {
	users := []*iam.User{}
	for i := 0; i < count; i++ {
		userName := fmt.Sprintf("Testing%d", i)
		arn := fmt.Sprintf("arn:uri:1234567890%d", i)
		path := "/"
		tmpUser := iam.User{
			UserName: &userName,
			Arn:      &arn,
			Path:     &path,
		}
		users = append(users, &tmpUser)
	}
	IAMList := iam.ListUsersOutput{
		Users: users,
	}

	return IAMList

}

func TestIAMUserNotExistsInLDAPList(t *testing.T) {
	LDAPUserList := []string{"realuser"}
	IAMList := generateRandomListUsersOutput(1)
	r := RemoveIAMUsers{
		&IAMList,
		LDAPUserList,
		[]string{},
		[]string{},
	}
	listToRemove := r.getUsersToRemove()
	if len(listToRemove) != 1 {
		t.Errorf("Incorrect number of users in listToRemove")
	}

	if listToRemove[0] != "Testing0" {
		t.Errorf("Incorrect entry in listToRemove")
	}

}
func TestIAMUserNotExistsInLDAPListWithIgnorePath(t *testing.T) {
	LDAPUserList := []string{"realuser"}
	ignorePath := "/path/to/ignore/"
	IAMList := generateRandomListUsersOutput(2)
	IAMList.Users[1].Path = &ignorePath
	r := RemoveIAMUsers{
		&IAMList,
		LDAPUserList,
		[]string{ignorePath},
		[]string{},
	}
	listToRemove := r.getUsersToRemove()
	if len(listToRemove) != 1 {
		t.Errorf("Incorrect number of users in listToRemove")
	}

	if listToRemove[0] != "Testing0" {
		t.Errorf("Incorrect entry in listToRemove")
	}

}
func TestIAMUserNotExistsInLDAPListWithIgnoreUser(t *testing.T) {
	LDAPUserList := []string{"realuser"}
	IAMList := generateRandomListUsersOutput(2)
	r := RemoveIAMUsers{
		&IAMList,
		LDAPUserList,
		[]string{},
		[]string{"Testing1"},
	}
	listToRemove := r.getUsersToRemove()
	if len(listToRemove) != 1 {
		t.Errorf("Incorrect number of users in listToRemove")
	}

	if listToRemove[0] != "Testing0" {
		t.Errorf("Incorrect entry in listToRemove")
	}

}
func TestIAMUserExistsInLDAPList(t *testing.T) {
	LDAPUserList := []string{"Testing0"}
	IAMList := generateRandomListUsersOutput(1)
	r := RemoveIAMUsers{
		&IAMList,
		LDAPUserList,
		[]string{},
		[]string{},
	}
	listToRemove := r.getUsersToRemove()
	if len(listToRemove) != 0 {
		t.Errorf("Incorrect number of users in listToRemove")
	}
}
