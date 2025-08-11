package notifier

import (
	"context"
	"fmt"

	"github.com/suhas-developer07/notification-service/internal/models"
)

type Notifier interface {
	Send(ctx context.Context, n models.Notification) error
}

func GetNotification(channel string) Notifier {
	switch channel {
	case "email":
		return &EmailNotifier{}
	case "sms":
		return &SMSNotifier{}
	case "push":
		return &PushNotifier{}
	default:
		fmt.Println("Unknown channel:", channel)
		return nil
	}
}
