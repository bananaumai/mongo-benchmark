package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	TestDoc struct {
		ID int64 `bson:"_id"`
	}
)

var (
	dbName = "testing"
	client *mongo.Client
)

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	client = c
}

func main() {
	ctx := context.Background()

	numOfDocs := int64(10000)
	if err := insertData(ctx, numOfDocs); err != nil {
		panic(err)
	}

	type getFunc func(context.Context, int64) (int64, int64, error)

	for _, f := range []getFunc{getFirstLastByAggregation, getFirstLastBySimpleQueries} {
		first, last, err := f(ctx, numOfDocs)
		if err != nil {
			panic(err)
		}
		fmt.Printf("first: %d, lst:%d\n", first, last)
	}

	if err := client.Database(dbName).Drop(ctx); err != nil {
		panic(err)
	}
}

func insertData(ctx context.Context, numOfDocs int64) error {
	docs := make([]interface{}, numOfDocs)
	for i := int64(0); i < numOfDocs; i++ {
		docs[i] = TestDoc{i}
	}

	_, err := client.Database(dbName).Collection(getCollectionName(numOfDocs)).InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	return nil
}

type aggregationResult struct {
	First int64 `bson:"first"`
	Last  int64 `bson:"last"`
}

func getFirstLastByAggregation(ctx context.Context, targetNumOfDocs int64) (int64, int64, error) {
	sortStage := bson.D{{"$sort", bson.D{{"_id", 1}}}}

	groupStage := bson.D{{"$group", bson.D{
		{"_id", 1},
		{"first", bson.D{{"$first", "$_id"}}},
		{"last", bson.D{{"$last", "$_id"}}},
	}}}

	collection := client.Database(dbName).Collection(getCollectionName(targetNumOfDocs))
	cur, err := collection.Aggregate(ctx, mongo.Pipeline{sortStage, groupStage})
	if err != nil {
		return 0, 0, nil
	}

	var res []aggregationResult
	if err := cur.All(ctx, &res); err != nil {
		return 0, 0, err
	}

	return res[0].First, res[0].Last, nil
}

func getFirstLastBySimpleQueries(ctx context.Context, targetNumOfDocs int64) (int64, int64, error) {
	var (
		doc   TestDoc
		first int64
		last  int64
	)

	collection := client.Database(dbName).Collection(getCollectionName(targetNumOfDocs))

	res := collection.FindOne(ctx, bson.D{}, &options.FindOneOptions{Sort: bson.M{"_id": 1}})
	if err := res.Decode(&doc); err != nil {
		return 0, 0, err
	}
	first = doc.ID

	res = collection.FindOne(ctx, bson.D{}, &options.FindOneOptions{Sort: bson.M{"_id": -1}})
	if err := res.Decode(&doc); err != nil {
		return 0, 0, err
	}
	last = doc.ID

	return first, last, nil
}

func getCollectionName(numOfDocs int64) string {
	return fmt.Sprintf("data-%d", numOfDocs)
}
