package emailservice

import "fmt"

type EmailConsumer struct {
	emailService *EmailService
}

func NewEmailConsumer(emailService *EmailService) *EmailConsumer {
	return &EmailConsumer{
		emailService: emailService,
	}

}

// Example consumer handler
func (c *EmailConsumer) HandleEmailMessage(msg EmailMessage) error {
	switch msg.Type {
	case MessageTypeMagicLink:
		payload, err := msg.DecodeMagicLink()
		if err != nil {
			return err
		}
		return c.emailService.SendMagicLink(payload.ToEmail, payload.Token, payload.VerifyURL)

	case MessageTypeInvitation:
		payload, err := msg.DecodeInvitation()
		if err != nil {
			return err
		}
		return c.emailService.SendInvitationLink(
			payload.ToEmail,
			payload.WorkspaceName,
			payload.Role,
			payload.Token,
			payload.InvitationURL,
		)

	case MessageTypePasswordReset:
		payload, err := msg.DecodePasswordReset()
		if err != nil {
			return err
		}
		return c.emailService.SendPasswordResetLink(payload.ToEmail, payload.Token, payload.ResetURL)

	case MessageTypeCustom:
		payload, err := msg.DecodeCustomEmail()
		if err != nil {
			return err
		}
		return c.emailService.SendCustomEmail(payload.ToEmail, payload.Template)

	default:
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}
}


