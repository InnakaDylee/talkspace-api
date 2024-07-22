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

func EmailOTP(to []string, template string, data interface{}) (bool, error) {
	config, err := configs.LoadConfig()
	if err != nil {
		logrus.Fatalf("failed to load smtp configuration: %v", err)
	}

	m := mail.NewMessage()
	m.SetHeader("From", config.SMTP.SMTP_USER)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", "TalkSpace : One Time Password")

	emailContent := strings.Replace(template, "{{.OTP}}", data.(string), -1)

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

func SendEmailOTP(emailAddress string, otp string) {
	go func() {
		filePath := "utils/helper/email/template/one-time-password.html"
		emailTemplate, errRead := os.ReadFile(filePath)
		if errRead != nil {
			log.Printf("failed to prepare email template: %v", errRead)
			return
		}

		emailTemplateStr := string(emailTemplate)
		emailContent := strings.Replace(emailTemplateStr, "{{.OTP}}", otp, -1)

		_, errEmail := EmailOTP([]string{emailAddress}, emailContent, otp)
		if errEmail != nil {
			log.Printf("failed to send otp email: %v", errEmail)
		}
	}()
}
