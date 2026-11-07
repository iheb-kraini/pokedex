package main

import (
	"strings"
	"testing"
)

const pidgeyJSON = `{"name":"pidgey","height":3,"weight":18,"stats":[` +
	`{"base_stat":40,"stat":{"name":"hp"}},` +
	`{"base_stat":45,"stat":{"name":"attack"}},` +
	`{"base_stat":40,"stat":{"name":"defense"}},` +
	`{"base_stat":35,"stat":{"name":"special-attack"}},` +
	`{"base_stat":35,"stat":{"name":"special-defense"}},` +
	`{"base_stat":56,"stat":{"name":"speed"}}],` +
	`"types":[{"type":{"name":"normal"}},{"type":{"name":"flying"}}]}`

func TestInspectCaughtPokemon(t *testing.T) {
	cfg := newTestConfig()
	cfg.args = []string{"pidgey"}
	cfg.cache.Add("pidgey", []byte(pidgeyJSON))
	p, err := getPokemon(cfg, "pidgey")
	if err != nil {
		t.Fatalf("getPokemon: %v", err)
	}
	cfg.pokedex["pidgey"] = p

	var err2 error
	out := captureStdout(t, func() { err2 = commandInspect(cfg) })
	if err2 != nil {
		t.Fatalf("commandInspect: %v", err2)
	}
	want := "Name: pidgey\n" +
		"Height: 3\n" +
		"Weight: 18\n" +
		"Stats:\n" +
		"  -hp: 40\n" +
		"  -attack: 45\n" +
		"  -defense: 40\n" +
		"  -special-attack: 35\n" +
		"  -special-defense: 35\n" +
		"  -speed: 56\n" +
		"Types:\n" +
		"  - normal\n" +
		"  - flying\n"
	if out != want {
		t.Errorf("expected inspect output %q, got %q", want, out)
	}
}

func TestInspectNotCaught(t *testing.T) {
	cfg := newTestConfig()
	cfg.args = []string{"pidgey"}

	var err error
	out := captureStdout(t, func() { err = commandInspect(cfg) })
	if err != nil {
		t.Fatalf("commandInspect: %v", err)
	}
	want := "you have not caught that pokemon\n"
	if out != want {
		t.Errorf("expected %q, got %q", want, out)
	}
}

func TestCatchThenInspect(t *testing.T) {
	cfg := newTestConfig()
	cfg.cache.Add("pidgey", []byte(pidgeyJSON))

	stubRoll(t, 799)
	cfg.args = []string{"pidgey"}
	if err := commandCatch(cfg); err != nil {
		t.Fatalf("commandCatch: %v", err)
	}

	var err error
	out := captureStdout(t, func() { err = commandInspect(cfg) })
	if err != nil {
		t.Fatalf("commandInspect: %v", err)
	}
	if !strings.Contains(out, "Height: 3") || !strings.Contains(out, "- flying") {
		t.Errorf("expected inspect details after catch, got %q", out)
	}
}
