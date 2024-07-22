package mailer

import (
	"log"
	"os"
	"strconv"
	"strings"
	"talkspace-api/app/configs"

	"github.com/sirupsen/logrus"
	"gopkg.in/mail.v2"
)

func EmailVerificationAccount(to []string, template string, data interface{}) (bool, error) {
	config, err := configs.LoadConfig()
	if err != nil {
		logrus.Fatalf("failed to load smtp configuration: %v", err)
	}

	m := mail.NewMessage()
	m.SetHeader("From", config.SMTP.SMTP_USER)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", "Account Verification")

	emailContent := strings.Replace(template, "{{.verification_account}}", data.(string), -1)

	m.SetBody("text/html", emailContent)

	SMTP_PORT, err := strconv.Atoi(config.SMTP.SMTP_PORT)
	if err != nil {
		return false, err
	}

	d := mail.NewDialer(
		config.SMTP.SMTP_HOST,
		SMTP_PORT,
		config.SMTP.SMTP_USER,
		config.SMTP.SMTP_PASS,
	)

	if err := d.DialAndSend(m); err != nil {
		return false, err
	}
	return true, nil
}

func SendEmailVerificationAccount(emailAddress string, token string) {
	go func() {
		verificationURL := "https://api.talkspace.my.id/verify-token?token=" + token
		filePath := "utils/helper/email/template/verification-account.html"
		emailTemplate, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("failed to prepare email template: %v", err)
			return
		}

		_, errEmail := EmailVerificationAccount([]string{emailAddress}, string(emailTemplate), verificationURL)
		if errEmail != nil {
			log.Printf("failed to send verification email: %v", errEmail)
		}
	}()
}
