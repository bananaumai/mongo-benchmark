# Queries to get first and last document in a collection

This code provides the comparison of the following two queries.

* Use $group aggregation query
* Query simple queries twice

The result was following:

```
2021/05/17 18:55:01 insert 1000 docs
2021/05/17 18:55:01 insert 10000 docs
2021/05/17 18:55:01 insert 100000 docs
2021/05/17 18:55:02 run test
goos: darwin
goarch: amd64
pkg: github.com/bananaumai/mongo-benchmark/query-first-last
cpu: Intel(R) Core(TM) i7-8700B CPU @ 3.20GHz
Benchmark_getFirstLastByAggregation_1000
Benchmark_getFirstLastByAggregation_1000-12        	     718	   1472230 ns/op
Benchmark_getFirstLastByAggregation_10000
Benchmark_getFirstLastByAggregation_10000-12       	      96	  12005398 ns/op
Benchmark_getFirstLastByAggregation_100000
Benchmark_getFirstLastByAggregation_100000-12      	       9	 116105565 ns/op
Benchmark_getFirstLastBySimpleQueries_1000
Benchmark_getFirstLastBySimpleQueries_1000-12      	    3661	    309350 ns/op
Benchmark_getFirstLastBySimpleQueries_10000
Benchmark_getFirstLastBySimpleQueries_10000-12     	    3790	    333709 ns/op
Benchmark_getFirstLastBySimpleQueries_100000
Benchmark_getFirstLastBySimpleQueries_100000-12    	    3319	    333679 ns/op
PASS
```

According to the result, it's way more efficient to use simple queries than to use aggregation.
