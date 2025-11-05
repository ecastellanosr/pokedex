package main

import (
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Map the next 20 Pokemon locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Map the previous 20 Pokemon locations",
			callback:    commandMapB,
		},
	}
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Printf("\nWelcome to the Pokedex!\nUsage:\n\n")
	for name, command := range getCommands() {
		fmt.Printf("%v: %v\n", name, command.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	pokeMap, err := getPokeMap(cfg.next)
	if err != nil {
		return err
	}
	cfg.next = pokeMap.Next
	cfg.previous = pokeMap.Previous
	for _, location := range pokeMap.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapB(cfg *config) error {
	if cfg.previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	pokeMap, err := getPokeMap(cfg.previous)
	if err != nil {
		return err
	}
	cfg.next = pokeMap.Next
	cfg.previous = pokeMap.Previous
	for _, location := range pokeMap.Results {
		fmt.Println(location.Name)
	}
	return nil
}
