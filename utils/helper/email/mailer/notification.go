package mailer

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"talkspace-api/app/configs"
	"text/template"

	"github.com/sirupsen/logrus"
	"gopkg.in/mail.v2"
)

func EmailNotificationAccount(to []string, templateContent string, data interface{}) (bool, error) {
	config, err := configs.LoadConfig()
	if err != nil {
		logrus.Fatalf("failed to load smtp configuration: %v", err)
	}

	m := mail.NewMessage()
	m.SetHeader("From", config.SMTP.SMTP_USER)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", "TalkSpace Notification")

	// Parse the template with data
	tmpl, err := template.New("emailTemplate").Parse(templateContent)
	if err != nil {
		return false, fmt.Errorf("failed to parse email template: %v", err)
	}

	var emailContent bytes.Buffer
	if err := tmpl.Execute(&emailContent, data); err != nil {
		return false, fmt.Errorf("failed to execute template: %v", err)
	}

	m.SetBody("text/html", emailContent.String())

	SMTP_PORT, err := strconv.Atoi(config.SMTP.SMTP_PORT)
	if err != nil {
		return false, fmt.Errorf("invalid SMTP port: %v", err)
	}

	d := mail.NewDialer(
		config.SMTP.SMTP_HOST,
		SMTP_PORT,
		config.SMTP.SMTP_USER,
		config.SMTP.SMTP_PASS,
	)

	if err := d.DialAndSend(m); err != nil {
		return false, fmt.Errorf("failed to send email: %v", err)
	}
	return true, nil
}

func SendEmailNotificationRegisterDoctor(fullname, licenseNumber, email, password string) {
	go func() {
		filePath := "utils/helper/email/template/register-doctor-success.html"
		emailTemplate, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("failed to load email template: %v", err)
			return
		}

		data := map[string]string{
			"Fullname":      fullname,
			"LicenseNumber": licenseNumber,
			"Email":         email,
			"Password":      password,
		}

		success, errEmail := EmailNotificationAccount([]string{email}, string(emailTemplate), data)
		if !success || errEmail != nil {
			log.Printf("failed to send notification email to %s: %v", email, errEmail)
		}
	}()
}

func SendEmailNotificationRegisterAccount(email string) {
	go func() {
		filePath := "utils/helper/email/template/register-success.html"
		emailTemplate, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("failed to load email template: %v", err)
			return
		}

		success, errEmail := EmailNotificationAccount([]string{email}, string(emailTemplate), nil)
		if !success || errEmail != nil {
			log.Printf("failed to send notification email to %s: %v", email, errEmail)
		}
	}()
}

func SendEmailNotificationLoginAccount(email string) {
	go func() {
		filePath := "utils/helper/email/template/login-success.html"
		emailTemplate, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("failed to load email template: %v", err)
			return
		}

		success, errEmail := EmailNotificationAccount([]string{email}, string(emailTemplate), nil)
		if !success || errEmail != nil {
			log.Printf("failed to send notification email to %s: %v", email, errEmail)
		}
	}()
}
