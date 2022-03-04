package main

import (
	"context"
	"fmt"
	"math/rand"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	client struct {
		mongoClient *mongo.Client
		rand        *rand.Rand

		dbName    string
		colNumber uint64
	}
)

const colNamePrefix = "bench"

func (c *client) createCols(ctx context.Context) error {
	for i := 0; i < int(c.colNumber); i++ {
		if err := c.mongoClient.Database(c.dbName).CreateCollection(ctx, fmt.Sprintf("%s-%d", colNamePrefix, i)); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) insert(ctx context.Context) error {
	colName := fmt.Sprintf("%s-%d", colNamePrefix, c.rand.Uint64()%c.colNumber)
	_, err := c.mongoClient.Database(c.dbName).Collection(colName).InsertOne(ctx, bson.M{"val": rand.Int()})
	if err != nil {
		return err
	}
	return nil
}
