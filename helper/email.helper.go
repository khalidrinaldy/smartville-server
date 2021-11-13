package helper

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

func SendEmail(to []string, subject, message string) error {
	body := "From: " + os.Getenv("EMAIL_SENDER_NAME") + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		message
	
	auth := smtp.PlainAuth("", os.Getenv("EMAIL_ADDR"), os.Getenv("EMAIL_PASSWORD"), os.Getenv("EMAIL_HOST"))
	smtpAddr := fmt.Sprintf("%s:%s", os.Getenv("EMAIL_HOST"), os.Getenv("EMAIL_PORT"))

	err := smtp.SendMail(smtpAddr, auth, os.Getenv("EMAIL_ADDR"), to, []byte(body))
    if err != nil {
        return err
    }
	return nil
}