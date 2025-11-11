package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	pokeapi "github.com/ecastellanosr/pokedex/internal/pokeAPI"
	"github.com/ecastellanosr/pokedex/internal/pokecache"
)

func cleanInput(text string) []string { //separate words and clean unnecesary spaces
	output := strings.ToLower(text)
	stringList := strings.Fields(output)
	return stringList
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
		pokedex:  map[string]pokeapi.Pokemon{},
	} //configuration with map, cache and pokedex

	for { // main loop
		fmt.Print("Pokedex >")
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

		}
		command, exists := getCommands()[firstLine]
		if exists {
			err := command.callback(&cfg) //execute command
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}
