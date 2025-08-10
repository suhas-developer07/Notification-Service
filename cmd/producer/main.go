package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/suhas-developer07/notification-service/internal/models"
	"github.com/suhas-developer07/notification-service/internal/queue"
)

func main() {
	ctx := context.Background()

	sqsClient, err := queue.NewSQSClient(ctx, "http://localhost:4566/000000000000/notify")

	if err != nil {
		log.Fatal(err)
	}

	notifcation := models.Notification{
		ID:      uuid.New().String(),
		Type:    "email",
		Message: fmt.Sprintf("Test Message at %v", time.Now()),
	}

	if err := sqsClient.SendMessage(ctx, notifcation); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Message sent:", notifcation)
}
