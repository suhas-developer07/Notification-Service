package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/joho/godotenv"

	"github.com/sirupsen/logrus"
	"github.com/suhas-developer07/notification-service/internal/models"
	"github.com/suhas-developer07/notification-service/internal/notifier"
	"github.com/suhas-developer07/notification-service/internal/queue"
)

func processNotification(ctx context.Context, notif models.Notification) error {

	logrus.Infof("Processing notification:%+v", notif)
	time.Sleep(1 * time.Second)
	return nil
}

func handleMessage(ctx context.Context, sqsClient *queue.SQSClient, msg types.Message) {
	var notif models.Notification

	if err := json.Unmarshal([]byte(*msg.Body), &notif); err != nil {
		logrus.WithError(err).Error("Error unmarshalling message")
		return
	}

	const maxRetries = 3
	var err error

	for attemt := 1; attemt < maxRetries; attemt++ {
		err = processNotification(ctx, notif)
		if err != nil {
			break
		}
		wait := time.Duration(attemt*2) * time.Second
		logrus.WithFields(logrus.Fields{
			"attemt": attemt,
			"wait":   wait,
		}).Warn("Retrying message processing")
		time.Sleep(wait)
	}

	if err != nil {
		logrus.WithError(err).Error("Fields after retries should sent to DLQ")
		return
	}

	if err := sqsClient.DeleteMessage(ctx, *msg.ReceiptHandle); err != nil {
		logrus.WithError(err).Error("Error Deleting Messages")
	}
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No env file found", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

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
				logrus.WithError(err).Error("Error Recieving Messages")
				continue
			}
			for _, m := range msgs {
				var notif models.Notification
				if err := json.Unmarshal([]byte(*m.Body), &notif); err != nil {
					log.Println("Error unmarshalling:", err)
					continue
				}
				notifier := notifier.GetNotification(notif.Channel)
				if notifier == nil {
					log.Println("no Notifier found for channel:", notif.Channel)
					continue
				}
				if err := notifier.Send(ctx, notif); err != nil {
					log.Println("Error sending notification:", err)
				}
				if len(msgs) == 0 {
					continue
				}
				var wg sync.WaitGroup
				for _, m := range msgs {
					wg.Add(1)
					go func(msg types.Message) {
						defer wg.Done()
						handleMessage(ctx, sqsClient, msg)
					}(m)
				}
				wg.Wait()
			}
		}
	}
}
