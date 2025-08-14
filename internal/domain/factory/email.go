package factory

import (
	"time"

	"github.com/flash-go/notifications-service/internal/domain/entity"
)

func NewEmail(data EmailData) *entity.Email {
	return &entity.Email{
		FromEmail: data.FromEmail,
		FromName:  data.FromName,
		Subject:   data.Subject,
		ToEmail:   data.ToEmail,
		Html:      data.Html,
		Text:      data.Text,
		Created:   time.Unix(0, data.Created.UnixNano()),
	}
}

type EmailData struct {
	FromEmail string
	FromName  string
	Subject   string
	ToEmail   string
	Html      string
	Text      string
	Created   time.Time
}
