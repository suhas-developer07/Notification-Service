package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/suhas-developer07/notification-service/internal/models"
)

type DynamoDBClient struct {
	Client *dynamodb.Client
}

func NewDynamoDBClient() (*DynamoDBClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("ap-southeast-1"), // Can be any valid AWS region
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				if service == dynamodb.ServiceID {
					return aws.Endpoint{
						URL:           "http://localhost:8000", // Local DynamoDB Docker
						SigningRegion: "ap-southeast-1",
					}, nil
				}
				return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
			})),
	)
	if err != nil {
		return nil, fmt.Errorf("error loading AWS config: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	return &DynamoDBClient{Client: client}, nil
}

func (d *DynamoDBClient) EnsureTableExists(ctx context.Context, tableName string) error {

	_, err := d.Client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err == nil {
		log.Println("Table was already created")
		return nil
	}

	log.Println("Creating Table..")

	_, err = d.Client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("notification_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("notification_id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})

	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	waiter := dynamodb.NewTableExistsWaiter(d.Client)

	err = waiter.Wait(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}, 5*time.Minute)

	if err != nil {
		return fmt.Errorf("Error waiting for table to become active : %v", err)
	}
	log.Println("Table Created successfully..")
	return nil
}

func (d *DynamoDBClient) SaveNotificationStatus(ctx context.Context, notification models.Notification, tableName string) error {

	if notification.NotificationID == "" {
		notification.NotificationID = uuid.NewString()
	}

	fmt.Println(notification.NotificationID)
	now := time.Now().UTC().Format(time.RFC3339)

	if notification.CreatedAt == "" {
		notification.CreatedAt = now
	}

	notification.UpdatedAt = now

	av, err := attributevalue.MarshalMap(notification)

	if err != nil {
		return fmt.Errorf("failed to marshal notification :%w", err)
	}

	_, err = d.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})

	if err != nil {
		return fmt.Errorf("failed to save notification :%w", err)
	}

	return nil
}

func (db *DynamoDBClient) UpdateNotificationStatus(ctx context.Context, id string, status string, tableName string) error {
	now := time.Now().UTC().Format(time.RFC3339)

	_, err := db.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
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
	})
	if err != nil {
		return fmt.Errorf("failed to update notification status :%w", err)
	}
	return nil
}
