package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type sessionMongo struct {
	client     *mongo.Client
	database   string
	timeout    time.Duration
	collection string
}

func NewSessionMongo(sessionStoreName string) *sessionMongo {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if err != nil {
		panic(err)
	}
	timeout := 20 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	return &sessionMongo{
		client:     client,
		database:   sessionStoreName,
		timeout:    timeout,
		collection: "userSession",
	}
}

func (sM *sessionMongo) Store(id string, key string, value interface{}) error {
	return nil
}
func (sM *sessionMongo) Find(id string, key string) (interface{}, error) {
	return nil, nil
}

func (sM *sessionMongo) IsSessionAvailable(id string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), sM.timeout)
	defer cancel()

	lstCollections, _ := sM.client.Database(sM.database).ListCollectionNames(ctx, bson.D{})

	for _, col := range lstCollections {
		if col == id {
			return true
		}
	}
	return false
}
