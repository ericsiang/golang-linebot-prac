package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func ConnectMongoDb(dsn string) (*mongo.Client, context.Context, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))

	if err != nil {
		return nil, ctx, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, ctx, err
	}
	log.Println("Connected to MongoDB!")
	return client, ctx, err
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("linebot").Collection(collectionName)
	return collection
}
