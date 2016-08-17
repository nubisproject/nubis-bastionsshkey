package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

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
