package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	dbName = fmt.Sprintf("testing-%d", time.Now().UnixNano())
	defer func() {
		_ = client.Database(dbName).Drop(ctx)
	}()

	for _, numOfDocs := range []int64{1000, 10000, 100000} {
		log.Printf("insert %d docs", numOfDocs)
		if err := insertData(ctx, numOfDocs); err != nil {
			panic(err)
		}
	}
	log.Printf("run test")

	m.Run()
}

func Benchmark_getFirstLastByAggregation_1000(b *testing.B)   { benchmarkAggregation(b, 1000) }
func Benchmark_getFirstLastByAggregation_10000(b *testing.B)  { benchmarkAggregation(b, 10000) }
func Benchmark_getFirstLastByAggregation_100000(b *testing.B) { benchmarkAggregation(b, 100000) }

func Benchmark_getFirstLastBySimpleQueries_1000(b *testing.B)   { benchmarkSimpleQueries(b, 1000) }
func Benchmark_getFirstLastBySimpleQueries_10000(b *testing.B)  { benchmarkSimpleQueries(b, 10000) }
func Benchmark_getFirstLastBySimpleQueries_100000(b *testing.B) { benchmarkSimpleQueries(b, 100000) }

func benchmarkAggregation(b *testing.B, numOfDocs int64) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		if _, _, err := getFirstLastByAggregation(ctx, numOfDocs); err != nil {
			b.Error(err)
		}
	}
}
func benchmarkSimpleQueries(b *testing.B, numOfDocs int64) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		if _, _, err := getFirstLastBySimpleQueries(ctx, numOfDocs); err != nil {
			b.Error(err)
		}
	}
}
