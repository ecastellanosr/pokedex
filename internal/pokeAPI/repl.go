package pokeapi

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type config struct {
	next     string
	previous string
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	stringList := strings.Fields(output)
	return stringList
}

func startRepl() {

	scanner := bufio.NewScanner(os.Stdin)
	cfg := config{
		next:     "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		previous: "",
	}
	for {
		fmt.Print("Pokedex >")
		if !scanner.Scan() {
			continue
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading input: %s\n", err)
		}
		line := scanner.Text()
		firstLine := cleanInput(line)[0]
		command, exists := getCommands()[firstLine]
		if exists {
			err := command.callback(&cfg)
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
