package emailservice

import (
	"encoding/json"
	"fmt"
)


// MessageType represents the type of email message
type MessageType string

const (
	MessageTypeMagicLink     MessageType = "email.magic_link"
	MessageTypeInvitation    MessageType = "email.invitation"
	MessageTypePasswordReset MessageType = "email.password_reset"
	MessageTypeCustom        MessageType = "email.custom"
)

// EmailMessage is the unified message type for all email queue messages
type EmailMessage struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Specific payload types for each message type
type MagicLinkPayload struct {
	ToEmail   string `json:"to_email"`
	Token     string `json:"token"`
	VerifyURL string `json:"verify_url"`
}

type InvitationPayload struct {
	ToEmail       string `json:"to_email"`
	WorkspaceName string `json:"workspace_name"`
	Role          string `json:"role"`
	Token         string `json:"token"`
	InvitationURL string `json:"invitation_url"`
}

type PasswordResetPayload struct {
	ToEmail  string `json:"to_email"`
	Token    string `json:"token"`
	ResetURL string `json:"reset_url"`
}

type CustomEmailPayload struct {
	ToEmail  string        `json:"to_email"`
	Template EmailTemplate `json:"template"`
}


// Decode methods to extract specific payloads
func (m *EmailMessage) DecodeMagicLink() (*MagicLinkPayload, error) {
	if m.Type != MessageTypeMagicLink {
		return nil, fmt.Errorf("expected message type %s, got %s", MessageTypeMagicLink, m.Type)
	}

	var payload MagicLinkPayload
	if err := json.Unmarshal(m.Payload, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal magic link payload: %w", err)
	}

	return &payload, nil
}

func (m *EmailMessage) DecodeInvitation() (*InvitationPayload, error) {
	if m.Type != MessageTypeInvitation {
		return nil, fmt.Errorf("expected message type %s, got %s", MessageTypeInvitation, m.Type)
	}

	var payload InvitationPayload
	if err := json.Unmarshal(m.Payload, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal invitation payload: %w", err)
	}

	return &payload, nil
}

func (m *EmailMessage) DecodePasswordReset() (*PasswordResetPayload, error) {
	if m.Type != MessageTypePasswordReset {
		return nil, fmt.Errorf("expected message type %s, got %s", MessageTypePasswordReset, m.Type)
	}

	var payload PasswordResetPayload
	if err := json.Unmarshal(m.Payload, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal password reset payload: %w", err)
	}

	return &payload, nil
}

func (m *EmailMessage) DecodeCustomEmail() (*CustomEmailPayload, error) {
	if m.Type != MessageTypeCustom {
		return nil, fmt.Errorf("expected message type %s, got %s", MessageTypeCustom, m.Type)
	}

	var payload CustomEmailPayload
	if err := json.Unmarshal(m.Payload, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal custom email payload: %w", err)
	}

	return &payload, nil
}
