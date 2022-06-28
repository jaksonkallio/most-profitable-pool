package main

import (
	"log"

	"github.com/jaksonkallio/coding-challenge-messari/src/marketdata"
)

func main() {
	// TODO: remove all this stuff in the main function, currently exists for testing connection to UniSwap subgraph.

	/*
		dateRangeStart, _ := marketdata.ParseDate("2022-01-01")
		dateRangeEnd, _ := marketdata.ParseDate("2022-02-28")
	*/

	dateRangeStart, _ := marketdata.ParseDate("2022-05-27")
	dateRangeEnd, _ := marketdata.ParseDate("2022-06-27")

	pools, err := marketdata.FetchAllPools(dateRangeStart, dateRangeEnd)
	if err != nil {
		log.Fatalf("could not fetch pools: %s", err)
	}

	for _, pool := range pools {
		log.Println(pool.Pretty())

		/*if i >= 20 {
			log.Printf("%d additional pools omitted from display", len(pools)-i)
			break
		}*/
	}
}
