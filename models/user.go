package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id        primitive.ObjectID
	UserId    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
