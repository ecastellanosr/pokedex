package pokeapi

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
)

type List interface {
	List() error
}

type PokeRegion struct { //Location in map
	ID        int                `json:"id"`
	Name      string             `json:"name"`
	Locations []NamedAPIResource `json:"locations"`
}

type PokeMapInfo struct { //Map struct
	Count    int                `json:"count"`
	Next     string             `json:"next"`
	Previous string             `json:"previous"`
	Results  []NamedAPIResource `json:"results"`
}

type PokeLocation struct {
	ID     int                `json:"id"`
	Name   string             `json:"name"`
	Region NamedAPIResource   `json:"region"`
	Areas  []NamedAPIResource `json:"areas"`
}

type NamedAPIResource struct { //general named resource
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokeArea struct { //Location in map
	ID                int                 `json:"id"`
	Name              string              `json:"name"`
	Location          NamedAPIResource    `json:"results"`
	PokemonEncounters []PokemonEncounters `json:"pokemon_encounters"`
}

type PokemonEncounters struct {
	Pokemon NamedAPIResource `json:"pokemon"`
}

type Pokemon struct {
	Name           string  `json:"name"`
	Url            string  `json:"url"`
	BaseExperience int     `json:"base_experience"`
	Height         int     `json:"Height"`
	Weight         int     `json:"Weight"`
	Order          int     `json:"order"`
	Stats          []stat  `json:"stats"`
	Types          []types `json:"types"`
}

type stat struct {
	Stat     NamedAPIResource `json:"stat"`
	Effort   int              `json:"effort"`
	BaseStat int              `json:"base_stat"`
}

type types struct {
	Slot        int              `json:"slot"`
	PokemonType NamedAPIResource `json:"type"`
}

func GetRegions() (PokeMapInfo, error) {
	var pokeMapInfo PokeMapInfo
	url := "https://pokeapi.co/api/v2/region/"
	res, err := http.Get(url) //HTTP get to API
	if err != nil {
		return pokeMapInfo, fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&pokeMapInfo); err != nil {
		return PokeMapInfo{}, err
	}
	return pokeMapInfo, nil
}

func GetRegion(url string) (PokeRegion, error) {
	var pokeRegion = PokeRegion{}
	res, err := http.Get(url) //HTTP get to API
	if err != nil {
		return pokeRegion, fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&pokeRegion); err != nil {
		return PokeRegion{}, err
	}
	return pokeRegion, nil
}

func GetLocation(url string) (PokeLocation, error) {
	var pokeLocation = PokeLocation{}
	res, err := http.Get(url) //HTTP get to API
	if err != nil {
		return pokeLocation, fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&pokeLocation); err != nil {
		return PokeLocation{}, err
	}
	return pokeLocation, nil
}

func GetPokeMap(url string) (PokeMapInfo, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	} // default URL
	var pokeMap = PokeMapInfo{}
	res, err := http.Get(url) //HTTP get to API
	if err != nil {
		return pokeMap, fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&pokeMap); err != nil {
		return PokeMapInfo{}, err
	}
	return pokeMap, nil
}

func GetPokemons(url string) (PokeArea, error) {
	var pokeArea = PokeArea{}
	res, err := http.Get(url) // no default url
	if err != nil {
		return pokeArea, fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&pokeArea); err != nil {
		return PokeArea{}, err
	}
	return pokeArea, nil
}

func GetPokemon(url string) (Pokemon, error) { //get pokemon struct
	var pokemon = Pokemon{}
	res, err := http.Get(url) // no default url
	if err != nil {
		return pokemon, fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&pokemon); err != nil {
		return Pokemon{}, err
	}
	return pokemon, nil
}

func CatchPokemon(pokemon Pokemon) bool { //function to determine if pokemon gets catched
	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon.Name)
	//chance is the percentage 0-100 of catching it, the function is based
	//in the base experience and is a sigmoid equation with 98% for 0 and 5% inf
	chance := 100.0 - (95.0 / (1.0 + math.Exp(-0.03*(float64(pokemon.BaseExperience)-130.0))))
	random := rand.IntN(101) // a random number [0,100]
	fmt.Printf("chance: %v,random: %v\n", chance, random)
	if random > int(math.Round(chance)) { //if random number is less than chance then it escapes
		fmt.Printf("%v escaped!\n", pokemon.Name)
		return false
	}
	fmt.Printf("%v was caught!\n", pokemon.Name)
	return true
}

func (p PokeMapInfo) List() error { //list locations in map
	if len(p.Results) == 0 {
		return fmt.Errorf("nothing to list")
	}
	for _, location := range p.Results {
		fmt.Printf(" - %v\n", location.Name)
	}
	return nil
}

func (p PokeLocation) List() error { //list locations in map
	if len(p.Areas) == 0 {
		return fmt.Errorf("nothing to list")
	}
	fmt.Printf("Areas in \"%v\" location:\n", p.Name)
	for _, location := range p.Areas {
		fmt.Printf(" - %v\n", location.Name)
	}
	fmt.Println("[Use \"explore [area]\" to enter a location]")
	return nil
}

func (p PokeArea) List() error { // list pokemon names in location
	if len(p.PokemonEncounters) == 0 {
		return fmt.Errorf("nothing to list")
	}
	for _, pokemonEncounter := range p.PokemonEncounters {
		fmt.Println(pokemonEncounter.Pokemon.Name)
	}
	return nil
}

func (p PokeRegion) List() error { // list pokemon names in location
	if len(p.Locations) == 0 {
		return fmt.Errorf("nothing to list")
	}
	fmt.Printf("Locations in \"%v\" region:\n", p.Name)

	for _, location := range p.Locations {
		fmt.Printf("- %v\n", location.Name)
	}
	fmt.Println("[Use \"travel [Region Name]\" to change regions]") //
	fmt.Println("[Use \"map [location]\" to enter a location]")
	return nil
}

func (p Pokemon) List() error { // list pokemon info
	if len(p.Stats) == 0 {
		return fmt.Errorf("nothing to list")
	}
	fmt.Printf("Name: %v\n", p.Name)
	fmt.Printf("Height: %v\n", p.Height)
	fmt.Printf("Weight: %v\n", p.Weight)
	fmt.Println("Stats:")
	for _, pokemonStats := range p.Stats {
		fmt.Printf("  -%v: %v\n", pokemonStats.Stat.Name, pokemonStats.BaseStat)
	}
	fmt.Println("Types:")
	for _, pokemonStats := range p.Types {
		fmt.Printf("  - %v\n", pokemonStats.PokemonType.Name)
	}
	return nil
}
