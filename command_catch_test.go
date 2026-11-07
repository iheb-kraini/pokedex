package main

import (
	"strings"
	"testing"
)

const pikachuJSON = `{"name":"pikachu","base_experience":112}`

func stubRoll(t *testing.T, roll int) {
	t.Helper()
	orig := rollCatch
	rollCatch = func(int) int { return roll }
	t.Cleanup(func() { rollCatch = orig })
}

func TestCatchCatchesPokemon(t *testing.T) {
	cfg := newTestConfig()
	cfg.args = []string{"pikachu"}
	cfg.cache.Add("pikachu", []byte(pikachuJSON))

	stubRoll(t, 799)

	var err error
	out := captureStdout(t, func() { err = commandCatch(cfg) })
	if err != nil {
		t.Fatalf("commandCatch: %v", err)
	}
	want := "Throwing a Pokeball at pikachu...\npikachu was caught!\nYou may now inspect it with the inspect command.\n"
	if out != want {
		t.Errorf("expected %q, got %q", want, out)
	}

	entry, ok := cfg.pokedex["pikachu"]
	if !ok {
		t.Fatal("expected pikachu in pokedex after catch")
	}
	if entry.Name != "pikachu" || entry.BaseExperience != 112 {
		t.Errorf("expected entry {pikachu 112}, got %+v", entry)
	}
}

func TestCatchEscapes(t *testing.T) {
	cfg := newTestConfig()
	cfg.args = []string{"pikachu"}
	cfg.cache.Add("pikachu", []byte(pikachuJSON))

	stubRoll(t, 111)

	var err error
	out := captureStdout(t, func() { err = commandCatch(cfg) })
	if err != nil {
		t.Fatalf("commandCatch: %v", err)
	}
	want := "Throwing a Pokeball at pikachu...\npikachu escaped!\n"
	if out != want {
		t.Errorf("expected %q, got %q", want, out)
	}
	if len(cfg.pokedex) != 0 {
		t.Errorf("expected empty pokedex after escape, got %+v", cfg.pokedex)
	}
}

func TestCatchRollBoundary(t *testing.T) {
	cfg := newTestConfig()
	cfg.args = []string{"pikachu"}
	cfg.cache.Add("pikachu", []byte(pikachuJSON))

	stubRoll(t, 112)

	var err error
	captureStdout(t, func() { err = commandCatch(cfg) })
	if err != nil {
		t.Fatalf("commandCatch: %v", err)
	}
	if _, ok := cfg.pokedex["pikachu"]; !ok {
		t.Errorf("expected catch when roll equals base experience")
	}
}

func TestCatchInvalidCachedJSON(t *testing.T) {
	cfg := newTestConfig()
	cfg.args = []string{"pikachu"}
	cfg.cache.Add("pikachu", []byte("not-json"))

	var err error
	out := captureStdout(t, func() { err = commandCatch(cfg) })
	if err == nil {
		t.Fatal("expected error for non-JSON cached body, got nil")
	}
	if !strings.Contains(out, "failed to catch pikachu") {
		t.Errorf("expected failure message in output, got %q", out)
	}
	if len(cfg.pokedex) != 0 {
		t.Errorf("expected empty pokedex after failed catch, got %+v", cfg.pokedex)
	}
}
