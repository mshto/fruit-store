package cache

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

func TestRedisPossitive(t *testing.T) {
	key := "test"
	value := "value"

	s, err := miniredis.Run()
	if err != nil {
		t.Fatal("failed to init miniredis")
	}
	defer s.Close()

	cache, err := New(Redis{
		Address: s.Addr(),
	})

	assert.NotNil(t, cache)
	assert.Nil(t, err)

	err = cache.Set(key, value, 2*time.Second)
	assert.Nil(t, err)

	result, err := cache.Get(key)
	assert.Equal(t, result, value)
	assert.Nil(t, err)

	err = cache.Del(key)
	assert.Nil(t, err)

	_, err = cache.Get(key)
	assert.NotNil(t, err)
}

func TestRedisDelNegative(t *testing.T) {
	key := "test"

	s, err := miniredis.Run()
	if err != nil {
		t.Fatal("failed to init miniredis")
	}
	defer s.Close()

	cache, err := New(Redis{
		Address: s.Addr(),
	})

	err = cache.Del(key)
	assert.NotNil(t, err)
}
