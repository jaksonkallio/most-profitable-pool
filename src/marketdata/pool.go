package marketdata

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/shurcooL/graphql"
)

// How many pools to get with each API call.
const FetchBatchCount int = 1000
const CouldNotParseErrMsg string = "could not parse %s value: %s"

// Representation of a particular pool with stats calculated from a date range.
type Pool struct {
	Id            string
	Token0Name    string
	Token1Name    string
	DayRangeStats PoolDayRangeStats
}

// Pool stats over a specified date range.
type PoolDayRangeStats struct {
	StartDate        time.Time
	EndDate          time.Time
	ProfitOverRange  float64
	ProfitAnnualized float64
}

// Iteratively fetches all pools with fee data in date range.
func FetchAllPools(dateRangeStart time.Time, dateRangeEnd time.Time, minTvl float64) ([]*Pool, error) {
	log.Printf("Fetching all liquidity pools and date range data from %v to %v where date TVL is at least %f", dateRangeStart, dateRangeEnd, minTvl)

	var res struct {
		Pools []struct {
			Id     graphql.String `graphql:"id"`
			Token0 struct {
				Name graphql.String `graphql:"name"`
			} `graphql:"token0"`
			Token1 struct {
				Name graphql.String `graphql:"name"`
			} `graphql:"token1"`
			PoolDays []struct {
				FeesUsd graphql.String `graphql:"feesUSD"`
				TvlUsd  graphql.String `graphql:"tvlUSD"`
				Date    graphql.Int    `graphql:"date"`
				// Assumes that we want to fetch dates in the start date and end date inclusively.
			} `graphql:"poolDayData(where: { date_gte: $dateRangeStart, date_lte: $dateRangeEnd })"`
			// Paginate using the best-practice outlined in The Graph documentation instead of using `skip`.
		} `graphql:"pools(first: $count, where: { id_gt: $lastId })"`
	}

	page := 0
	lastId := ""
	pools := make([]*Pool, 0)

	for {
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

		for _, resPool := range res.Pools {
			// Init a range stats object.
			poolDayRangeStats := PoolDayRangeStats{
				StartDate: dateRangeStart,
				EndDate:   dateRangeEnd,
			}

			exceedsMinTvl := false
			var profit float64

			// Iterate over each pool day.
			for _, resPoolDay := range resPool.PoolDays {
				// Parse the day's fees represented in USD.
				feesUsd, err := strconv.ParseFloat(string(resPoolDay.FeesUsd), 64)
				if err != nil {
					return nil, fmt.Errorf(CouldNotParseErrMsg, "fees USD", err)
				}

				// Parse the day's TVL represented in USD.
				tvlUsd, err := strconv.ParseFloat(string(resPoolDay.TvlUsd), 64)
				if err != nil {
					return nil, fmt.Errorf(CouldNotParseErrMsg, "tvl USD", err)
				}

				// Check if the TVL on this date exceeds minimum.
				if !exceedsMinTvl && tvlUsd > minTvl {
					exceedsMinTvl = true
				}

				// Skip where TVL is non-positive, i.e. we'd get a divide-by-zero issue when we calculate profit.
				if tvlUsd <= 0 {
					continue
				}

				// The profit for this day is calculated as the fees collected divided by the total locked value.
				// This gives us the profit per $1 USD locked ("deposited") on this day.
				// We're summing all of the daily profit-per-$1 numbers into a `profit` float.
				profit += feesUsd / tvlUsd
			}

			// Verify that on some date in the range, the min TVL was exceeded.
			// Ignore pools where the min TVL was never exceeded in the provided date range.
			if exceedsMinTvl {
				poolDayRangeStats.ProfitOverRange = profit

				if poolDayRangeStats.Length() > 0 {
					// Annualized profit assumes there are exactly 365 days in a year.
					// Divide profit received over the range by the number of days in the range to get the average daily return, then multiply by 365 to get the annual return.
					poolDayRangeStats.ProfitAnnualized = (poolDayRangeStats.ProfitOverRange / poolDayRangeStats.Length()) * 365
				}

				// Add new pool to the resulting pools slice.
				pools = append(
					pools,
					&Pool{
						Id:            string(resPool.Id),
						Token0Name:    string(resPool.Token0.Name),
						Token1Name:    string(resPool.Token1.Name),
						DayRangeStats: poolDayRangeStats,
					},
				)
			}
		}

		if len(res.Pools) > 0 {
			lastId = string(res.Pools[len(res.Pools)-1].Id)
			log.Printf("Fetched %d pools", len(res.Pools))
		} else {
			// Break once we no longer receive any pools from our query.
			break
		}

		// Increment page counter.
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
		"\n\tPool Address: %s\n\tTokens: %s <-> %s\n\tRange Length: %f\n\tProfit Over Range (Earned per $1 USD Deposited): %f\n\tProfit Annualized (APR): %.2f%%",
		pool.Id,
		pool.Token0Name,
		pool.Token1Name,
		pool.DayRangeStats.Length(),
		pool.DayRangeStats.ProfitOverRange,
		math.Round(pool.DayRangeStats.ProfitAnnualized*10000)/float64(100),
	)
}

// Returns the length of the pool stats range in days.
func (poolDayRangeStats *PoolDayRangeStats) Length() float64 {
	return poolDayRangeStats.EndDate.Sub(poolDayRangeStats.StartDate).Hours() / 24
}
