package marketdata

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/shurcooL/graphql"
)

const FetchBatchCount int = 10
const CouldNotParseErrMsg string = "could not parse %s value: %s"

type Pool struct {
	Id                  string
	TotalValueLockedUsd float64
	Token0Name          string
	Token1Name          string
	Days                []PoolDay
}

type PoolDay struct {
	Date    time.Time
	FeesUsd float64
	TvlUsd  float64
}

type DayRangeStats struct {
	Volume           float64
	Fees             float64
	ProfitOverRange  float64
	ProfitAnnualized float64
	Length           int
}

// Iteratively fetches all pools with fee data in date range.
func FetchAllPools(dateRangeStart time.Time, dateRangeEnd time.Time) ([]Pool, error) {
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
			// TODO: revert this by removing orderBy volume and orderDirection
		} `graphql:"pools(first: $count, where: { id_gt: $lastId }, orderBy: volumeUSD, orderDirection:desc )"`
	}

	page := 0
	lastId := ""
	pools := make([]Pool, 0)

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

			poolDays := make([]PoolDay, 0)

			for _, resPoolDay := range resPool.PoolDays {
				feesUsd, err := strconv.ParseFloat(string(resPoolDay.FeesUsd), 64)
				if err != nil {
					return nil, fmt.Errorf(CouldNotParseErrMsg, "fees USD", err)
				}

				tvlUsd, err := strconv.ParseFloat(string(resPoolDay.TvlUsd), 64)
				if err != nil {
					return nil, fmt.Errorf(CouldNotParseErrMsg, "tvl USD", err)
				}

				poolDays = append(
					poolDays,
					PoolDay{
						FeesUsd: feesUsd,
						TvlUsd:  tvlUsd,
						Date:    time.Unix(int64(resPoolDay.Date), 0),
					},
				)
			}

			pools = append(
				pools,
				Pool{
					Id:                  string(resPool.Id),
					TotalValueLockedUsd: totalValueLockedUsd,
					Token0Name:          string(resPool.Token0.Name),
					Token1Name:          string(resPool.Token1.Name),
					Days:                poolDays,
				},
			)
		}

		if len(res.Pools) > 0 {
			log.Printf("Fetched %d pools", len(res.Pools))
			lastId = string(res.Pools[len(res.Pools)-1].Id)
		} else {
			lastId = ""
		}

		page += 1

		// TODO: remove
		break
	}

	return pools, nil
}

func (pool *Pool) DayRangeStats() DayRangeStats {
	dayRangeStats := DayRangeStats{
		Length: len(pool.Days),
	}

	for _, day := range pool.Days {
		dayRangeStats.Volume += day.TvlUsd
		dayRangeStats.Fees += day.FeesUsd
	}

	// TODO: fix this calculation... I'm pretty sure its wrong
	if dayRangeStats.Volume != 0 {
		dayRangeStats.ProfitOverRange = dayRangeStats.Fees / dayRangeStats.Volume
	}

	if dayRangeStats.Length > 0 {
		dayRangeStats.ProfitAnnualized = (dayRangeStats.ProfitOverRange / float64(dayRangeStats.Length)) * 365
	}

	return dayRangeStats
}

func (pool *Pool) Pretty() string {
	dayRangeStats := pool.DayRangeStats()

	return fmt.Sprintf(
		"Pool %s / %s\n\tID: %s\n\tRange Length: %d\n\tProfit Rate Over Range: %.2f%%\n\tProfit Annualized (APR): %.2f%%",
		pool.Token0Name,
		pool.Token1Name,
		pool.Id,
		dayRangeStats.Length,
		math.Round(dayRangeStats.ProfitOverRange*10000)/float64(100),
		math.Round(dayRangeStats.ProfitAnnualized*10000)/float64(100),
	)
}
