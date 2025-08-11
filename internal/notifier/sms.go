package notifier

import (
	"context"
	"fmt"

	"github.com/suhas-developer07/notification-service/internal/models"
)

type SMSNotifier struct {
}

func (s *SMSNotifier) Send(ctx context.Context, n models.Notification) error {
	fmt.Printf("Sending SMS to %s : %s\n", n.Target, n.Message)
	return nil
}
