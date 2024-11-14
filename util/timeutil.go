package util

import (
	"time"
)

// Returns the current time in UTC.
func CurrentTimeUTC() time.Time {
	return time.Now().UTC()
}
