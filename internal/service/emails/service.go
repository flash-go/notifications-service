package service

import (
	"context"
	"time"

	"github.com/flash-go/notifications-service/internal/domain/entity"
	"github.com/flash-go/notifications-service/internal/domain/factory"
	emailsRepositoryAdapterPort "github.com/flash-go/notifications-service/internal/port/adapter/repository/emails"
	emailsServicePort "github.com/flash-go/notifications-service/internal/port/service/emails"
)

type Config struct {
	EmailsRepository emailsRepositoryAdapterPort.Interface
}

func New(config *Config) emailsServicePort.Interface {
	return &service{
		config.EmailsRepository,
	}
}

type service struct {
	emailsRepository emailsRepositoryAdapterPort.Interface
}

func (s *service) Send(ctx context.Context, data *emailsServicePort.SendData) (*entity.Email, error) {
	// Create email entity
	email := factory.NewEmail(
		factory.EmailData{
			FromEmail: data.FromEmail,
			FromName:  data.FromName,
			Subject:   data.Subject,
			ToEmail:   data.ToEmail,
			Html:      data.Html,
			Text:      data.Text,
			Created:   time.Now(),
		},
	)

	// Send email
	if err := s.emailsRepository.Send(ctx, email); err != nil {
		return nil, err
	}

	return email, nil
}
