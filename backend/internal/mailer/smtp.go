package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"gopkg.in/gomail.v2"
)

type SMTPMailer struct {
	host     string
	port     int
	username string
	password string
	email    string
}

func NewSMTPMailer(host string, port int, username, password, email string) *SMTPMailer {
	return &SMTPMailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
		email:    email,
	}
}

func (m *SMTPMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return fmt.Errorf("failed to execute subject template: %w", err)
	}
	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return fmt.Errorf("failed to execute body template: %w", err)
	}

	message := gomail.NewMessage()
	message.SetHeader("From", m.email)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject.String())
	message.SetBody("text/html", body.String())

	dialer := gomail.NewDialer(m.host, m.port, m.username, m.password)

	for i := 0; i < maxRetries; i++ {
		if isSandbox {
			log.Printf("Sandbox mode: Email to %s not sent. Subject: %s", email, subject.String())
			log.Printf("Email body: %s", body.String())
			return nil
		} else {
			log.Printf("Attempting to send email to %s (attempt %d)", email, i+1)
		}

		if err := dialer.DialAndSend(message); err != nil {
			if i == maxRetries-1 {
				return err
			}
			continue
		}

		log.Printf("Email sent to %s", email)
		return nil
	}

	return fmt.Errorf("failed to send email to %s after %d attempts", email, maxRetries)
}
