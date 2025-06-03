package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)



type SendGridMailer struct {
	fromEmail string
	apiKey 	   string
	client     *sendgrid.Client
}


func NewSendgrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey: apiKey,
		client: client,
	}
}

func (m *SendGridMailer) Send(templateFile, username, email string , data interface{} , isSandbox bool) error {
	from := mail.NewEmail(FromName,m.fromEmail)
	to := mail.NewEmail(username,email)

	tmpl , err := template.ParseFS(Fs, "templates/"+templateFile)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}
	subject := new(bytes.Buffer)
	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}


	message := mail.NewSingleEmail(from, subject.String(), to , "", body.String())

	message.SetMailSettings(
		&mail.MailSettings{
			SandboxMode: &mail.Setting{
				Enable: &isSandbox,
			},
		},
	)

	for i := 0; i < maxRetries; i++ {
		response, err := m.client.Send(message)
		if err != nil {
			log.Printf("Failed to send email to %v, attempt %d of %d", email, i+1, maxRetries)
			log.Printf("Error: %v", err)

			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		if response.StatusCode >= 200 && response.StatusCode < 300 {
			return nil
		} else {
			log.Printf("Failed to send email to %v, attempt %d of %d", email, i+1, maxRetries)
			log.Printf("Status Code: %d", response.StatusCode)
			log.Printf("Response Body: %s", response.Body)

			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
	}

	return fmt.Errorf("failed to send email to %v after %d attempts", email, maxRetries)
}