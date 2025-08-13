package models

type Notification struct {
	NotificationID string `dynamodbav:"notification_id"`
	Recipient      string `dynamodbav:"recipient"`
	Channel        string `dynamodbav:"channel"`
	Message        string `dynamodbav:"message"`
	Status         string `dynamodbav:"status"`
	CreatedAt      string `dynamodbav:"created_at"`
	UpdatedAt      string `dynamodbav:"updated_at"`
}

type NotificationPayload struct {
	Recipient string `json:"recipient"`
	Channel   string `json:"channel"`
	Message   string `json:"message"`
}
