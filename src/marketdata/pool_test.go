package marketdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMostProfitablePool(t *testing.T) {
	pools := []*Pool{
		{
			Id: "0x123",
			DayRangeStats: PoolDayRangeStats{
				ProfitOverRange: 0.004,
			},
		},
		{
			Id: "0x456",
			DayRangeStats: PoolDayRangeStats{
				ProfitOverRange: 0.001,
			},
		},
		{
			Id: "0x567",
			DayRangeStats: PoolDayRangeStats{
				ProfitOverRange: 0.000,
			},
		},
		{
			Id: "0x89a",
			DayRangeStats: PoolDayRangeStats{
				ProfitOverRange: 0.145,
			},
		},
	}

	assert.Equal(t, pools[3].Id, MostProfitablePool(pools).Id)
}
