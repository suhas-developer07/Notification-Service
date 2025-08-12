package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/suhas-developer07/notification-service/internal/models"
)

type DynamoDBClient struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBClient(cfg aws.Config, tableName string) *DynamoDBClient {
	return &DynamoDBClient{
		client:    dynamodb.NewFromConfig(cfg),
		tableName: tableName,
	}
}

func (db *DynamoDBClient) SaveNotification(ctx context.Context, notification models.Notification) error {

	av, err := attributevalue.MarshalMap(notification)

	if err != nil {
		return fmt.Errorf("failed to marshal notification :%w", err)
	}

	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(db.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(notification_id)"),
	})

	if err != nil {
		return fmt.Errorf("failed to save notification :%w", err)
	}

	return nil
}

func (db *DynamoDBClient) UpdateNotificationStatus(ctx context.Context, id string, status string) error {
	now := time.Now().UTC().Format(time.RFC3339)

	_, err := db.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(db.tableName),
		Key: map[string]types.AttributeValue{
			"notification_id": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression: aws.String("SET #s = :status, updated_at = :updated_at"),
		ExpressionAttributeNames: map[string]string{
			"#s": "status", // 'status' is a reserved keyword in DynamoDB, so alias it
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status":     &types.AttributeValueMemberS{Value: status},
			":updated_at": &types.AttributeValueMemberS{Value: now},
		},
		ConditionExpression: aws.String("attribute_exists(notification_id)"), // Avoid updating non-existing
	})
	if err != nil {
		return fmt.Errorf("failed to update notification status :%w", err)
	}
	return nil
}
