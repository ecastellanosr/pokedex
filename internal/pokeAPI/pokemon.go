package pokeapi

import (
	"fmt"
	"math"
	"math/rand/v2"
)

type Pokemon struct { //all structs relate to the pokemon struct which is a simpler struct than the Hpokemon which has all info but no calculation
	Name           string
	nature         Nature
	weight         int
	growthRate     GrowthRate
	BaseExperience int
	Exp            int
	CurrentHealth  int
	Level          int //random 1 to 5
	Stats          map[string]*statM
	Types          []string                 //get name from types
	Moves          map[int]NamedAPIResource // get name and url from move
}

func ChangePokemon(HPokemon HPokemon) (Pokemon, error) {
	nature, err := GetNature()
	if err != nil {
		return Pokemon{}, fmt.Errorf("error while changing pokemon struct: %v", err)
	}
	growthRate, err := GetGrowthRate()
	if err != nil {
		return Pokemon{}, fmt.Errorf("error while changing pokemon struct: %v", err)
	}
	level := (rand.IntN(6) + 1)
	statm := changetoStatM(HPokemon.Stats, nature, level)
	typesM := changetoTypes(HPokemon.Types)
	movesM := changetoMoves(HPokemon.Move, level)
	HP := statm["hp"].Amount
	exp := calculateExp(level, growthRate)
	pokemon := Pokemon{
		Name:           HPokemon.Name,
		nature:         nature,
		growthRate:     growthRate,
		weight:         HPokemon.Weight,
		BaseExperience: HPokemon.BaseExperience,
		Exp:            exp, //change
		CurrentHealth:  HP,
		Level:          level,
		Stats:          statm,
		Types:          typesM,
		Moves:          movesM,
	}

	return pokemon, nil
}

func (p Pokemon) List() error { // list pokemon info
	if len(p.Stats) == 0 {
		return fmt.Errorf("nothing to list")
	}
	fmt.Printf("Name: %v\n", p.Name)
	fmt.Printf("Weight: %v\n", p.weight)
	fmt.Printf("Nature: %v\n", p.nature.Name)
	fmt.Printf("Growth Rate:%v\n", p.growthRate.Name)
	fmt.Printf("Level: %v\n", p.Level)
	fmt.Printf("Current XP: %v\n", p.Exp)
	fmt.Printf("Current Health: %v\n", p.CurrentHealth)
	fmt.Println("Stats:")
	for _, pokemonStats := range p.Stats {
		fmt.Printf("  -%v: %v\n", pokemonStats.Name, pokemonStats.Amount)
	}
	fmt.Println("Types:")
	for _, pokemonTypes := range p.Types {
		fmt.Printf("  - %v\n", pokemonTypes)
	}
	for _, move := range p.Moves {
		fmt.Printf("  - %v\n", move.Name)
	}
	return nil
}

func CatchPokemon(pokemon Pokemon, species Species, pokeballtype string) bool { //function to determine if pokemon gets catched
	pokeballModifier := 1.0
	switch pokeballtype {
	case "greatball":
		pokeballModifier = 1.5
	case "ultraball":
		pokeballModifier = 2.0
	case "masterball":
		return true
	}
	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon.Name)
	//chance is the percentage 0-100 of catching it, the function is based
	//in the base experience and is a sigmoid equation with 98% for 0 and 5% inf
	top := ((3 * pokemon.Stats["hp"].Amount) - (2 * pokemon.CurrentHealth))
	bottom := (3 * pokemon.Stats["hp"].Amount)
	modifiers := int(4096.0 * float64(species.CaptureRate) * pokeballModifier)
	catchrate := (top / bottom) * modifiers
	chancedivision := (catchrate / 1044480.0)
	chance := math.Pow(float64(chancedivision), 0.75)

	random := rand.IntN(101) // a random number [0,100]

	if random > int(math.Round(chance*100.0)) { //if random number is less than chance then it escapes
		fmt.Printf("%v escaped!\n", pokemon.Name)
		return false
	}
	fmt.Printf("%v was caught!\n", pokemon.Name)
	return true
}

func RandomEncounter(pokemonList []string) string {
	random := rand.IntN(len(pokemonList))
	return pokemonList[random]
}

func RandomPokemon(pokemonArea map[string]bool) string {

	for pokemon := range pokemonArea {
		random := rand.IntN(5)
		if random == 2 {
			return pokemon
		}
	}
	return ""
}
