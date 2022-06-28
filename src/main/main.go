package main

import (
	"log"

	"github.com/jaksonkallio/coding-challenge-messari/src/marketdata"
)

func main() {
	dateRangeStart, err := marketdata.ParseDate("2022-01-01")
	if err != nil {
		log.Fatalf("Bad range start date: %s", err)
	}

	dateRangeEnd, err := marketdata.ParseDate("2022-02-28")
	if err != nil {
		log.Fatalf("Bad range end date: %s", err)
	}

	if dateRangeStart.After(dateRangeEnd) {
		log.Fatalf("Date range start is after date range end")
	}

	pools, err := marketdata.FetchAllPools(dateRangeStart, dateRangeEnd)
	if err != nil {
		log.Fatalf("Could not fetch pools: %s", err)
	}

	mostProfitablePool := marketdata.MostProfitablePool(pools)
	if mostProfitablePool == nil {
		log.Printf("No pool is considered most profitable")
	} else {
		log.Printf("Most profitable pool: %s", mostProfitablePool.Pretty())
	}
}
