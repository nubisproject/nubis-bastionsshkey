package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"log"
)

type CreateIAMUserResult struct {
	Username  string
	AccessKey string
	SecretKey string
}

// encryptMailBody retrieves the PGP fingerprint of a recipient from ldap, then
// queries the gpg server to retrieve the public key and encrypts the body with it.
func EncryptMailBody(origBody []byte, key []byte, rcpt string) (body []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("encryptMailBody-> %v", e)
		}
	}()
	if err != nil {
		panic(err)
	}
	el, err := openpgp.ReadArmoredKeyRing(bytes.NewBuffer(key))
	if err != nil {
		panic(err)
	}
	encbuf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(encbuf, el, nil, nil, nil)
	if err != nil {
		panic(err)
	}
	_, err = w.Write([]byte(origBody))
	if err != nil {
		panic(err)
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}
	armbuf := bytes.NewBuffer(nil)
	w, err = armor.Encode(armbuf, "PGP MESSAGE", nil)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(encbuf.Bytes())
	if err != nil {
		panic(err)
	}
	w.Close()
	body = armbuf.Bytes()
	return
}
func DeleteIAMUser(config Configuration, username string) (bool, error) {
	sess := GetSession(config)
	svc := iam.New(sess)
	accessKeysParams := &iam.ListAccessKeysInput{
		UserName: aws.String(username),
	}
	accessKeysResult, _ := svc.ListAccessKeys(accessKeysParams)
	for _, accessKey := range accessKeysResult.AccessKeyMetadata {
		deleteAccessKeyParams := &iam.DeleteAccessKeyInput{
			AccessKeyId: aws.String(*accessKey.AccessKeyId),
			UserName:    aws.String(username),
		}

		_, deleteAccessKeyErr := svc.DeleteAccessKey(deleteAccessKeyParams)

		if deleteAccessKeyErr != nil {
			fmt.Printf("Unable to delete AccessKeyId: %s with error: %s", accessKey.AccessKeyId, deleteAccessKeyErr)
		}
	}

	params := &iam.DeleteUserInput{
		UserName: aws.String(username), // Required
	}
	_, err := svc.DeleteUser(params)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	return true, nil

}
func CreateIAMUser(config Configuration, username string, path string) (CreateIAMUserResult, error) {
	sess := GetSession(config)

	svc := iam.New(sess)
	params := &iam.CreateUserInput{
		UserName: aws.String(username), // Required
		Path:     aws.String(path),
	}
	resp, err := svc.CreateUser(params)

	if err != nil {
		fmt.Println(err.Error())
		return CreateIAMUserResult{}, err
	}
	IAMUserResult := CreateIAMUserResult{*resp.User.UserName, "", ""}
	createAccessKeysParams := &iam.CreateAccessKeyInput{
		UserName: aws.String(username),
	}
	createAccessKeyResp, createAccessKeyErr := svc.CreateAccessKey(createAccessKeysParams)
	if createAccessKeyErr != nil {
		fmt.Println(err.Error())
		return CreateIAMUserResult{}, err
	}
	IAMUserResult.AccessKey = *createAccessKeyResp.AccessKey.AccessKeyId
	IAMUserResult.SecretKey = *createAccessKeyResp.AccessKey.SecretAccessKey
	return IAMUserResult, nil

}

func GetAllIAMUsers(config Configuration) (*iam.ListUsersOutput, error) {
	sess := GetSession(config)

	svc := iam.New(sess)

	params := &iam.ListUsersInput{
		PathPrefix: aws.String("/"),
	}
	resp, err := svc.ListUsers(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return nil, err
	}

	// Pretty-print the response data.
	//fmt.Println(resp)
	return resp, nil

}

// Helper function to get ARN of a particular IAM user
//
// Example: arn, err := GetUserArn(conf, "test-user")
func GetUserArn(config Configuration, username string) (Arn string, err error) {
	sess := GetSession(config)
	svc := iam.New(sess)

	params := &iam.GetUserInput{
		UserName: aws.String(username),
	}
	resp, err := svc.GetUser(params)
	if err != nil {
		log.Printf("User %s not found", username)
		return "", err
	}
	return *resp.User.Arn, nil
}

// Helper function to get arn of an IAM Role
//
// Example: arn, err := GetRoleArn(conf, "test-role")
func GetRoleArn(config Configuration, rolename string) (Arn string, err error) {
	sess := GetSession(config)
	svc := iam.New(sess)

	params := &iam.GetRoleInput{
		RoleName: aws.String(rolename), // Required
	}
	resp, err := svc.GetRole(params)
	if err != nil {
		log.Printf("Role %s does not exist", rolename)
		return "", err
	}
	return *resp.Role.Arn, nil
}

// Create role
func CreateRole(config Configuration, rolename string, userarn string, rolepath string) error {

	sess := GetSession(config)
	svc := iam.New(sess)

	// Policy doc for this is a constant from the iam_roles.go file
	policydoc := fmt.Sprintf(AssumeRolePolicy, userarn)
	params := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(policydoc), //required
		RoleName:                 aws.String(rolename),
		Path:                     aws.String(rolepath),
	}
	_, err := svc.CreateRole(params)
	if err != nil {
		log.Printf("Unable to create role %s: %v", rolename, err.Error())
		return err
	}
	return nil
}

// Main function to attach policy to a role
func AttachPolicy(config Configuration, rolearn string, rolename string) (*iam.AttachRolePolicyOutput, error) {
	sess := GetSession(config)
	svc := iam.New(sess)

	params := &iam.AttachRolePolicyInput{
		PolicyArn: aws.String(rolearn),  // Required
		RoleName:  aws.String(rolename), // Required
	}

	resp, err := svc.AttachRolePolicy(params)
	if err != nil {
		log.Printf("Unable to attach policy %s: %v", rolearn, err.Error())
		return &iam.AttachRolePolicyOutput{}, err
	}
	return resp, nil
}

// Attaches a readonly policy to a role that we specify
// The ReadOnlyPolicyArn is a constant value that comes with
// every AWS account
func AttachReadOnlyPolicy(config Configuration, rolename string) (*iam.AttachRolePolicyOutput, error) {
	resp, err := AttachPolicy(config, ReadOnlyPolicyArn, rolename)
	if err != nil {
		log.Printf("Unable to attach policy %s: %v", rolename, err.Error())
		return &iam.AttachRolePolicyOutput{}, err
	}
	return resp, nil
}

// Attaches an admin policy to a rolename
// The Administrator Policy is a constant value
// that comes with the AWS account
func AttachAdminPolicy(config Configuration, rolename string) (*iam.PutRolePolicyOutput, error) {
	sess := GetSession(config)
	svc := iam.New(sess)

	params := &iam.PutRolePolicyInput{
		PolicyDocument: aws.String(AdminPolicy), //required, this is a constant
		PolicyName:     aws.String(rolename),
		RoleName:       aws.String(rolename),
	}
	resp, err := svc.PutRolePolicy(params)
	if err != nil {
		log.Printf("Unable to attach inline policy %s: %v", rolename, err.Error())
		return &iam.PutRolePolicyOutput{}, err
	}
	return resp, nil
}

// Attaches admin group
func AttachAdminGroup(config Configuration, username string) {
	sess := GetSession(config)
	svc := iam.New(sess)

	params := &iam.AddUserToGroupInput{
		GroupName: aws.String("Administrators"),
		UserName:  aws.String(username),
	}

	_, err := svc.AddUserToGroup(params)
	if err != nil {
		log.Printf("Unable to add user %s to group Administrator: %v", username, err.Error())
	}
	return
}

func PolicyEnforcer(config Configuration, username string, path string) {

	if path == "/nubis/admin/" {
		if noop {
			log.Printf("NOOP: Creating reate role: %s with admin privilege", username)
		} else {
			log.Printf("Creating role: %s", username)
			userarn, usererr := GetUserArn(config, username)
			if usererr != nil {
				log.Fatalf("Unable to get user arn %s: %v", username, usererr.Error())
			}
			err := CreateRole(config, username, userarn, path)
			if err != nil {
				log.Printf("Not creating role")
			}
			//_, attacherr := AttachAdminPolicy(config, username)
			//if attacherr != nil {
			//	log.Printf("Unable to attach admin policy for %s", username)
			//}
			//AttachAdminGroup(config, username)
		}
	} else if path == "/nubis/readonly/" {
		if noop {
			log.Printf("NOOP: Attempting to attach readonly role to user %s", username)
		}

	} else {
		log.Fatalf("Invalid IAM path: %s", path)
	}
}
