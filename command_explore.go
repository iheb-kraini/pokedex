package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type pokemon struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
type foundRes struct {
	PokemonEncounters []struct {
		Pokemon pokemon `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func getExplore(c *config, mapAreaKey string) (foundRes, error) {
	var apiRes foundRes

	data, ok := c.cache.Get(mapAreaKey)
	if !ok {
		url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", mapAreaKey)
		res, err := http.Get(url)
		if err != nil {
			return apiRes, err
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return apiRes, err
		}
		if res.StatusCode > 299 {
			return apiRes, fmt.Errorf("response failed with status %d: %s", res.StatusCode, data)
		}
		c.cache.Add(mapAreaKey, data)
	}

	err := json.Unmarshal(data, &apiRes)
	return apiRes, err
}

func showPokemon(c *config, mapAreaKey string) error {
	apiRes, err := getExplore(c, mapAreaKey)
	if err != nil {
		return err
	}
	for _, item := range apiRes.PokemonEncounters {
		fmt.Println(item.Pokemon.Name)
	}
	return nil
}

func commandExplore(c *config) error {
	if len(c.args) == 0 {
		return nil
	}
	return showPokemon(c, c.args[0])
}
