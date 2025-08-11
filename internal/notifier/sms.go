package notifier

import (
	"context"
	"fmt"
	"os"

	"github.com/suhas-developer07/notification-service/internal/models"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMSNotifier struct {
}

func (s *SMSNotifier) Send(ctx context.Context, n models.Notification) error {

	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	from := os.Getenv("TWILIO_FROM")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo("+18777804236")
	params.SetFrom(from)
	params.SetBody("Hello prasanna")

	_, err := client.Api.CreateMessage(params)

	if err != nil {
		return fmt.Errorf("failed to send SMS:%w", err)
	}

	return nil
}
