package cache

import(
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestCache_SetGetDel(t *testing.T) {
	cache := NewCache[string](5 * time.Second)

	// Setting and getting a value from the cache
	cache.Set("key1", "value1")
	value, found := cache.Get("key1")
	assert.True(t, found)
	assert.Equal(t, "value1", value)

	// Getting a non-existent value from the cache
	_, found = cache.Get("key2")
	assert.False(t, found)

	// Setting and getting multiple values from the cache
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	value, found = cache.Get("key2")
	assert.True(t, found)
	assert.Equal(t, "value2", value)
	value, found = cache.Get("key3")
	assert.True(t, found)
	assert.Equal(t, "value3", value)

	// Deleting a value from the cache
	cache.Delete("key2")
	_, found = cache.Get("key2")
	assert.False(t, found)
}

func TestCache_Expiration(t *testing.T) {
	cache := NewCache[string](5 * time.Second)

	// Set, get, and check a value in the cache with no delay
	cache.Set("key1", "value1")
	value, found := cache.Get("key1")
	assert.True(t, found)
	assert.Equal(t, "value1", value)

	// Set, get, and check a value in the cache with a small delay 
	cache.Set("key2", "value2")
	time.Sleep(2 * time.Second)
	value, found = cache.Get("key2")
	assert.True(t, found)
	assert.Equal(t, "value2", value)

	// Wait for the value to expire
	time.Sleep(4 * time.Second)

	// Check if the value has expired
	_, found = cache.Get("key1")
	assert.False(t, found)
}
