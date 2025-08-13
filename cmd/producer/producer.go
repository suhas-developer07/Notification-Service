package producer

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/suhas-developer07/notification-service/internal/models"
	"github.com/suhas-developer07/notification-service/internal/queue"
)

func ProduceMessage(w http.ResponseWriter, r *http.Request) {
	var payload models.NotificationPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "failed to decode Body", http.StatusBadRequest)
		return
	}
	ctx := context.Background()

	sqsClient, err := queue.NewSQSClient(ctx, "http://localhost:4566/000000000000/notification")

	if err != nil {
		http.Error(w, "failed to connect sqs", http.StatusInternalServerError)
		return
	}

	notifcation := models.Notification{
		NotificationID: uuid.New().String(),
		Channel:        payload.Channel,
		Message:        payload.Message,
		Recipient:      payload.Recipient,
	}

	if err := sqsClient.SendMessage(ctx, notifcation); err != nil {
		http.Error(w, "failed to send message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Message Sent Successfully",
	})
}
