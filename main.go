package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Models *Models
}

func main() {
	mgClient, err := connectToDB()

	if err != nil {
		log.Panic("Unable to connect to mongo")
	}
	app := Config{
		Models: New(mgClient),
	}

	go app.startGRPC()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	defer func() {
		if err = mgClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	forever := make(chan int)

	<-forever
}

const (
	mongoURL = "mongodb://mongo:27017"
)

func connectToDB() (*mongo.Client, error) {
	cOptions := options.Client().ApplyURI(mongoURL)

	cOptions.SetAuth(options.Credential{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	conn, err := mongo.Connect(context.TODO(), cOptions)

	if err != nil {
		log.Println("Error connecting mongo")
	}

	return conn, nil
}
