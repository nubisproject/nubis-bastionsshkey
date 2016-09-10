package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
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
	return fmt.Sprintf("%s", *resp.Role.Arn), nil
}

// Create role
func CreateRole(config Configuration, rolename string, userarn string, rolepath string) (*iam.CreateRoleOutput, error) {

	sess := GetSession(config)
	svc := iam.New(sess)

	// Policy doc for this is a constant from the iam_roles.go file
	policydoc := fmt.Sprintf(AssumeRolePolicy, userarn)
	params := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(policydoc), //required
		RoleName:                 aws.String(rolename),
		Path:                     aws.String(rolepath),
	}
	resp, err := svc.CreateRole(params)
	if err != nil {
		log.Fatalf("Unable to create role, role %s might already exist", rolename)
		return nil, err
	}
	return resp, nil
}

// Attach an inline policy to the role
//
// Example: resp, err := AttachInlinePolicy(conf, rolepolicystring, "test-role")
func AttachInlinePolicy(config Configuration, rolepolicy string, rolename string) (*iam.PutRolePolicyOutput, error) {

	sess := GetSession(config)
	svc := iam.New(sess)

	// rolepolicy should be a json document
	params := &iam.PutRolePolicyInput{
		PolicyDocument: aws.String(rolepolicy), //required
		PolicyName:     aws.String(rolename),
		RoleName:       aws.String(rolename),
	}
	resp, err := svc.PutRolePolicy(params)
	if err != nil {
		log.Fatalf("Attach policy error: %v", err.Error())
		return nil, err
	}
	return resp, nil
}

// Attaches a readonly Policy to the IAM user, the Readonly policy
// that we have is actually a managed policy that comes standard
// with every AWS account
//
// Example: resp, err := AttachReadOnlyPolicy(conf, "testuser")
func AttachReadOnlyPolicy(config Configuration, username string) (*iam.AttachUserPolicyOutput, error) {
	sess := GetSession(config)
	svc := iam.New(sess)

	// Pretty sure this is a standard arn for readonly access
	ReadOnlyArn := "arn:aws:iam::aws:policy/ReadOnlyAccess"
	params := &iam.AttachUserPolicyInput{
		PolicyArn: aws.String(ReadOnlyArn),
		UserName:  aws.String(username),
	}
	resp, err := svc.AttachUserPolicy(params)
	if err != nil {
		log.Fatalf("Unable to attach readonly policy %s to user %s: %v", ReadOnlyArn, username, err.Error())
		return nil, err
	}
	return resp, nil
}

func PolicyEnforcer(config Configuration, username string) {
	// This is pretty inefficient
	for _, x := range config.LdapServer.IAMGroupMapping {
		priv := x.PrivilegeLevel
		if priv == "admin" {
			if noop {
				log.Printf("NOOP: Will attempt to create role: %s with admin privilege", username)
			} else {
				log.Printf("Calling CreateRole function")
				//_, err := CreateRole(config, userarn)
				//_, err := CreateRole(config, username, userarn, x.IAMPath)
				//if err != nil {
				//	log.Fatalf("Unable to create role: %s for user %s", username, username)
				//}
				//_, err := AttachInlinePolicy(config, AdminPolicy, username)
			}

		} else if priv == "readonly" {
			if noop {
				log.Printf("NOOP: Attempt to attach readonly role to user %s", username)
			}

		} else {
			log.Fatalf("Invalid PrivelegeLevel value: %s", priv)
		}
	}
}
