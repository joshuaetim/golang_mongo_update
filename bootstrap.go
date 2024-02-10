package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// Mongo initializes a mongo db connection.
func Mongo() *mongo.Database {
	uri := "mongodb://localhost:27017"

	bsonOpts := &options.BSONOptions{
		UseJSONStructTags: false,
	}
	clientOpts := options.Client().
		ApplyURI(uri).
		SetBSONOptions(bsonOpts)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal(err)
	}
	db := client.Database("book_store")
	return db
}
