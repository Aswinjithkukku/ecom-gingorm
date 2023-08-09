package utils

import (
	"log"
	"net/smtp"
	"os"
)

func SendEmailWithoutHTML(toAddress []string, subject string, content string) {

	auth := smtp.PlainAuth(
		"",
		os.Getenv("SMTP_EMAIL"),
		os.Getenv("SMTP_SECRET_APP_PASSWORD"),
		"smtp.gmail.com",
	)

	message := "Subject: " + subject + "\n" + content + "\n" + "Thank you!!"

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		os.Getenv("SMTP_EMAIL"),
		toAddress,
		[]byte(message),
	)

	if err != nil {
		log.Fatalf("failed to send mail: %s", err)
	}
}
