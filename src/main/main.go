package main

import (
	"log"

	"github.com/jaksonkallio/coding-challenge-messari/src/marketdata"
)

func main() {
	// TODO: remove all this stuff in the main function, currently exists for testing connection to UniSwap subgraph.

	pools, err := marketdata.FetchAllPools()
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	for i, pool := range pools {
		log.Printf("%#v", pool)

		if i >= 20 {
			log.Printf("%d additional pools omitted from display", len(pools)-i)
			break
		}
	}
}
