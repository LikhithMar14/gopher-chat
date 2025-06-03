package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/vorobeyme/mailtrap-go/mailtrap"
)

type MailTrapClient struct {
	apiKey    string
	fromEmail string
	client    *mailtrap.SendingClient
}

func NewMailTrap(apiKey, fromEmail string) (*MailTrapClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("mailtrap API key is required")
	}

	if fromEmail == "" {
		return nil, fmt.Errorf("from email is required")
	}

	// Basic email validation
	if !strings.Contains(fromEmail, "@") || !strings.Contains(fromEmail, ".") {
		return nil, fmt.Errorf("invalid from email format: %s", fromEmail)
	}

	client, err := mailtrap.NewSendingClient(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create mailtrap client: %w", err)
	}

	return &MailTrapClient{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		client:    client,
	}, nil
}

func (m *MailTrapClient) Send(templateFile, username, email string, data interface{}, isSandbox bool) error {
	// Validate inputs
	if templateFile == "" {
		return fmt.Errorf("template file is required")
	}

	if email == "" {
		return fmt.Errorf("recipient email is required")
	}

	// Basic email validation
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return fmt.Errorf("invalid recipient email format: %s", email)
	}

	// Normalize username if empty
	if username == "" {
		username = strings.Split(email, "@")[0]
	}

	// Parse template
	tmpl, err := template.ParseFS(Fs, "templates/"+templateFile)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create buffers for template execution
	subject := new(bytes.Buffer)
	body := new(bytes.Buffer)

	// Execute templates
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return fmt.Errorf("failed to execute subject template: %w", err)
	}

	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return fmt.Errorf("failed to execute body template: %w", err)
	}

	// Validate template output
	if subject.Len() == 0 {
		return fmt.Errorf("subject template produced empty result")
	}

	if body.Len() == 0 {
		return fmt.Errorf("body template produced empty result")
	}

	// Create a new email request
	emailReq := &mailtrap.SendEmailRequest{
		From: mailtrap.EmailAddress{
			Email: m.fromEmail,
			Name:  FromName,
		},
		To: []mailtrap.EmailAddress{
			{
				Email: email,
				Name:  username,
			},
		},
		Subject:  subject.String(),
		Text:     "", // We're using HTML only
		HTML:     body.String(),
		Category: "user_notification",
	}

	// If using sandbox mode, add it to headers
	if isSandbox {
		if emailReq.Headers == nil {
			emailReq.Headers = make(map[string]string)
		}
		emailReq.Headers["X-Sandbox-Mode"] = "true"
	}

	// Send the email with retries
	for i := 0; i < maxRetries; i++ {
		_, _, err := m.client.Send(emailReq)
		if err == nil {
			log.Printf("Successfully sent email to %s", email)
			return nil
		}

		log.Printf("Failed to send email to %v, attempt %d of %d", email, i+1, maxRetries)
		log.Printf("Error: %v", err)

		if i < maxRetries-1 {
			// Wait before retrying with exponential backoff
			backoffDuration := time.Duration(i+1) * time.Second
			log.Printf("Retrying in %v", backoffDuration)
			time.Sleep(backoffDuration)
		}
	}

	return fmt.Errorf("failed to send email to %v after %d attempts", email, maxRetries)
}
