package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Title    string             `json:"title"`
	Author   string             `json:"author"`
	Price    float64            `json:"price"`
	Category string             `json:"category"`
	Formats  []string           `json:"formats"`
	Count    int                `json:"count"` // number of books left
}
