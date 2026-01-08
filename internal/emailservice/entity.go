package emailservice

import (
	"os"
)

type EmailConfig struct {
	SMTPHost    string
	SMTPPort    string
	SenderEmail string
	SenderName  string
	AppPassword string
	FrontendURL string
}

type EmailTemplate struct {
	Subject     string
	Heading     string
	Greeting    string
	MainMessage string
	ButtonText  string
	ButtonURL   string
	FooterNote  string
	ExpiryNote  string
}

func DefaultEmailConfig() EmailConfig {
	return EmailConfig{
		SMTPHost:    os.Getenv("SMTP_HOST"),
		SMTPPort:    os.Getenv("SMTP_PORT"),
		SenderEmail: os.Getenv("SMTP_USER"),
		SenderName:  os.Getenv("STMP_SENDER_NAME"),
		AppPassword: os.Getenv("SMTP_PASS"),
		FrontendURL: os.Getenv("FRONTEND_URL"),
	}
}
