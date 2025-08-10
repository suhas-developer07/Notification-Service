package models

type Notification struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Message string `json:"message"`
}
