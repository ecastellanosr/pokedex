package main

import (
	"fmt"
	"os"

	pokeapi "github.com/ecastellanosr/pokedex/internal/pokeAPI"
)

// CLI Pokedex command format
type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

// get all current commands
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
		"explore": {
			name:        "explore",
			description: "Explore an specific location; needs the name as an argument",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a specific pokemon, needs the name as an argument",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "List a pokemon's traits in your pokedex, needs the name as an argument",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List your pokemon in your pokedex",
			callback:    commandPokedex,
		},
	}
}

// exit CLI
func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// List all commands
func commandHelp(cfg *config) error {
	fmt.Printf("\nWelcome to the Pokedex!\nUsage:\n\n")
	for name, command := range getCommands() {
		fmt.Printf("%v: %v\n", name, command.description)
	}
	return nil
}

// Map next area
func commandMap(cfg *config) error {
	url := cfg.next
	ok, err := cacheShowList(cfg, url, "map") //Check if area is in the cache for easy access
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	pokeMap, err := pokeapi.GetPokeMap(url) // Get the map
	if err != nil {
		return err
	}
	if err = cfgUpdate(cfg, pokeMap, url); err != nil { //update configuration for next and prev map
		return err
	}
	if err = pokeMap.List(); err != nil { //List the map areas
		return err
	}
	return nil
}

// Map previous area
func commandMapB(cfg *config) error { //same as commandMap
	url := cfg.previous
	if url == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	ok, err := cacheShowList(cfg, url, "map")
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	pokeMap, err := pokeapi.GetPokeMap(url)
	if err != nil {
		return err
	}
	if err = cfgUpdate(cfg, pokeMap, url); err != nil {
		return err
	}
	if err = pokeMap.List(); err != nil {
		return err
	}
	return nil
}

// explore a especific location
func commandExplore(cfg *config) error {
	place := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v/", cfg.arg) //place URL
	ok, err := cacheShowList(cfg, place, "explore")                              //Check cache
	if err != nil {
		return fmt.Errorf("Error while cacheShowList: %v", err)
	}
	if ok {
		return nil
	}
	pokeArea, err := pokeapi.GetPokemons(place) //Get pokemons in the area
	if err != nil {
		return fmt.Errorf("Error while GetPokeMap: %v", err)
	}
	if err = cfgUpdate(cfg, pokeArea, place); err != nil { //update cfg for cache
		return fmt.Errorf("Error while cfgUpdate: %v", err)
	}

	if err = pokeArea.List(); err != nil { //list pokemons in area
		return fmt.Errorf("Error while Listing Area: %v", err)
	}
	return nil
}

func commandCatch(cfg *config) error { //catch pokemon after going to an area
	cachePokemon, ok := cfg.pokedex[cfg.arg]
	if ok {
		cfg.arg = ""
		pokeapi.CatchPokemon(cachePokemon)
		return nil
	}
	pokemonURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v/", cfg.arg)
	pokemon, err := pokeapi.GetPokemon(pokemonURL)
	if err != nil {
		return fmt.Errorf("Error while getting pokemon: %v", err)
	}
	catched := pokeapi.CatchPokemon(pokemon)
	if catched == false {
		return nil
	}
	err = pokedexUpdate(cfg, pokemon)
	if err != nil {
		return fmt.Errorf("Error while updating pokedex: %v", err)
	}
	return nil
}

func commandInspect(cfg *config) error { //inspect a pokemon in your pokedex
	cachePokemon, ok := cfg.pokedex[cfg.arg] //check for pokemon in pokedex
	cfg.arg = ""
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	err := cachePokemon.List() //Show pokemon info
	if err != nil {
		return fmt.Errorf("error while listing pokemon traits: %v", err)
	}
	return nil
}

func commandPokedex(cfg *config) error { //list pokemons in pokedex
	if len(cfg.pokedex) == 0 {
		fmt.Println("No pokemon in pokedex")
		return nil
	}
	fmt.Println("Your pokedex:")
	for _, pokemon := range cfg.pokedex {
		fmt.Printf(" - %v\n", pokemon.Name)
	}
	return nil
}
