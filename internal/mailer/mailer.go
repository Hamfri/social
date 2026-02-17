package mailer

import (
	"embed"
)

var (
	MaxRetries             = 3
	UserActivationTemplate = "user_activation.html"
)

//go:embed "templates"
var tempalateFS embed.FS

type Mailer interface {
	Send(recipient, templateFile string, data any) error
}
