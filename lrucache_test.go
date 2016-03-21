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
			assert.Equal(t, string(value), string(e.value))
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
	assert.Equal(t, c.size, int64(0))

	// Check that size is overhead + len(key) + len(value)
	c.Set("some", []byte("text"))
	assert.Equal(t, c.size, int64(entryOverhead+8))

	// Replace key
	c.Set("some", []byte("longer text"))
	assert.Equal(t, c.size, int64(entryOverhead+15))

	assert.Equal(t, c.Size(), c.size)

	c.Delete("some")
	assert.Equal(t, c.size, int64(0))
}

func TestEvict(t *testing.T) {
	c := New(entryOverhead*2 + 20)

	for _, e := range entries {
		c.Set(e.key, []byte(e.value))
	}

	// Make sure only the last two entries were kept.
	assert.Equal(t, c.size, int64(entryOverhead*2+10))
}

func TestRace(t *testing.T) {
	c := New(100000)

	for worker := 0; worker < 8; worker++ {
		go testRaceWorker(c)
	}
}

func testRaceWorker(c *LruCache) {
	v := []byte("value")

	for n := 0; n < 1000; n++ {
		c.Set(randKey(100), v)
		_, _ = c.Get(randKey(200))
		c.Delete(randKey(100))
		_ = c.Size()
	}
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
			_, _ = c.Get(randKey(20000))
		}
	})
}

func BenchmarkSize(b *testing.B) {
	c := benchSetup(b, 10000000, 10000)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = c.Size()
		}
	})
}

func BenchmarkSetGetDeleteSize(b *testing.B) {
	v := []byte("value")

	c := benchSetup(b, 10000000, 10000)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Set(randKey(10000), v)
			_, _ = c.Get(randKey(20000))
			c.Delete(randKey(10000))
			_ = c.Size()
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
