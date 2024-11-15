package util

import (
	"time"
)

// a utility library for functions inolving time

// returns the current time in UTC.
func CurrentTimeUTC() time.Time {
	return time.Now().UTC()
}
