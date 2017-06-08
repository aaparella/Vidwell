package notifier

import "github.com/aaparella/vidwell/models"

type EmailSource interface {
	Subject() string
	Content() string
}

type EmailNotifier struct {
	source EmailSource
}

type Email struct {
	Subject, Content string
}

func (en *EmailNotifier) Notify(user *models.User) error {
	email := en.generateEmail()
	return en.deliverEmail(user.Email, email)
}

func (en *EmailNotifier) generateEmail() *Email {
	return &Email{
		Subject: en.source.Subject(),
		Content: en.source.Content(),
	}
}

func (en *EmailNotifier) deliverEmail(address string, email *Email) error {
	// Deliver the actual email
	return nil
}
