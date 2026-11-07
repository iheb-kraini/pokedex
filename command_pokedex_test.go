package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestPokedexListsCaught(t *testing.T) {
	cfg := newTestConfig()
	cfg.pokedex["pidgey"] = Pokemon{Name: "pidgey", BaseExperience: 50}
	cfg.pokedex["caterpie"] = Pokemon{Name: "caterpie", BaseExperience: 39}

	var err error
	out := captureStdout(t, func() { err = commandPokedex(cfg) })
	if err != nil {
		t.Fatalf("commandPokedex: %v", err)
	}
	if !strings.Contains(out, "Your Pokedex:\n") {
		t.Errorf("expected header in output, got %q", out)
	}
	for _, name := range []string{"pidgey", "caterpie"} {
		if !strings.Contains(out, fmt.Sprintf(" - %s\n", name)) {
			t.Errorf("expected bullet for %q in output, got %q", name, out)
		}
	}
	if n := strings.Count(out, " - "); n != 2 {
		t.Errorf("expected exactly 2 bullets, got %d in %q", n, out)
	}
}

func TestPokedexEmpty(t *testing.T) {
	cfg := newTestConfig()

	var err error
	out := captureStdout(t, func() { err = commandPokedex(cfg) })
	if err != nil {
		t.Fatalf("commandPokedex: %v", err)
	}
	want := "your pokedex is empty\n"
	if out != want {
		t.Errorf("expected %q, got %q", want, out)
	}
}
