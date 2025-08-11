package notifier

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/suhas-developer07/notification-service/internal/config"
	"github.com/suhas-developer07/notification-service/internal/models"
)

type EmailNotifier struct {
}

func (e *EmailNotifier) Send(ctx context.Context, n models.Notification) error {

	cfg := config.LoadConfig()
	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)

	msg := []byte(fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\n\r\n%s\r\n",
		n.Target, n.Message, n.Message,
	))

	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)
	return smtp.SendMail(addr, auth, cfg.SMTPUser, []string{n.Target}, msg)
}
