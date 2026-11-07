package cache

import (
	"bytes"
	"testing"
	"time"
)

// TestCacheAddAndGet tests basic adding and retrieving of cache entries.
func TestCacheAddAndGet(t *testing.T) {
	cache := NewCache(5 * time.Minute)
	key := "https://pokeapi.co/api/v2/pokemon/ditto"
	val := []byte("some-json-data-for-ditto")

	cache.Add(key, val)

	got, exists := cache.Get(key)
	if !exists {
		t.Fatalf("expected key %q to exist in cache, but it didn't", key)
	}

	if !bytes.Equal(got, val) {
		t.Errorf("expected cached bytes %q, got %q", val, got)
	}
}

// TestCacheGetMissing tests behavior when requesting a key that does not exist.
func TestCacheGetMissing(t *testing.T) {
	cache := NewCache(5 * time.Minute)
	_, exists := cache.Get("non_existent_key")
	if exists {
		t.Errorf("expected key to not exist in cache, but Get returned true")
	}
}

// TestCacheReplaceValue tests that adding an existing key replaces its old value.
func TestCacheReplaceValue(t *testing.T) {
	cache := NewCache(5 * time.Minute)
	key := "pokemon-pikachu"
	val1 := []byte("old-data")
	val2 := []byte("new-data")

	cache.Add(key, val1)
	cache.Add(key, val2) // Should replace val1

	got, exists := cache.Get(key)
	if !exists {
		t.Fatalf("expected key %q to exist", key)
	}

	if !bytes.Equal(got, val2) {
		t.Errorf("expected replaced value %q, got %q", val2, got)
	}
}

// TestCacheMultipleKeys ensures different keys in the same cache don't overwrite each other.
func TestCacheMultipleKeys(t *testing.T) {
	cache := NewCache(5 * time.Minute)

	entries := map[string][]byte{
		"bulbasaur":  []byte("data-bulbasaur"),
		"charmander": []byte("data-charmander"),
		"squirtle":   []byte("data-squirtle"),
	}

	// Add all entries
	for k, v := range entries {
		cache.Add(k, v)
	}

	// Verify all entries exist and match their respective values
	for k, expectedVal := range entries {
		got, exists := cache.Get(k)
		if !exists {
			t.Errorf("expected key %q to exist", k)
			continue
		}
		if !bytes.Equal(got, expectedVal) {
			t.Errorf("for key %q: expected %q, got %q", k, expectedVal, got)
		}
	}
}

// TestCacheExpiration checks that fresh entries stay available and old entries are eventually removed.
func TestCacheExpiration(t *testing.T) {
	interval := 50 * time.Millisecond
	cache := NewCache(interval)
	key := "expired-key"
	val := []byte("data")

	cache.Add(key, val)

	// 1. Check that the fresh entry is available immediately
	if _, exists := cache.Get(key); !exists {
		t.Fatalf("expected fresh entry to exist immediately in cache")
	}

	// 2. Wait slightly longer than the interval for the background reap loop to clear it
	time.Sleep(100 * time.Millisecond)

	// 3. Check that the entry has been removed
	if _, exists := cache.Get(key); exists {
		t.Errorf("expected entry to have expired and been removed from cache")
	}
}
