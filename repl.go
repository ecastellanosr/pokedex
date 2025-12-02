package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	pokeapi "github.com/ecastellanosr/pokedex/internal/pokeAPI"
	pokebattle "github.com/ecastellanosr/pokedex/internal/pokeBattle"
	"github.com/ecastellanosr/pokedex/internal/pokecache"
)

func cleanInput(text string) []string { //separate words and clean unnecesary spaces
	output := strings.ToLower(text)
	stringList := strings.Fields(output)
	return stringList
}

func separateDashedString(text string) []string {
	return strings.Split(text, "-")
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)     //scanner for stdIn
	interval, err := time.ParseDuration("5s") // cache interval
	if err != nil {
		fmt.Printf("error while getting the interval for cache: %v", err)
	}
	cfg := config{
		next:     "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		previous: "",
		cache:    pokecache.NewCache(interval),
		pokedex:  map[string]pokeapi.HPokemon{},
		starters: [3]string{
			"bulbasaur",
			"squirtle",
			"charmander",
		},
		hasStarter:    false,
		currentRegion: "kanto",
		playerinfo:    pokebattle.NewPlayer(),
		areaPokemon:   map[string]bool{},
	} //configuration with map, cache and pokedex
	fmt.Println("Choose one pokemon between these starter pokemon")
	fmt.Println(cfg.starters)
	for {
		fmt.Print("PokedexCLI >")
		if !scanner.Scan() {
			continue
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading input: %s\n", err)
		}
		line := scanner.Text()
		input := cleanInput(line) // user input
		firstLine := input[0]     //command
		if len(input) == 2 {
			cfg.arg = input[1]
		} else if len(input) == 3 {
			cfg.arg = fmt.Sprintf("%v-%v", input[1], input[2])
		} else {
			cfg.arg = ""
		}
		if !cfg.hasStarter {
			introductionLoop(&cfg, firstLine)
			continue
		}
		actionLoop(&cfg, firstLine)
	}
}

func introductionLoop(cfg *config, commandName string) {
	//get a function here for the name and starter
	if commandName == "start" {
		actionLoop(cfg, commandName)
		fmt.Printf("You are right now in the \"%v\" Region, explore it with \"region kanto\"!\n", cfg.currentRegion)
		return
	} else {
		fmt.Println("First, choose one pokemon between these starter pokemon")
		fmt.Println(cfg.starters)
		return
	}
}

func actionLoop(cfg *config, commandName string) {
	//get a function here for the name and starter
	command, exists := getCommands()[commandName]
	if exists {
		err := command.callback(cfg) //execute command
		if err != nil {
			fmt.Println(err)
		}
		return
	} else {
		fmt.Println("Unknown command")
		return
	}

}
