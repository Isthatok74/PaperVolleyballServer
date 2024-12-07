package util

import (
	"sync"
	"testing"
)

func TestSyncMapSize(t *testing.T) {
	m := &sync.Map{}
	m.Store("key1", "value1")
	m.Store("key2", "value2")
	count := GetSyncMapSize(m)
	t.Logf("sync.Map has %d items", count)

	if count != 2 {
		t.Errorf("sync.Map size = %d; want %d", count, 2)
	}
}

func TestSyncMapCopy(t *testing.T) {
	m := &sync.Map{}
	m.Store("key1", "value1")
	m.Store("key2", "value2")

	copy := CopySyncMap(m)
	val, ok := copy.Load("key2")
	if !ok {
		t.Errorf("error copying sync.Map")
	}
	if val != "value2" {
		t.Errorf("sync.Map copy, value stored for key2 = %s; want %s", val, "value2")
	}
}
