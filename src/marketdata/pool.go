package marketdata

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/shurcooL/graphql"
)

const FetchBatchCount int = 1000
const CouldNotParseErrMsg string = "could not parse %s value: %s"

type Pool struct {
	Id                  string
	TotalValueLockedUsd float64
	Token0Name          string
	Token1Name          string
	DayRangeStats       PoolDayRangeStats
}

type PoolDayRangeStats struct {
	SumTotalValueLocked float64
	SumFees             float64
	ProfitOverRange     float64
	ProfitAnnualized    float64
	Length              int
}

// Iteratively fetches all pools with fee data in date range.
func FetchAllPools(dateRangeStart time.Time, dateRangeEnd time.Time) ([]*Pool, error) {
	log.Printf("Fetching all liquidity pools and date range data from %v to %v", dateRangeStart, dateRangeEnd)

	var res struct {
		Pools []struct {
			Id                  graphql.String `graphql:"id"`
			TotalValueLockedUsd graphql.String `graphql:"totalValueLockedUSD"`
			Token0              struct {
				Name graphql.String `graphql:"name"`
			} `graphql:"token0"`
			Token1 struct {
				Name graphql.String `graphql:"name"`
			} `graphql:"token1"`
			PoolDays []struct {
				FeesUsd graphql.String `graphql:"feesUSD"`
				TvlUsd  graphql.String `graphql:"tvlUSD"`
				Date    graphql.Int    `graphql:"date"`
			} `graphql:"poolDayData(where: { date_gte: $dateRangeStart, date_lte: $dateRangeEnd })"`
		} `graphql:"pools(first: $count, where: { id_gt: $lastId } )"`
	}

	page := 0
	lastId := ""
	pools := make([]*Pool, 0)

	for page == 0 || len(lastId) > 0 {
		err := GraphQuery(
			&res,
			map[string]interface{}{
				"count":          graphql.Int(FetchBatchCount),
				"lastId":         graphql.String(lastId),
				"dateRangeStart": graphql.Int(dateRangeStart.Unix()),
				"dateRangeEnd":   graphql.Int(dateRangeEnd.Unix()),
			},
		)
		if err != nil {
			return nil, fmt.Errorf("could not get pools: %s", err)
		}

		// TODO: more efficient way of appending to resulting pools slice.
		for _, resPool := range res.Pools {
			totalValueLockedUsd, err := strconv.ParseFloat(string(resPool.TotalValueLockedUsd), 64)
			if err != nil {
				return nil, fmt.Errorf(CouldNotParseErrMsg, "total value locked USD", err)
			}

			poolDayRangeStats := PoolDayRangeStats{
				Length: len(resPool.PoolDays),
			}

			for _, resPoolDay := range resPool.PoolDays {
				feesUsd, err := strconv.ParseFloat(string(resPoolDay.FeesUsd), 64)
				if err != nil {
					return nil, fmt.Errorf(CouldNotParseErrMsg, "fees USD", err)
				}

				tvlUsd, err := strconv.ParseFloat(string(resPoolDay.TvlUsd), 64)
				if err != nil {
					return nil, fmt.Errorf(CouldNotParseErrMsg, "tvl USD", err)
				}

				poolDayRangeStats.SumFees += feesUsd
				poolDayRangeStats.SumTotalValueLocked += tvlUsd
			}

			if poolDayRangeStats.SumTotalValueLocked != 0 {
				poolDayRangeStats.ProfitOverRange = poolDayRangeStats.SumFees / poolDayRangeStats.SumTotalValueLocked
			}

			if poolDayRangeStats.Length > 0 {
				poolDayRangeStats.ProfitAnnualized = (poolDayRangeStats.ProfitOverRange / float64(poolDayRangeStats.Length)) * 365
			}

			pools = append(
				pools,
				&Pool{
					Id:                  string(resPool.Id),
					TotalValueLockedUsd: totalValueLockedUsd,
					Token0Name:          string(resPool.Token0.Name),
					Token1Name:          string(resPool.Token1.Name),
					DayRangeStats:       poolDayRangeStats,
				},
			)
		}

		if len(res.Pools) > 0 {
			lastId = string(res.Pools[len(res.Pools)-1].Id)
			log.Printf("Fetched %d pools", len(res.Pools))
		} else {
			lastId = ""
		}

		page += 1
	}

	return pools, nil
}

// Finds the most profitable pool, given a slice of pools.
// If two pools are equally profitable, chooses the one that comes earlier.
func MostProfitablePool(pools []*Pool) *Pool {
	var mostProfitablePool *Pool

	for _, pool := range pools {
		if mostProfitablePool == nil || pool.DayRangeStats.ProfitOverRange > mostProfitablePool.DayRangeStats.ProfitOverRange {
			mostProfitablePool = pool
		}
	}

	return mostProfitablePool
}

// Pretty-prints the pool by returning a string.
func (pool *Pool) Pretty() string {
	return fmt.Sprintf(
		"\n\tPool Address: %s\n\tTokens: %s <-> %s\n\tRange Length: %d\n\tProfit Rate Over Range: %.2f%%\n\tProfit Annualized (APR): %.2f%%",
		pool.Id,
		pool.Token0Name,
		pool.Token1Name,
		pool.DayRangeStats.Length,
		math.Round(pool.DayRangeStats.ProfitOverRange*10000)/float64(100),
		math.Round(pool.DayRangeStats.ProfitAnnualized*10000)/float64(100),
	)
}
