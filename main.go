/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"homework/command"
	"homework/mongodb"
	"log"
)

var config *viper.Viper

func main() {
	command.DockerUp()
	config := initConfigure()
	dsn := fmt.Sprintf("mongodb://%s:%s@%s:%d", config.Get("database.user"), config.Get("database.password"), config.Get("database.host"), config.Get("database.port"))
	client, ctx, err := mongodb.ConnectMongoDb(dsn)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(client, ctx)
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
