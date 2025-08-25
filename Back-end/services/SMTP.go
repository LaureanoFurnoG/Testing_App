package services

import (
	"net/smtp"
	"os"
	"fmt"
)

func SendEmail(to, otp string) error {
	from := os.Getenv("EMAIL_USER")
	pass := os.Getenv("EMAIL_PASS")

	msg := fmt.Sprintf("Subject: Your Code is %s", otp)

	return smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))
}
