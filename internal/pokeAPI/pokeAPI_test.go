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
		pokeArea, err := GetArea("https://pokeapi.co/api/v2/location-area/pastoria-city-area/")
		if err != nil {
			t.Errorf("problem while generating pokemap: %v", err)
		}
		if pokeArea.PokemonEncounters[0].HPokemon.Name != "tentacools" {
			t.Errorf("pokemon was not tentacool, it was %v", pokeArea.PokemonEncounters[0].HPokemon.Name)
		}
	})
}

func TestChangePokemon(t *testing.T) {
	t.Run(fmt.Sprintf("Test case TestChangePokemon"), func(t *testing.T) {
		rattata, err := GetPokemon("https://pokeapi.co/api/v2/pokemon/rattata/")
		if err != nil {
			t.Errorf("problem while generating pokemap: %v", err)
		}
		newPokemon, err := ChangePokemon(rattata)
		if err != nil {
			t.Errorf("error while changing pokemon %v", err)

		}
		pokemon := Pokemon{
			weight: 35,
			Types: []string{
				"normal",
			},
		}
		rattata.List()
		fmt.Println("---------------------------------------------------------------------------------")
		newPokemon.List()
		if newPokemon.weight != pokemon.weight {
			t.Errorf("weight is wrong, it was %v, and should be %v", newPokemon.weight, pokemon.weight)
		}
		if pokemon.Types[0] != newPokemon.Types[0] {
			t.Errorf("type is wrong, it was %v, and should be %v", newPokemon.Types[0], pokemon.Types[0])
		}
	})
}

func TestLevelUp(t *testing.T) {
	t.Run(fmt.Sprintf("Test case TestLevelUp"), func(t *testing.T) {
		blissey, err := GetPokemon("https://pokeapi.co/api/v2/pokemon/blissey")
		rattata, err := GetPokemon("https://pokeapi.co/api/v2/pokemon/rattata/")
		if err != nil {
			t.Errorf("problem while generating pokemap: %v", err)
		}
		levelUpPokemon, err := ChangePokemon(rattata)
		defeatedPokemon, err := ChangePokemon(blissey)
		if err != nil {
			t.Errorf("error while changing pokemon %v", err)

		}
		currentExp := levelUpPokemon.Exp
		currentEV := levelUpPokemon.Stats["hp"].effort
		GiveXP(&levelUpPokemon, &defeatedPokemon)
		LevelUp(&levelUpPokemon, &defeatedPokemon)

		if currentExp == levelUpPokemon.Exp {
			t.Errorf("currentExp (%v) and pokemon Exp (%v) should not be the same", currentExp, levelUpPokemon.Exp)
		}
		if currentEV == levelUpPokemon.Stats["hp"].effort {
			t.Errorf("currentLevel (%v) and pokemon level (%v) should not be the same", currentEV, levelUpPokemon.Stats["hp"].effort)
		}
	})
}
