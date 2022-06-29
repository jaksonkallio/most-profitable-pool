# Most Profitable Pool
Simple project to calculate the most profitable UniSwap v3 pool over a given time range.

## Testing
Run `<repo root>/test.sh`.

## Building
Run `<repo root>/build.sh`. This will create the binary `most-profitable-pool` in `<repo root>/dist` subdirectory.

## Running
There are three command line arguments:
1. Start date in `YYYY-MM-DD` (ISO) format (inclusive).
2. End date in `YYYY-MM-DD` (ISO) format (inclusive).
3. Minimum TVL. Will exclude pools that never exceeded the specified minimum TVL value in the specified date range. Helpful for filtering out outliers. Set to `0` to allow all pools.

### Example
For example, `./dist/most-profitable-pool 2022-01-01 2022-02-28 100000` will return:
```
2022/06/29 08:14:39 Fetching all liquidity pools and date range data from 2022-01-01 00:00:00 +0000 UTC to 2022-02-28 00:00:00 +0000 UTC where date TVL is at least 100000.000000
2022/06/29 08:14:46 Fetched 1000 pools
2022/06/29 08:14:48 Fetched 1000 pools
2022/06/29 08:14:52 Fetched 1000 pools
2022/06/29 08:14:55 Fetched 1000 pools
2022/06/29 08:14:57 Fetched 1000 pools
2022/06/29 08:15:00 Fetched 1000 pools
2022/06/29 08:15:03 Fetched 1000 pools
2022/06/29 08:15:03 Fetched 137 pools
2022/06/29 08:15:03 Most profitable pool:
	Pool Address: 0x9396c357befc79abfef7f229a3bd8dd0ae8e6bfd
	Tokens: Shping Coin <-> Wrapped Ether
	Range Length: 58.000000
	Profit Over Range (Earned per $1 USD Deposited): 3.479155
	Profit Annualized (APR): 2189.47%
```