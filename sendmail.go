package main

import (
	"fmt"
	"log"
	"net/smtp"
)

func SendWelcomeMail(config Configuration, dest string, message []byte) {
	dialString := fmt.Sprintf("%s:%s", config.AWS.SMTPHostname, config.AWS.SMTPPort)
	fmt.Println(dialString)
	fmt.Println(dest)
	fmt.Println(config.AWS.SMTPFromAddress)
	// Set the sender and recipient.
	msg := []byte("To: " + dest + "\r\n" +
		"Subject: Nubis Account Credentials!\r\n" +
		"\r\n" +
		"Please decrypt the following message for AWS AccessKeyID & Secret Key.\r\n" +
		string(message) + "\r\n")

	auth := smtp.PlainAuth(
		"",
		config.AWS.SMTPUsername,
		config.AWS.SMTPPassword,
		config.AWS.SMTPHostname,
	)
	sendErr := smtp.SendMail(
		dialString,
		auth,
		config.AWS.SMTPFromAddress,
		[]string{dest},
		msg,
	)
	if sendErr != nil {
		log.Fatal(sendErr)
	}
}
