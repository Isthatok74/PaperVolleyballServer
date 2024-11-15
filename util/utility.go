package util

import (
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
