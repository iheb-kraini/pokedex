package main

import (
	"strings"
	"testing"
)

const areaEncountersJSON = `{"pokemon_encounters":[` +
	`{"pokemon":{"name":"pikachu","url":"https://pokeapi.co/api/v2/pokemon/25/"}},` +
	`{"pokemon":{"name":"geodude","url":"https://pokeapi.co/api/v2/pokemon/74/"}}]}`

func TestCommandExplorePrintsEncounters(t *testing.T) {
	cfg := newTestConfig()
	cfg.args = []string{"stage"}
	if _, ok := cfg.cache.Add("stage", []byte(areaEncountersJSON)); ok {
		t.Fatalf("expected fresh cache key")
	}

	var err error
	out := captureStdout(t, func() { err = commandExplore(cfg) })
	if err != nil {
		t.Fatalf("commandExplore: %v", err)
	}
	if !strings.Contains(out, "pikachu") || !strings.Contains(out, "geodude") {
		t.Errorf("expected encounter names in output, got %q", out)
	}
}

func TestCommandExploreInvalidCachedJSON(t *testing.T) {
	cfg := newTestConfig()
	cfg.args = []string{"stage"}
	cfg.cache.Add("stage", []byte("not-json"))

	var err error
	captureStdout(t, func() { err = commandExplore(cfg) })
	if err == nil {
		t.Fatal("expected unmarshal error for non-JSON cached body, got nil")
	}
}
