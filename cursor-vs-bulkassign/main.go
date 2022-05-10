package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

func travers(ctx context.Context, cur *mongo.Cursor) {
	for cur.Next(ctx) {
		var d interface{}
		if err := cur.Decode(&d); err != nil {
			panic(err)
		}
	}
}

func assign(ctx context.Context, cur *mongo.Cursor) {
	var data []interface{}
	if err := cur.All(ctx, &data); err != nil {
		panic(err)
	}
}
