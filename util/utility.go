package util

import (
	"sync"
)

func GetSyncMapSize(m *sync.Map) int {
	count := 0
	m.Range(func(key, value interface{}) bool {
		count++
		return true // Continue iterating
	})
	return count
}
