package notifier

import (
	"context"
	"fmt"

	"github.com/suhas-developer07/notification-service/internal/models"
)

type PushNotifier struct {
}

func (p *PushNotifier) Send(ctx context.Context, n models.Notification) error {
	fmt.Printf("Sending Push Notification to %s : %s\n", n.Recipient, n.Message)
	return nil
}
