package main

import (
	"time"

	cache "github.com/iheb-kraini/pokedex/pokecache"
)

func main() {
	conf := &config{
		commands: map[string]cliCommand{
			"exit": {
				name:        "exit",
				description: "Exit the Pokedex",
				callback:    commandExit,
			},
			"help": {
				name:        "help",
				description: "Exit the Pokedex",
				callback:    command_help,
			},

			"map": {
				name:        "map",
				description: "map",
				callback:    command_map,
			},
			"mapb": {
				name:        "bmap",
				description: "bmap",
				callback:    command_bmap,
			},
			"explore": {
				name:        "explore",
				description: "explore the possibel pokemen to encouter in an area",
				callback:    commandExplore,
			},
			"catch": {
				name:        "catch",
				description: "try to catch a Pokemon",
				callback:    commandCatch,
			},
			"inspect": {
				name:        "inspect",
				description: "details about a caught Pokemon",
				callback:    commandInspect,
			},
		},

		next:    "https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
		prev:    "",
		cache:   cache.NewCache(5 * time.Minute),
		args:    []string{},
		pokedex: map[string]Pokemon{},
	}

	startRepl(conf)
}
