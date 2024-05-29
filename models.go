package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

type Models struct {
	LogEntry *LogEntry
}

type LogEntry struct {
	ID            string `bson:"_id,omitempty"`
	Event         string `bson:"event"`
	CallerService string `bson:"callerService"`
	Timestamp     int64  `bson:"timestamp"`
	Details       string `bson:"details"`
}

var client *mongo.Client

func New(mongo *mongo.Client) *Models {
	client = mongo
	return &Models{
		LogEntry: &LogEntry{},
	}
}

func (*LogEntry) Insert(entry *LogEntry) error {
	collection := client.Database(os.Getenv("DATABASE")).Collection(os.Getenv("COLLECTION"))

	_, err := collection.InsertOne(context.TODO(), entry)

	if err != nil {
		log.Println("Error inserting log: ", err)
		return err
	}

	return nil
}
