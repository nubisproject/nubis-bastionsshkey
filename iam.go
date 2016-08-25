package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type CreateIAMUserResult struct {
	Username  string
	AccessKey string
	SecretKey string
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
