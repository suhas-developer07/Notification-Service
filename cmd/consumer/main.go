package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/suhas-developer07/notification-service/internal/models"
	"github.com/suhas-developer07/notification-service/internal/queue"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sqsClient, err := queue.NewSQSClient(ctx, "http://localhost:4566/000000000000/notify")

	if err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down worker...")
		cancel()
	}()

	fmt.Println("Worker started. Listening for messages...")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msgs, err := sqsClient.ReceiveMessage(ctx, 5)
			if err != nil {
				log.Println("Error receiving messages:", err)
				continue
			}
			for _, m := range msgs {
				var notif models.Notification
				if err := json.Unmarshal([]byte(*m.Body), &notif); err != nil {
					log.Println("Error unmarshalling:", err)
					continue
				}
				log.Printf("Processing: %+v", notif)

				if err := sqsClient.DeleteMessage(ctx, *m.ReceiptHandle); err != nil {
					log.Println("Error deleting message:", err)
				}
			}
		}
	}
}
