/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"homework/command"
	"homework/models"
	"homework/mongodb"
	"log"
	"net/http"
	"time"
)

var config *viper.Viper

func main() {
	command.DockerUp()
	config := initConfigure()
	bot, err := linebot.New(fmt.Sprintf("%s", config.Get("linebot.channelSecret")), fmt.Sprintf("%s", config.Get("linebot.channelAccessToken")))
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.POST("/callback", func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				fmt.Println("ErrInvalidSignature")
				log.Print(err)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				dsn := fmt.Sprintf("mongodb://%s:%s@%s:%d", config.Get("database.user"), config.Get("database.password"), config.Get("database.host"), config.Get("database.port"))
				client, ctx, err := mongodb.ConnectMongoDb(dsn)

				if err != nil {
					log.Fatal(err)
				}
				messageType := ""
				messageText := ""
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					messageType = "text"
					messageText = message.Text
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("已收到文字訊息！")).Do(); err != nil {
						log.Print(err)
					}
				case *linebot.ImageMessage:
					messageType = "image"
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("已收到圖片訊息！")).Do(); err != nil {
						log.Print(err)
					}
				case *linebot.VideoMessage:
					messageType = "video"
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("已收到影片訊息！")).Do(); err != nil {
						log.Print(err)
					}
				case *linebot.AudioMessage:
					messageType = "audio"
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("已收到語音訊息！")).Do(); err != nil {
						log.Print(err)
					}
				case *linebot.FileMessage:
					messageType = "file"
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("已收到檔案訊息！")).Do(); err != nil {
						log.Print(err)
					}
				case *linebot.LocationMessage:
					messageType = "location"
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("已收到定位訊息！")).Do(); err != nil {
						log.Print(err)
					}
				case *linebot.StickerMessage:
					messageType = "sticker"
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("已收到貼圖訊息！")).Do(); err != nil {
						log.Print(err)
					}
				default:
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("已收到訊息！")).Do(); err != nil {
						log.Print(err)
					}
				}

				userCollection := mongodb.GetCollection(client, "users")
				userCount, err := userCollection.CountDocuments(ctx, bson.M{"userid": event.Source.UserID})
				if err != nil {
					log.Fatal(userCount)
				}
				if userCount == 0 {
					user := models.User{
						Id:        primitive.NewObjectID(),
						UserId:    event.Source.UserID,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}
					result, err := userCollection.InsertOne(ctx, user)
					if err != nil {
						log.Print(err)
					}
					fmt.Println(result)
				}

				lineMessageCollection := mongodb.GetCollection(client, "lineMessage")
				lineMessage := models.LineMessage{
					Id:             primitive.NewObjectID(),
					EventType:      "message",
					MessageType:    messageType,
					MessageText:    messageText,
					UserId:         event.Source.UserID,
					ReplyToken:     event.ReplyToken,
					WebhookEventId: event.WebhookEventID,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}
				result, err := lineMessageCollection.InsertOne(ctx, lineMessage)
				if err != nil {
					log.Print(err)
				}
				fmt.Println(result)
			}
		}

	})

	v1 := router.Group("/api/v1")
	v1.POST("/sendMessage/:userId", func(c *gin.Context) {
		userId := c.Param("userId")
		pushMessage := c.DefaultQuery("message", "send empty message")
		message := linebot.NewTextMessage(pushMessage)

		_, err := bot.PushMessage(userId, message).Do()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			log.Print(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "send message successfully!"})
	})

	router.Run(":80")
}

func initConfigure() *viper.Viper {
	config := viper.New()
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	config.AddConfigPath("./config")
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	config.WatchConfig()
	config.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	return config
}
