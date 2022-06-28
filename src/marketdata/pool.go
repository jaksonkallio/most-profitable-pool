package marketdata

import (
	"fmt"
	"log"
	"strconv"

	"github.com/shurcooL/graphql"
)

const FetchBatchCount int = 1000

type Pool struct {
	Id                  string
	TotalValueLockedUsd float64
	Token0Name          string
	Token1Name          string
}

// Iteratively fetches all pools.
func FetchAllPools() ([]Pool, error) {
	var res struct {
		Pools []struct {
			Id                  graphql.String `graphql:"id"`
			TotalValueLockedUsd graphql.String `graphql:"totalValueLockedUSD"`
			Token0              struct {
				Name graphql.String `graphql:"name"`
			}
			Token1 struct {
				Name graphql.String `graphql:"name"`
			}
		} `graphql:"pools(first: $count, where: { id_gt: $lastId })"`
	}

	page := 0
	lastId := ""
	pools := make([]Pool, 0)

	for page == 0 || len(lastId) > 0 {
		err := GraphQuery(
			&res,
			map[string]interface{}{
				"count":  graphql.Int(FetchBatchCount),
				"lastId": graphql.String(lastId),
			},
		)
		if err != nil {
			return nil, fmt.Errorf("could not get pools: %s", err)
		}

		// TODO: more efficient way of appending to resulting pools slice.
		for _, resPool := range res.Pools {
			totalValueLockedUsd, err := strconv.ParseFloat(string(resPool.TotalValueLockedUsd), 64)
			if err != nil {
				return nil, fmt.Errorf("could not parse total value locked USD value: %s", err)
			}

			pools = append(
				pools,
				Pool{
					Id:                  string(resPool.Id),
					TotalValueLockedUsd: totalValueLockedUsd,
					Token0Name:          string(resPool.Token0.Name),
					Token1Name:          string(resPool.Token0.Name),
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
	}

	return pools, nil
}
