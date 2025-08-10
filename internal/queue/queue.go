package queue

import (
	"context"

	"github.com/suhas-developer07/notification-service/internal/models"
)

type NotificationMessage interface {
	SendMessage(ctx context.Context, notification models.Notification) error
	ReceiveMessage(ctx context.Context, max int32) ([]NotificationMessage, error)
	DeleteMessage(ctx context.Context, reciptHandle string) error
}
