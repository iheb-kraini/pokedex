package main 
import "fmt"

func command_help(c *config) error {
	message :=`Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex`
	fmt.Println(message);
	return nil;
}
