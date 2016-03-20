package lrucache

import (
	"github.com/gregjones/httpcache"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
)

var entries = []struct {
	key   string
	value string
}{
	{"1", "one"},
	{"2", "two"},
	{"3", "three"},
	{"4", "four"},
	{"5", "five"},
}

func TestInterface(t *testing.T) {
	var h httpcache.Cache
	h = New(1000000)
	if assert.NotNil(t, h) {
		_, ok := h.Get("missing")
		assert.Equal(t, ok, false)
	}
}

func TestCache(t *testing.T) {
	c := New(1000000)

	for _, e := range entries {
		c.Set(e.key, []byte(e.value))
	}

	c.Delete("missing")
	_, ok := c.Get("missing")
	assert.Equal(t, ok, false)

	for _, e := range entries {
		value, ok := c.Get(e.key)
		if assert.Equal(t, ok, true) {
			assert.Equal(t, string(value), e.value)
		}
	}

	for _, e := range entries {
		c.Delete(e.key)

		_, ok := c.Get(e.key)
		assert.Equal(t, ok, false)
	}
}

func TestSize(t *testing.T) {
	c := New(1000000)
	assert.Equal(t, c.size, 0)

	// Check that size is overhead + len(key) + len(value)
	c.Set("some", []byte("text"))
	assert.Equal(t, c.size, entryOverhead+8)

	// Replace key
	c.Set("some", []byte("longer text"))
	assert.Equal(t, c.size, entryOverhead+15)

	c.Delete("some")
	assert.Equal(t, c.size, 0)
}

func TestEvict(t *testing.T) {
	c := New(entryOverhead*2 + 20)

	for _, e := range entries {
		c.Set(e.key, []byte(e.value))
	}

	// Make sure only the last two entries were kept.
	assert.Equal(t, c.size, entryOverhead*2+10)
}

func BenchmarkSet(b *testing.B) {
	v := []byte("value")

	c := benchSetup(b, 10000000, 10000)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Set(randKey(10000), v)
		}
	})
}

func BenchmarkGet(b *testing.B) {
	c := benchSetup(b, 10000000, 10000)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Get(randKey(20000))
		}
	})
}

func BenchmarkSetGetDelete(b *testing.B) {
	v := []byte("value")

	c := benchSetup(b, 10000000, 10000)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Set(randKey(10000), v)
			_, _ = c.Get(randKey(20000))
			c.Delete(randKey(10000))
		}
	})
}

func benchSetup(b *testing.B, size int64, entries int) *LruCache {
	c := New(size)

	v := []byte("value")
	for i := 0; i < entries; i++ {
		c.Set(strconv.Itoa(i), v)
	}

	b.ResetTimer()

	return c
}

func randKey(n int32) string {
	return strconv.Itoa(int(rand.Int31n(n)))
}
