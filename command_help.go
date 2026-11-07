package main 
import "fmt"

func command_help(c *config) error {
	message :=`Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex
catch: Try to catch a Pokemon and add it to your Pokedex
inspect: Details about a caught Pokemon`
	fmt.Println(message);
	return nil;
}
