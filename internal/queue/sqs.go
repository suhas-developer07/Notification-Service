package queue

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/go-delve/delve/pkg/config"
)

type SQSClient struct {
	client   *sqs.Client
	queueURL string
}

func NewSQSClient(ctx context.Context, queueURL string) (*SQSClient, error) {
	cfg, err := config.ConfigDefaoult(ctx)
	if err != nil {
		return nil, err
	}
	return &SQSClient{
		client:   sqs.NewFromConfig(cfg),
		queueURL: queueURL,
	}, nil
}
