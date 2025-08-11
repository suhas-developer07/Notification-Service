package queue

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSClient struct {
	client   *sqs.Client
	queueURL string
}

func NewSQSClient(ctx context.Context, queueURL string) (*SQSClient, error) {
	
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           "http://localhost:4566", // LocalStack endpoint
			SigningRegion: "ap-southeast-1",
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("ap-southeast-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
		config.WithEndpointResolverWithOptions(customResolver),
	)

	if err != nil {
		return nil, err
	}

	return &SQSClient{
		client:   sqs.NewFromConfig(cfg),
		queueURL: queueURL,
	}, nil
}
func (q *SQSClient) SendMessage(ctx context.Context, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = q.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &q.queueURL,
		MessageBody: aws.String(string(body)),
	})
	return err
}

func (q *SQSClient) ReceiveMessage(ctx context.Context, max int32) ([]types.Message, error) {
	resp, err := q.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &q.queueURL,
		MaxNumberOfMessages: max,
		WaitTimeSeconds:     2,
	})

	if err != nil {
		return nil, err
	}

	return resp.Messages, nil
}

func (q *SQSClient) DeleteMessage(ctx context.Context, receiptHandle string) error {
	_, err := q.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &q.queueURL,
		ReceiptHandle: &receiptHandle,
	})
	return err
}
