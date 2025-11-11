package pokeapi

import (
	"fmt"
	"testing"
)

func TestGetPokeMap(t *testing.T) {
	t.Run(fmt.Sprintf("Test case TestGetPokeMap"), func(t *testing.T) {
		pokemap, err := GetPokeMap("https://pokeapi.co/api/v2/location-area/?offset=0&limit=20")
		if err != nil {
			t.Errorf("problem while generating pokemap: %v", err)
		}
		if pokemap.Results[0].Name != "canalave-city-area" {
			t.Errorf("different map area, the first map area should be: canalave-city-area")
		}
	})
}

func TestAddGetPokemons(t *testing.T) {
	t.Run(fmt.Sprintf("Test case TestAddGet"), func(t *testing.T) {
		pokeArea, err := GetPokemons("https://pokeapi.co/api/v2/location-area/pastoria-city-area/")
		if err != nil {
			t.Errorf("problem while generating pokemap: %v", err)
		}
		if pokeArea.PokemonEncounters[0].Pokemon.Name != "tentacools" {
			t.Errorf("pokemon was not tentacool, it was %v", pokeArea.PokemonEncounters[0].Pokemon.Name)
		}
	})
}
