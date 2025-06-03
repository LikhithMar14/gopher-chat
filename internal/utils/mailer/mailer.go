package mailer

import "embed"

//go:embed templates/*
var Fs embed.FS

const (
	FromName = "Gopher Chat"
	maxRetries = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

type Client interface {
	Send(templateFile, username, email string, data interface{}, isSandbox bool) error
}