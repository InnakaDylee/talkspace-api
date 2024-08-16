
package mailer

import (
    "log"
    "os"
    "fmt"
    "talkspace-api/app/configs"
    "gopkg.in/mail.v2"
    "strconv"
    "strings"
    "github.com/sirupsen/logrus"
)

func EmailNotificationAccount(to []string, template string, data interface{}) (bool, error) {
    config, err := configs.LoadConfig()
    if err != nil {
        logrus.Fatalf("failed to load smtp configuration: %v", err)
    }

    m := mail.NewMessage()
    m.SetHeader("From", config.SMTP.SMTP_USER)
    m.SetHeader("To", to...)
    m.SetHeader("Subject", "TalkSpace Notification")

    emailContent := strings.Replace(template, "{{.Data}}", data.(string), -1)

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

func SendEmailNotificationRegisterDoctor(fullname, licenseNumber, email, password string) {
    go func() {

        filePath := "utils/helper/email/template/register-doctor-success.html"
        emailTemplate, err := os.ReadFile(filePath)
        if err != nil {
            log.Printf("failed to prepare email template: %v", err)
            return
        }

        data := fmt.Sprintf(
            "Fullname: %s<br>License Number: %s<br>Email: %s<br>Password: %s",
            fullname, licenseNumber, email, password,
        )

        _, errEmail := EmailNotificationAccount([]string{email}, string(emailTemplate), data)
        if errEmail != nil {
            log.Printf("failed to send notification email: %v", errEmail)
        }
    }()
}


func SendEmailNotificationRegisterAccount(email string) {
	go func() {

		filePath := "utils/helper/email/template/register-success.html"
		emailTemplate, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("failed to prepare email template: %v", err)
			return
		}

		_, errEmail := EmailNotificationAccount([]string{email}, string(emailTemplate), "")
		if errEmail != nil {
			log.Printf("failed to send notification email: %v", errEmail)
		}
	}()
}

func SendEmailNotificationLoginAccount(email string) {
	go func() {

		filePath := "utils/helper/email/template/login-success.html"
		emailTemplate, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("failed to prepare email template: %v", err)
			return
		}

		_, errEmail := EmailNotificationAccount([]string{email}, string(emailTemplate), "")
		if errEmail != nil {
			log.Printf("failed to send notification email: %v", errEmail)
		}
	}()
}
