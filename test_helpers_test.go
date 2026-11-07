package main

import (
	"io"
	"os"
	"testing"
	"time"

	cache "github.com/iheb-kraini/pokedex/pokecache"
)

// captureStdout redirects os.Stdout to a pipe while fn runs and returns everything fn printed.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	screen := os.Stdout
	os.Stdout = w
	defer func() {
		os.Stdout = screen
		r.Close()
	}()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("close pipe writer: %v", err)
	}
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read captured stdout: %v", err)
	}
	return string(data)
}

// newTestConfig returns a config with isolated fields and its own cache,
// safe to mutate in a test without leaking state between tests.
func newTestConfig() *config {
	return &config{
		commands: map[string]cliCommand{},
		next:     "",
		prev:     "",
		cache:    cache.NewCache(5 * time.Minute),
		args:     []string{},
		pokedex:  map[string]Pokemon{},
	}
}
