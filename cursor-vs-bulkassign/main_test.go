package main

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var col *mongo.Collection

const numOfDocs = 100000

func TestMain(m *testing.M) {
	ctx := context.Background()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
		return
	}

	dbName := "bench-travers-vs-assign"

	db := mongoClient.Database(dbName)

	defer func() {
		_ = db.Drop(context.Background())
	}()

	data := make([]interface{}, numOfDocs)
	for i := 0; i < numOfDocs; i++ {
		data[i] = bson.M{"v": i}
	}

	col = db.Collection("bench")
	if _, err := col.InsertMany(ctx, data); err != nil {
		panic(err)
	}

	m.Run()
}

func Benchmark_travers(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		cur, err := col.Find(ctx, bson.M{})
		if err != nil {
			panic(err)
		}
		travers(ctx, cur)
	}
}

func Benchmark_assign(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		cur, err := col.Find(ctx, bson.M{})
		if err != nil {
			panic(err)
		}
		assign(ctx, cur)
	}
}
