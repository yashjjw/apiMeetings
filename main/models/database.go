package models

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDatabase(mongoURL string) *mongo.Client {
	clientOptions := options.Client().ApplyURI(mongoURL)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database successfully")
	return client
}

type Meeting struct {
	ID              string        `bson:"id" json:"id"`
	Title           string        `bson:"title" json:"title"`
	Participants    []Participant `bson:"participants" json:"participants"`
	StartTime       int64         `bson:"startTime" json:"startTime"`
	EndTime         int64         `bson:"endTime" json:"endTime"`
	CreationTime    int64         `bson:"creationTime" json:"creationTime"`
}

type Participant struct {
	Name  string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	RSVP  string `bson:"rsvp" json:"rsvp"`
}

