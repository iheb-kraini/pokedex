package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	cache "github.com/iheb-kraini/pokedex/pokecache"
)

type config struct {
	commands map[string]cliCommand
	prev     string
	next     string
	cache    cache.Cache
}
type cliCommand struct {
	name        string
	description string
	callback    func(c *config) error
}

func cleanInput(text string) []string {
	words := strings.Split(text, " ")
	out := make([]string, 0)
	for _, word := range words {
		if len(word) != 0 {
			out = append(out, strings.ToLower(word))
		}
	}
	return out
}
func startRepl(c *config) {
	c_map := c.commands
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() { // Reads up to a newline character
			input := cleanInput(scanner.Text()) // Converts the input bytes into a string
			val, ok := c_map[input[0]]
			if ok {
				val.callback(c)
			}
		}
	}
}
