package smap

import (
	"hash/fnv"
	"sync"
)

// ConcurrentStringMap version of the segmentation implementation
type ConcurrentStringMap struct {
	segCount uint32
	segments []*SimpleConcurrentStringMap
}

func NewStringMap(seg int) *ConcurrentStringMap {
	m := &ConcurrentStringMap{segCount: uint32(seg)}
	for i := 0; i < seg; i++ {
		m.segments = append(m.segments, &SimpleConcurrentStringMap{
			mu:        sync.RWMutex{},
			stringMap: make(map[string]interface{}),
		})
	}
	return m
}

func (m *ConcurrentStringMap) Store(key string, value interface{}) {
	m.segment(key).Store(key, value)
}

func (m *ConcurrentStringMap) Load(key string) (value interface{}, ok bool) {
	return m.segment(key).Load(key)
}

func (m *ConcurrentStringMap) Range(f func(key string, value interface{}) bool) {
	for _, seg := range m.segments {
		seg.Range(f)
	}
}

func (m *ConcurrentStringMap) segment(key string) *SimpleConcurrentStringMap {
	hash := HashString(key)
	idx := hash % m.segCount
	return m.segments[idx]
}

// SimpleConcurrentStringMap simple implementation version based on rwmutex
type SimpleConcurrentStringMap struct {
	mu        sync.RWMutex
	stringMap map[string]interface{}
}

func (m *SimpleConcurrentStringMap) Store(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stringMap[key] = value
}

func (m *SimpleConcurrentStringMap) Load(key string) (value interface{}, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok = m.stringMap[key]
	return value, ok
}

func (m *SimpleConcurrentStringMap) Range(f func(key string, value interface{}) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.stringMap {
		if ok := f(k, v); !ok {
			break
		}
	}
}

func (m *SimpleConcurrentStringMap) Len() int {
	return len(m.stringMap)
}

func HashString(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}
