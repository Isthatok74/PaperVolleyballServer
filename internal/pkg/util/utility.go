package util

import (
	"fmt"
	"net/http"
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

// creates a copy of a sync map without copying the mutex lock
func CopySyncMap(m *sync.Map) *sync.Map {
	copyMap := &sync.Map{}

	m.Range(func(key, value interface{}) bool {
		copyMap.Store(key, value)
		return true
	})
	return copyMap
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

// returns the size of the http header
func HTTPHeaderSize(header http.Header) uint64 {
	headerSize := 0
	for name, values := range header {
		for _, value := range values {
			headerSize += len(name) + len(value)
		}
	}
	// Add the final "\r\n" after headers
	headerSize += 2
	return uint64(headerSize)
}
