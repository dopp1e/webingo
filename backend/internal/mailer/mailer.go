package mailer

import "embed"

const (
	FromName               = "Webingo"
	maxRetries             = 3
	UserInvitationTemplate = "user_invitation.html"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) error
}
