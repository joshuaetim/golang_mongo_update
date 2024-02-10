package main

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func insertBook(db *mongo.Database) {
	book := Book{
		Title:    "The Alchemist",
		Author:   "Paul Coelho",
		Price:    25.00,
		Category: "fiction",
		Formats:  []string{"paper", "pdf", "epub"},
		Count:    35,
	}

	collection := db.Collection("books")
	inserted, err := collection.InsertOne(context.TODO(), &book)
	if err != nil {
		log.Fatalf("error inserting book: %v", err)
	}

	fmt.Printf("book inserted: %v\n", inserted.InsertedID)
}

func findOneBook(db *mongo.Database) {
	title := "The Alchemist"
	filter := bson.D{{"title", title}}

	collection := db.Collection("books")

	var book Book
	err := collection.FindOne(context.TODO(), filter).Decode(&book)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Printf("no document found with title: %s\n", title)
			return
		}
		log.Fatalf("error retrieving book: %v", err)
	}

	fmt.Println(book)
}

func getAllBooks(db *mongo.Database) {
	collection := db.Collection("books")

	var books []Book
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatalf("error fetching all books: %v", err)
	}
	if err := cursor.All(context.TODO(), &books); err != nil {
		log.Fatalf("error decoding cursor data: %v", err)
	}

	for _, book := range books {
		fmt.Println(book.Title)
	}
}

func deleteBook(db *mongo.Database) {
	collection := db.Collection("books")

	filter := bson.D{{"title", "The Alchemist"}}

	res, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatalf("error deleting book: %v", err)
	}
	fmt.Println(res.DeletedCount, "item deleted")
}
