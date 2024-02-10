package main

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"reflect"
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

func updateBook(db *mongo.Database) {
	collection := db.Collection("books")

	// create updates variable to hold all the update fields
	updates := bson.D{}

	newBook := Book{
		Author: "Paulo Coelho",
	}
	// get the type of struct == Book
	typeData := reflect.TypeOf(newBook)

	// get the values from the provided object: author -> Paulo Coelho
	values := reflect.ValueOf(newBook)

	// starting from index 1 to exclude the ID field
	for i := 1; i < typeData.NumField(); i++ {
		field := typeData.Field(i)   // get the field from the struct definition
		val := values.Field(i)       // get the value from the specified field position
		tag := field.Tag.Get("json") // from the field, get the json struct tag

		// we want to avoid zero values, as the omitted fields from newBook
		// corresponds to their zero values, and we only want provided fields
		if !isZeroType(val) {
			update := bson.E{Key: tag, Value: val.Interface()}
			updates = append(updates, update)
		}
	}

	filter := bson.D{{"title", "The Alchemist"}}
	updateFilter := bson.D{{"$set", updates}}
	_, err := collection.UpdateOne(context.TODO(), filter, updateFilter)
	if err != nil {
		log.Fatalf("error updating book: %v", err)
	}
}

// isZeroType checks if the value from the struct is the zero value of its type
func isZeroType(value reflect.Value) bool {
	zero := reflect.Zero(value.Type()).Interface()

	switch value.Kind() {
	case reflect.Slice, reflect.Array, reflect.Chan, reflect.Map:
		return value.Len() == 0
	default:
		return reflect.DeepEqual(zero, value.Interface())
	}
}
