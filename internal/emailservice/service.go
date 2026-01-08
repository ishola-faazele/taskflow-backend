package emailservice

import (
	"fmt"
)

type EmailService struct {
	config EmailConfig
}

func NewEmailService(config EmailConfig) *EmailService {
	return &EmailService{
		config: config,
	}
}

// SendMagicLink sends a magic link email to the user
func (e *EmailService) SendMagicLink(toEmail, token, verifyURL string) error {
	// Construct the verification URL
	fullURL := fmt.Sprintf("%s%s%s", e.config.FrontendURL, verifyURL, token)

	template := EmailTemplate{
		Subject:     "Your Magic Link to Sign In",
		Heading:     "Sign in to TaskFlow",
		Greeting:    "Hello,",
		MainMessage: "Click the button below to sign in to your account. This link will expire in 15 minutes for security reasons.",
		ButtonText:  "Sign In Now",
		ButtonURL:   fullURL,
		FooterNote:  "If you didn't request this email, you can safely ignore it.",
		ExpiryNote:  "This link will expire in 15 minutes.",
	}

	body := e.createEmailHTML(template)
	return e.sendEmail(toEmail, template.Subject, body)
}

// SendInvitationLink sends an invitation email to join a workspace
func (e *EmailService) SendInvitationLink(toEmail, workspaceName, role, token, invitationURL string) error {
	// Construct the invitation URL
	fullURL := fmt.Sprintf("%s%s%s", e.config.FrontendURL, invitationURL, token)

	template := EmailTemplate{
		Subject:     fmt.Sprintf("You've been invited to join %s", workspaceName),
		Heading:     fmt.Sprintf("Join %s on TaskFlow", workspaceName),
		Greeting:    "Hello,",
		MainMessage: fmt.Sprintf("You have been invited to join the %s workspace as a %s. Click the button below to accept the invitation and get started.", workspaceName, role),
		ButtonText:  "Accept Invitation",
		ButtonURL:   fullURL,
		FooterNote:  "If you don't want to accept this invitation, you can safely ignore this email.",
		ExpiryNote:  "This invitation will expire in 24 hours.",
	}

	body := e.createEmailHTML(template)
	return e.sendEmail(toEmail, template.Subject, body)
}

// SendPasswordResetLink sends a password reset email
func (e *EmailService) SendPasswordResetLink(toEmail, token, resetURL string) error {
	// Construct the reset URL
	fullURL := fmt.Sprintf("%s%s%s", e.config.FrontendURL, resetURL, token)

	template := EmailTemplate{
		Subject:     "Reset Your Password",
		Heading:     "Password Reset Request",
		Greeting:    "Hello,",
		MainMessage: "We received a request to reset your password. Click the button below to create a new password.",
		ButtonText:  "Reset Password",
		ButtonURL:   fullURL,
		FooterNote:  "If you didn't request a password reset, you can safely ignore this email. Your password will remain unchanged.",
		ExpiryNote:  "This link will expire in 1 hour.",
	}
	body := e.createEmailHTML(template)

	return e.sendEmail(toEmail, template.Subject, body)
}

// SendCustomEmail sends an email with custom template
func (e *EmailService) SendCustomEmail(toEmail string, template EmailTemplate) error {
	body := e.createEmailHTML(template)
	return e.sendEmail(toEmail, template.Subject, body)
}
