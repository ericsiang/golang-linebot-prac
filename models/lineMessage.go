package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LineMessage struct {
	Id             primitive.ObjectID
	EventType      string
	MessageType    string
	MessageText    string
	UserId         string
	ReplyToken     string
	WebhookEventId string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
