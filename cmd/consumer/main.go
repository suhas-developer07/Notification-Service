package main

import (
	"context"
	"fmt"
	"log"

	"github.com/suhas-developer07/notification-service/internal/queue"
)

func main() {
	ctx := context.Background()

	sqsClient, err := queue.NewSQSClient(ctx, "http://localhost:4566/000000000000/notify")

	if err != nil {
		log.Fatal(err)
	}

	msgs, err := sqsClient.ReceiveMessage(ctx, 5)

	if err != nil {
		log.Fatal(err)
	}

	for _, msg := range msgs {
		fmt.Println("Recieved:", msg)
	}
}
