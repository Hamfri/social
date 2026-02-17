package mailer

import (
	"bytes"
	"time"

	"github.com/wneessen/go-mail"

	ht "html/template"
	tt "text/template"
)

type SMTPMailer struct {
	client *mail.Client
	sender string
}

func NewMailtrap(host string, port int, username, password, sender string) (*SMTPMailer, error) {
	client, err := mail.NewClient(
		host,
		mail.WithSMTPAuth(mail.SMTPAuthLogin),
		mail.WithPort(port),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}

	return &SMTPMailer{
		client: client,
		sender: sender,
	}, nil
}

func (m *SMTPMailer) Send(recipient, templateFile string, data any) error {
	textTmpl, err := tt.New("").ParseFS(tempalateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = textTmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = textTmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlTmpl, err := ht.New("").ParseFS(tempalateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = htmlTmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMsg()
	err = msg.To(recipient)
	if err != nil {
		return err
	}

	if err = msg.From(m.sender); err != nil {
		return err
	}

	msg.Subject(subject.String())
	msg.SetBodyString(mail.TypeTextPlain, plainBody.String())
	msg.AddAlternativeString(mail.TypeTextHTML, htmlBody.String())

	for i := 1; i <= MaxRetries; i++ {
		if err = m.client.DialAndSend(msg); err == nil {
			return nil
		}

		// Exponential backoff: 500ms, 1000ms, 2000ms
		time.Sleep(time.Millisecond * time.Duration(i*500))
	}

	return err
}
