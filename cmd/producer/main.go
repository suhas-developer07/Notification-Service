package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/suhas-developer07/notification-service/internal/models"
	"github.com/suhas-developer07/notification-service/internal/queue"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Env file is not loaded")
	}
	ctx := context.Background()

	sqsClient, err := queue.NewSQSClient(ctx, "http://localhost:4566/000000000000/notify")

	if err != nil {
		log.Fatal(err)
	}

	notifcation := models.Notification{
		ID:      uuid.New().String(),
		Channel: "email",
		Message: "Hello Dear",
		Target:  "suhasdeveloper07@gmail.com",
	}

	if err := sqsClient.SendMessage(ctx, notifcation); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Message sent:", notifcation)
}
