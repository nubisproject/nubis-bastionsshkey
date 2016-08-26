package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
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
	sess := session.New(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials(config.AWS.AccessKey, config.AWS.SecretKey, ""),
	})
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
	sess := session.New(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials(config.AWS.AccessKey, config.AWS.SecretKey, ""),
	})

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
	sess := session.New(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials(config.AWS.AccessKey, config.AWS.SecretKey, ""),
	})

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
