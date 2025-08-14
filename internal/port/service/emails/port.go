package port

import (
	"context"

	"github.com/flash-go/notifications-service/internal/domain/entity"
)

type Interface interface {
	Send(ctx context.Context, data *SendData) (*entity.Email, error)
}

// Args

type SendData struct {
	FromEmail string
	FromName  string
	Subject   string
	ToEmail   string
	Html      string
	Text      string
}
