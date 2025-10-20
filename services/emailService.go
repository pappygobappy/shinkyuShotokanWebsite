package services

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

func SendPasswordResetEmail(toEmail, resetURL string) error {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		host = "smtp.gmail.com"
	}

	port := 587
	if portStr := os.Getenv("SMTP_PORT"); portStr != "" {
		parsedPort, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("invalid SMTP_PORT: %w", err)
		}
		port = parsedPort
	}

	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	if username == "" || password == "" {
		return fmt.Errorf("email configuration incomplete: missing username or password")
	}

	if from == "" {
		from = username
	}

	auth := smtp.PlainAuth("", username, password, host)
	addr := fmt.Sprintf("%s:%d", host, port)

	headers := []string{
		"From: " + from,
		"To: " + toEmail,
		"Subject: Shinkyu Shotokan password reset",
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=\"utf-8\"",
	}

	body := fmt.Sprintf(
		"We received a request to reset the password for your Shinkyu Shotokan account.\n\n"+
		"Reset your password using the link below (valid for one hour):\n%s\n\n"+
		"If you did not request a password reset, you can ignore this email.\n",
		resetURL,
	)
	message := strings.Join(headers, "\r\n") + "\r\n\r\n" + body

	return smtp.SendMail(addr, auth, from, []string{toEmail}, []byte(message))
}
