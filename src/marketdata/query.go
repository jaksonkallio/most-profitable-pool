package marketdata

import (
	"context"
	"fmt"

	"github.com/shurcooL/graphql"
)

const endpoint string = "https://api.thegraph.com/subgraphs/name/ianlapham/uniswap-v3-subgraph"

var graphQlClient *graphql.Client

func init() {
	graphQlClient = graphql.NewClient(endpoint, nil)
}

// Queries the Uniswap V3 subgraph.
func GraphQuery(query interface{}, variables map[string]interface{}) error {
	err := graphQlClient.Query(context.Background(), query, variables)
	if err != nil {
		return fmt.Errorf("query failed: %s", err)
	}

	return nil
}
