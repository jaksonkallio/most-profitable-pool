package marketdata

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDate(t *testing.T) {
	date, err := ParseDate("2022-05-06")
	if err != nil {
		t.Fatalf("Bad range start date: %s", err)
	}

	assert.Equal(t, 2022, date.Year())
	assert.Equal(t, time.Month(5), date.Month())
	assert.Equal(t, 6, date.Day())
	assert.Equal(t, int64(1651795200), date.Unix())
}
