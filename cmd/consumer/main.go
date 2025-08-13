package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/suhas-developer07/notification-service/internal/db"
	"github.com/suhas-developer07/notification-service/internal/models"
	"github.com/suhas-developer07/notification-service/internal/notifier"
	"github.com/suhas-developer07/notification-service/internal/queue"
)

func handleMessage(ctx context.Context, sqsClient *queue.SQSClient, dbClient *db.DynamoDBClient, msg types.Message) error {
	var notif models.Notification
	if err := json.Unmarshal([]byte(*msg.Body), &notif); err != nil {
		logrus.WithError(err).Error("Error unmarshalling message")
		_ = sqsClient.DeleteMessage(ctx, *msg.ReceiptHandle)
		return err
	}

	// Save as Pending
	if err := dbClient.SaveNotificationStatus(ctx, notif, "Notification"); err != nil {
		logrus.WithError(err).Error("Error saving notification")
	}

	// Get correct notifier
	n := notifier.GetNotification(notif.Channel)
	if n == nil {
		logrus.WithField("channel", notif.Channel).Error("No notifier found")
		_ = dbClient.UpdateNotificationStatus(ctx, notif.NotificationID, "failed", "Notification")

		// Delete message so it doesn't loop forever
		_ = sqsClient.DeleteMessage(ctx, *msg.ReceiptHandle)
		return fmt.Errorf("unknown channel: %s", notif.Channel)
	}

	// Send notification
	if err := n.Send(ctx, notif); err != nil {
		logrus.WithError(err).Error("Error sending notification")
		_ = dbClient.UpdateNotificationStatus(ctx, notif.NotificationID, "failed", "Notification")

		// Could send to DLQ here
		_ = sqsClient.DeleteMessage(ctx, *msg.ReceiptHandle)
		return err
	}

	// Update status as sent
	_ = dbClient.UpdateNotificationStatus(ctx, notif.NotificationID, "sent", "Notification")

	// Remove from SQS
	if err := sqsClient.DeleteMessage(ctx, *msg.ReceiptHandle); err != nil {
		return fmt.Errorf("failed to delete SQS message: %w", err)
	}

	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No env file found", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// SQS setup
	sqsClient, err := queue.NewSQSClient(ctx, "http://localhost:4566/000000000000/notification")
	if err != nil {
		log.Fatal(err)
	}

	// DynamoDB setup
	dynamoClient, err := db.NewDynamoDBClient()
	if err != nil {
		log.Fatalf("Error creating DynamoDB client: %v", err)
	}
	if err = dynamoClient.EnsureTableExists(ctx, "Notification"); err != nil {
		log.Fatalf("Error ensuring table exists: %v", err)
	}

	if err != nil {
		log.Fatal("Error loading AWS config:", err)
	}

	// Handle shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nShutting down worker...")
		cancel()
	}()

	fmt.Println("Worker started. Listening for messages...")

	// Main loop
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msgs, err := sqsClient.ReceiveMessage(ctx, 5)
			if err != nil {
				logrus.WithError(err).Error("Error receiving messages from SQS")
				continue
			}

			if len(msgs) == 0 {
				time.Sleep(1 * time.Second)
				continue
			}

			for _, m := range msgs {
				handleMessage(ctx, sqsClient, dynamoClient, m)
			}
		}
	}
}
