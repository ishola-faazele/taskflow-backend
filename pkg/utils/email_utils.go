package utils

import (
	"fmt"
	"net/smtp"
	"regexp"
	"strings"
)

func IsValidEmail(email string) bool {
	// Simple regex for email validation
	const emailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// EmailConfig holds SMTP configuration for Gmail
type EmailConfig struct {
	SMTPHost    string
	SMTPPort    string
	SenderEmail string
	SenderName  string
	AppPassword string // Gmail app password
	FrontendURL string // Base URL of your frontend
}

// EmailService handles email sending
type EmailService struct {
	config EmailConfig
}

// NewEmailService creates a new email service
func NewEmailService(config EmailConfig) *EmailService {
	return &EmailService{
		config: config,
	}
}

// SendMagicLink sends a magic link email to the user
func (e *EmailService) SendMagicLink(toEmail, token string) error {
	// Construct the verification URL
	verifyURL := fmt.Sprintf("%s/api/user/verify?token=%s", e.config.FrontendURL, token)

	// Create email subject and body
	subject := "Your Magic Link to Sign In"
	body := e.createMagicLinkHTML(toEmail, verifyURL)

	return e.sendEmail(toEmail, subject, body)
}

// createMagicLinkHTML creates a nicely formatted HTML email
func (e *EmailService) createMagicLinkHTML(email, verifyURL string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Magic Link Sign In</title>
</head>
<body style="margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #f4f4f4;">
    <table role="presentation" style="width: 100%%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; border-collapse: collapse; background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 20px 40px; text-align: center;">
                            <h1 style="margin: 0; color: #333333; font-size: 24px;">Sign in to TaskFlow</h1>
                        </td>
                    </tr>
                    
                    <!-- Body -->
                    <tr>
                        <td style="padding: 20px 40px;">
                            <p style="margin: 0 0 20px 0; color: #666666; font-size: 16px; line-height: 1.5;">
                                Hello,
                            </p>
                            <p style="margin: 0 0 20px 0; color: #666666; font-size: 16px; line-height: 1.5;">
                                Click the button below to sign in to your account. This link will expire in 1 hour for security reasons.
                            </p>
                            
                            <!-- Button -->
                            <table role="presentation" style="margin: 30px 0;">
                                <tr>
                                    <td align="center">
                                        <a href="%s" style="display: inline-block; padding: 16px 40px; background-color: #4F46E5; color: #ffffff; text-decoration: none; border-radius: 6px; font-size: 16px; font-weight: 600;">
                                            Sign In Now
                                        </a>
                                    </td>
                                </tr>
                            </table>
                            
                            <p style="margin: 20px 0 0 0; color: #999999; font-size: 14px; line-height: 1.5;">
                                Or copy and paste this link into your browser:
                            </p>
                            <p style="margin: 10px 0 0 0; color: #4F46E5; font-size: 14px; word-break: break-all;">
                                %s
                            </p>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 30px 40px; border-top: 1px solid #eeeeee;">
                            <p style="margin: 0; color: #999999; font-size: 12px; line-height: 1.5;">
                                If you didn't request this email, you can safely ignore it.
                            </p>
                            <p style="margin: 10px 0 0 0; color: #999999; font-size: 12px;">
                                This link will expire in 15 minutes.
                            </p>
                        </td>
                    </tr>
                </table>
                
                <!-- Footer text -->
                <table role="presentation" style="width: 600px; margin-top: 20px;">
                    <tr>
                        <td style="text-align: center; color: #999999; font-size: 12px;">
                            © 2024 TaskFlow. All rights reserved.
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`, verifyURL, verifyURL)
}

// sendEmail sends an email using Gmail SMTP
func (e *EmailService) sendEmail(to, subject, htmlBody string) error {
	// Set up authentication
	auth := smtp.PlainAuth(
		"",
		e.config.SenderEmail,
		e.config.AppPassword,
		e.config.SMTPHost,
	)

	// Compose the email
	from := fmt.Sprintf("%s <%s>", e.config.SenderName, e.config.SenderEmail)

	// Build email headers and body
	msg := e.buildEmailMessage(from, to, subject, htmlBody)

	// Send the email
	addr := fmt.Sprintf("%s:%s", e.config.SMTPHost, e.config.SMTPPort)
	err := smtp.SendMail(addr, auth, e.config.SenderEmail, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// buildEmailMessage constructs the email message with proper headers
func (e *EmailService) buildEmailMessage(from, to, subject, htmlBody string) string {
	var msg strings.Builder

	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	return msg.String()
}

// SendPlainTextMagicLink sends a plain text version (fallback)
func (e *EmailService) SendPlainTextMagicLink(toEmail, token string) error {
	verifyURL := fmt.Sprintf("%s/api/user/verify?token=%s", e.config.FrontendURL, token)

	subject := "Your Magic Link to Sign In"
	body := fmt.Sprintf(`Hello,

		Click the link below to sign in to your TaskFlow account:

		%s

		This link will expire in 1 hour for security reasons.

		If you didn't request this email, you can safely ignore it.

		Best regards,
		TaskFlow Team`,
		verifyURL)

	// Set up authentication
	auth := smtp.PlainAuth(
		"",
		e.config.SenderEmail,
		e.config.AppPassword,
		e.config.SMTPHost,
	)

	// Build plain text message
	from := fmt.Sprintf("%s <%s>", e.config.SenderName, e.config.SenderEmail)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		from, toEmail, subject, body)

	// Send email
	addr := fmt.Sprintf("%s:%s", e.config.SMTPHost, e.config.SMTPPort)
	err := smtp.SendMail(addr, auth, e.config.SenderEmail, []string{toEmail}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
