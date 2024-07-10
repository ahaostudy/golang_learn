package smap

import (
	"strings"
	"sync"
	"testing"

	"github.com/alecthomas/assert"
)

func TestNewStringMap(t *testing.T) {
	const cnt = 1000000
	m := NewStringMap(100)
	ch := make(chan struct {
		K string
		V int
	}, cnt)

	wg := sync.WaitGroup{}
	for i := 0; i < cnt; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			k := conv.Itoa(v)
			ch <- struct {
				K string
				V int
			}{K: k, V: v}
			m.Store(k, v)
		}(i)
	}
	wg.Wait()
	close(ch)

	for x := range ch {
		_, ok := m.Load(x.K)
		assert.True(t, ok)
		// since the key can be the same and the value will be overwritten by the new value, it is not necessarily the same here
		// assert.Equal(t, v, x.V)
	}

	maxLen, minLen := 0, cnt
	for _, seg := range m.segments {
		if seg.Len() > maxLen {
			maxLen = seg.Len()
		}
		if seg.Len() < minLen {
			minLen = seg.Len()
		}
	}
	t.Logf("Maximum segment length: %d, minimum segment length: %d", maxLen, minLen)
}

func BenchmarkConcurrentStringMap_Store(b *testing.B) {
	m := NewStringMap(100)
	wg := sync.WaitGroup{}
	for n := 0; n < b.N; n++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			k := conv.Itoa(i)
			m.Store(k, i)
		}(n)
	}
	wg.Wait()
}

func BenchmarkSyncMap_Store(b *testing.B) {
	m := sync.Map{}
	wg := sync.WaitGroup{}
	for n := 0; n < b.N; n++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			k := conv.Itoa(i)
			m.Store(k, i)
		}(n)
	}
	wg.Wait()
}

var (
	storedConcurrentStringMap = NewStringMap(100)
	storedSyncMap             = sync.Map{}
)

func init() {
	for i := 0; i < 123456; i++ {
		k := conv.Itoa(i)
		storedConcurrentStringMap.Store(k, i)
		storedSyncMap.Store(k, i)
	}
}

func BenchmarkConcurrentStringMap_Load(b *testing.B) {
	wg := sync.WaitGroup{}
	for n := 0; n < b.N; n++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			k := conv.Itoa(i)
			storedConcurrentStringMap.Load(k)
		}(n)
	}
	wg.Wait()
}

func BenchmarkSyncMap_Load(b *testing.B) {
	wg := sync.WaitGroup{}
	for n := 0; n < b.N; n++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			k := conv.Itoa(i)
			storedSyncMap.Load(k)
		}(n)
	}
	wg.Wait()
}

func BenchmarkConcurrentStringMap_StoreAndLoad(b *testing.B) {
	m := NewStringMap(100)
	wg := sync.WaitGroup{}
	for n := 0; n < b.N; n++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			k := conv.Itoa(i)
			if i%2 == 0 {
				m.Load(k)
			} else {
				m.Store(k, i)
			}
		}(n)
	}
	wg.Wait()
}

func BenchmarkSyncMap_StoreAndLoad(b *testing.B) {
	m := sync.Map{}
	wg := sync.WaitGroup{}
	for n := 0; n < b.N; n++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			k := conv.Itoa(i)
			if i%2 == 0 {
				m.Load(k)
			} else {
				m.Store(k, i)
			}
		}(n)
	}
	wg.Wait()
}

func TestHashString(t *testing.T) {
	for i := 0; i < 10000; i++ {
		s := conv.Itoa(i)
		t.Logf("hash: %s - %d", s, HashString(s))
	}
}

var conv = Converter{base: 123456}

type Converter struct {
	base int
}

func (conv *Converter) Itoa(n int) string {
	const (
		Int64Max = (1 << 63) - 1
		charset  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+/"
	)
	var x = int64(n % conv.base * (Int64Max / conv.base))
	var builder strings.Builder
	builder.Grow(11)
	for i := 0; i < 64; i += 6 {
		index := (x >> i) & 0x3F
		builder.WriteByte(charset[index])
	}
	return builder.String()
}
