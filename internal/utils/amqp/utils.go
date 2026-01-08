package amqp

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	. "github.com/ishola-faazele/taskflow/internal/emailservice"
)

func PublishEmailMessage(ch *amqp.Channel, msg *EmailMessage) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal email message: %w", err)
	}

	err = ch.Publish(
		"",            // exchange
		"email_queue", // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         msgBytes,
		})
	if err != nil {
		return fmt.Errorf("failed to publish email message: %w", err)
	}

	return nil
}

// Helper methods to create messages
func NewMagicLinkMessage(toEmail, token, verifyURL string) (*EmailMessage, error) {
	payload := MagicLinkPayload{
		ToEmail:   toEmail,
		Token:     token,
		VerifyURL: verifyURL,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal magic link payload: %w", err)
	}

	return &EmailMessage{
		Type:    MessageTypeMagicLink,
		Payload: payloadBytes,
	}, nil
}

func NewInvitationMessage(toEmail, workspaceName, role, token, invitationURL string) (*EmailMessage, error) {
	payload := InvitationPayload{
		ToEmail:       toEmail,
		WorkspaceName: workspaceName,
		Role:          role,
		Token:         token,
		InvitationURL: invitationURL,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal invitation payload: %w", err)
	}

	return &EmailMessage{
		Type:    MessageTypeInvitation,
		Payload: payloadBytes,
	}, nil
}

func NewPasswordResetMessage(toEmail, token, resetURL string) (*EmailMessage, error) {
	payload := PasswordResetPayload{
		ToEmail:  toEmail,
		Token:    token,
		ResetURL: resetURL,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal password reset payload: %w", err)
	}

	return &EmailMessage{
		Type:    MessageTypePasswordReset,
		Payload: payloadBytes,
	}, nil
}

func NewCustomEmailMessage(toEmail string, template EmailTemplate) (*EmailMessage, error) {
	payload := CustomEmailPayload{
		ToEmail:  toEmail,
		Template: template,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal custom email payload: %w", err)
	}

	return &EmailMessage{
		Type:    MessageTypeCustom,
		Payload: payloadBytes,
	}, nil
}
