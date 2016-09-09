package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func GetSession(config Configuration) *session.Session {
	var sess *session.Session
	if config.AWS.AccessKey == "" {
		sess = session.New(&aws.Config{
			Region: aws.String(config.AWS.Region),
		})
	} else {
		sess = session.New(&aws.Config{
			Region:      aws.String(config.AWS.Region),
			Credentials: credentials.NewStaticCredentials(config.AWS.AccessKey, config.AWS.SecretKey, ""),
		})
	}
	return sess

}
