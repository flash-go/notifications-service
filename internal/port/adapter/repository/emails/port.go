package port

import (
	"context"

	"github.com/flash-go/notifications-service/internal/domain/entity"
)

type Interface interface {
	Send(ctx context.Context, email *entity.Email) error
}
