package util

import (
	"time"
)

func GetSubmitUnixSeconds() int {
	// 2055/01/01 00:00:00 - Now
	MaxUnixTime := time.Date(2055, time.January, 1, 0, 0, 0, 0, time.UTC)
	now := time.Now()
	secondsUntilMaxTime := MaxUnixTime.Unix() - now.Unix()

	return int(secondsUntilMaxTime)
}
