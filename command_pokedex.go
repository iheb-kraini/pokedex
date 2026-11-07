package main

import (
	"fmt"
)

func commandPokedex(c *config) error {
	if len(c.pokedex) == 0 {
		fmt.Println("your pokedex is empty")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name := range c.pokedex {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}
