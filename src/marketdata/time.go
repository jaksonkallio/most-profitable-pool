package marketdata

import "time"

const StandardDateFormat string = "2006-01-02"

func ParseDate(dateStr string) (time.Time, error) {
	t, err := time.Parse(StandardDateFormat, dateStr)
	if err != nil {
		return time.Now(), err
	}

	return t, nil
}
