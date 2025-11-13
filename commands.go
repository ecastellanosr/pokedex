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
		"region": {
			name:        "region",
			description: "Map location in a region, needs the name as an argument",
			callback:    commandRegion,
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
		"start": {
			name:        "start",
			description: "Pick a starter pokemon",
			callback:    commandStart,
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

func commandStart(cfg *config) error { //catch pokemon after going to an area
	if cfg.hasStarter {
		fmt.Println("You already have a starter!")
		return nil
	}
	pokemonURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v/", cfg.arg)
	for _, pokemonName := range cfg.starters {
		if pokemonName == cfg.arg {
			cfg.hasStarter = true
			pokemon, err := pokeapi.GetPokemon(pokemonURL)
			if err != nil {
				return fmt.Errorf("Error while getting pokemon: %v", err)
			}
			err = pokedexUpdate(cfg, pokemon)
			if err != nil {
				return fmt.Errorf("Error while updating pokedex: %v", err)
			}
			fmt.Printf("%v has been selected\n", pokemon.Name)
			fmt.Println()
			return nil
		} else {
			continue
		}
	}
	fmt.Printf("%v is not a starter, please select a starter from\n%v\n", cfg.arg, cfg.starters)
	cfg.arg = ""
	return nil
}

// Map Region location or list all regions
func commandRegion(cfg *config) error {
	if cfg.arg == "" {
		regions, err := pokeapi.GetRegions()
		if err != nil {
			return err
		}
		if err = regions.List(); err != nil { //list pokemons in area
			return fmt.Errorf("Error while Listing Regions: %v", err)
		}
		return nil
	}
	region := fmt.Sprintf("https://pokeapi.co/api/v2/region/%v/", cfg.arg) //place URL
	ok, err := cacheShowList(cfg, cfg.arg, "region")                       //Check cache
	if err != nil {
		return fmt.Errorf("Error while cacheShowList: %v", err)
	}
	if ok {
		return nil
	}
	pokeRegion, err := pokeapi.GetRegion(region) //Get pokemons in the area
	if err != nil {
		return fmt.Errorf("Error while GetPokeMap: %v", err)
	}
	if err = cfgUpdate(cfg, pokeRegion, cfg.arg); err != nil { //update cfg for cache
		return fmt.Errorf("Error while cfgUpdate: %v", err)
	}

	if err = pokeRegion.List(); err != nil { //list pokemons in area
		return fmt.Errorf("Error while Listing Locations in a region: %v", err)
	}
	return nil
}

// Map Areas in a location
func commandMap(cfg *config) error {
	if cfg.arg != "" {
		err := commandArgMap(cfg)
		if err != nil {
			return fmt.Errorf("error in commandArgMap: %v", err)
		}
		return nil
	}
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

func commandArgMap(cfg *config) error {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location/%v/", cfg.arg)
	ok, err := cacheShowList(cfg, url, "argmap") //Check if area is in the cache for easy access
	if err != nil {
		return fmt.Errorf("error in cacheshowlist: %v", err)
	}
	if ok {
		return nil
	}
	pokeLocation, err := pokeapi.GetLocation(url) // Get the map
	if err != nil {
		return fmt.Errorf("pokelocation: %v", err)
	}
	if err = cfgUpdate(cfg, pokeLocation, url); err != nil { //update configuration for next and prev map
		return fmt.Errorf("cfgUpdate: %v", err)
	}
	if err = pokeLocation.List(); err != nil { //List the map areas
		return fmt.Errorf("list error: %v", err)
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
