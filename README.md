# Most Profitable Pool
Simple project to calculate the most profitable UniSwap v3 pool over a given time range.

## Building
Run `<repo root>/build.sh`. This will create the binary `most-profitable-pool` in `<repo root>/dist` subdirectory.

## Running
There are three command line arguments:
1. Start date in `YYYY-MM-DD` (ISO) format.
2. End date in `YYYY-MM-DD` (ISO) format.
3. Minimum TVL. Will exclude pools that never exceeded the minimum TVL value in the specified date range. Helpful for filtering out outliers caused by sparse data. Set to `0` to allow all pools.

### Example
For example, `./dist/most-profitable-pool 2022-01-01 2022-02-28 100000` will return:
```
2022/06/28 13:45:50 Fetching all liquidity pools and date range data from 2022-01-01 00:00:00 +0000 UTC to 2022-02-28 00:00:00 +0000 UTC where date TVL is at least 100000.000000
2022/06/28 13:45:55 Fetched 1000 pools
2022/06/28 13:45:58 Fetched 1000 pools
2022/06/28 13:46:02 Fetched 1000 pools
2022/06/28 13:46:06 Fetched 1000 pools
2022/06/28 13:46:11 Fetched 1000 pools
2022/06/28 13:46:13 Fetched 1000 pools
2022/06/28 13:46:17 Fetched 1000 pools
2022/06/28 13:46:18 Fetched 130 pools
2022/06/28 13:46:18 Most profitable pool:
	Pool Address: 0x63805e5d951398bc1c1bec242d303f59fa7732e3
	Tokens: X2Y2Token <-> Wrapped Ether
	Range Length: 13
	Profit Rate Over Range: 7.42%
	Profit Annualized (APR): 208.44%
```