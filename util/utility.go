package util

import (
	"fmt"
	"sync"
)

// a general utility library

// returns the size of a sync.Map object (i.e. the number of entries that it currently has)
func GetSyncMapSize(m *sync.Map) int {
	count := 0
	m.Range(func(key, value interface{}) bool {
		count++
		return true // Continue iterating
	})
	return count
}

// convert a number of bytes to a legible string
func FormatBytes(bytes uint64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
		tb = gb * 1024
	)
	switch {
	case bytes >= tb:
		return fmt.Sprintf("%.3g TB", float64(bytes)/float64(tb))
	case bytes >= gb:
		return fmt.Sprintf("%.3g GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.3g MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.3g KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}
