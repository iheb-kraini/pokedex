package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

// Pokemon is a species from the /pokemon endpoint, used as a pokedex entry.
type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

const maxBaseExp = 800

// rollCatch is the randomness seam for catching; tests stub it for deterministic outcomes.
var rollCatch = func(n int) int { return rand.Intn(n) }

func getPokemon(c *config, name string) (Pokemon, error) {
	var p Pokemon

	data, ok := c.cache.Get(name)
	if !ok {
		url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", name)
		res, err := http.Get(url)
		if err != nil {
			return p, err
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return p, err
		}
		if res.StatusCode > 299 {
			return p, fmt.Errorf("response failed with status %d: %s", res.StatusCode, data)
		}
		c.cache.Add(name, data)
	}

	err := json.Unmarshal(data, &p)
	return p, err
}

func commandCatch(c *config) error {
	name := c.args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	p, err := getPokemon(c, name)
	if err != nil {
		fmt.Printf("failed to catch %s: %v\n", name, err)
		return err
	}

	if rollCatch(maxBaseExp) >= p.BaseExperience {
		c.pokedex[name] = p
		fmt.Printf("%s was caught!\n", name)
	} else {
		fmt.Printf("%s escaped!\n", name)
	}
	return nil
}
