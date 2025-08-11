package models

type Notification struct {
	ID      string `json:"id"`
	Channel string `json:"type"`
	Message string `json:"message"`
	Target  string `json:"target"`
}
