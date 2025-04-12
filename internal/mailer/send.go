package mailer

import (
	"context"
	"os"

	"github.com/mailersend/mailersend-go"
)

var ms *mailersend.Mailersend

func InitMailer() {
	ms = mailersend.NewMailersend(os.Getenv("MAILER_SEND_API_TOKEN"))
}

func SendEmail(fromEmail, fromName, toEmail, toName, subject, body string) error {
	message := ms.Email.NewMessage()

	message.SetFrom(mailersend.Recipient{
		Email: fromEmail,
		Name:  fromName,
	})
	message.SetRecipients([]mailersend.Recipient{
		{
			Email: toEmail,
			Name:  toName,
		},
	})
	message.SetSubject(subject)
	message.SetHTML(body)

	_, err := ms.Email.Send(context.Background(), message)
	return err
}
