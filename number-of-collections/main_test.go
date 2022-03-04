package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var c *client

func TestMain(m *testing.M) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
		return
	}

	dbName := fmt.Sprintf("benchmark-%d", time.Now().UnixNano())

	colNumEnv, ok := os.LookupEnv("COL_NUM")
	if !ok {
		colNumEnv = "1"
	}
	colNum, err := strconv.ParseUint(colNumEnv, 10, 64)
	if err != nil {
		colNum = 1
	}
	log.Printf("colNum: %d", colNum)

	c = &client{
		mongoClient: mongoClient,
		rand:        rnd,
		dbName:      dbName,
		colNumber:   colNum,
	}

	defer func() {
		_ = mongoClient.Database(dbName).Drop(context.Background())
	}()

	if err := c.createCols(context.Background()); err != nil {
		panic(err)
	}

	m.Run()
}

func Benchmark_client_insert(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		if err := c.insert(ctx); err != nil {
			b.Error(err)
		}
	}
}
