package pokeapi

import (
	"encoding/json"
	"fmt"
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
	HPokemon NamedAPIResource `json:"pokemon"`
}

type HPokemon struct {
	Name           string  `json:"name"`
	Url            string  `json:"url"`
	BaseExperience int     `json:"base_experience"`
	Height         int     `json:"Height"`
	Weight         int     `json:"Weight"`
	Order          int     `json:"order"`
	Stats          []stat  `json:"stats"`
	Types          []types `json:"types"`
	Move           []Move  `json:"moves"`
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

type Move struct {
	Move    NamedAPIResource `json:"move"`
	Details []MoveDetails    `json:"version_group_details"`
}

type MoveDetails struct {
	MoveLearnedMethod NamedAPIResource `json:"move_learn_method"`
	VersionGroup      NamedAPIResource `json:"version_group"`
	AtLevel           int              `json:"level_learned_at"`
	Order             int              `json:"order"` //Order by which the pokemon will learn the move. A newly learnt move will replace the move with lowest order.
}

type GrowthRate struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Nature struct {
	ID             int              `json:"id"`
	Name           string           `json:"name"`
	Decreased_stat NamedAPIResource `json:"decreased_stat"`
	Increased_stat NamedAPIResource `json:"increased_stat"`
}

type Species struct {
	Name        string `json:"name"`
	CaptureRate int    `json:"capture_rate"`
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

func GetArea(url string) (PokeArea, error) {
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

func GetPokemon(url string) (HPokemon, error) { //get pokemon struct
	var pokemon = HPokemon{}
	res, err := http.Get(url) // no default url
	if err != nil {
		return pokemon, fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&pokemon); err != nil {
		return HPokemon{}, err
	}
	return pokemon, nil
}

func GetSpecies(pokemon string) (Species, error) { //get pokemon struct
	var species = Species{}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon-species/%v/", pokemon)
	res, err := http.Get(url) // no default url
	if err != nil {
		return species, fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&pokemon); err != nil {
		return Species{}, err
	}
	return species, nil
}

func GetNature() (Nature, error) {
	random := (rand.IntN(25) + 1)
	url := fmt.Sprintf("https://pokeapi.co/api/v2/nature/%v/", random)
	var nature = Nature{}
	res, err := http.Get(url) // no default url
	if err != nil {
		return nature, fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&nature); err != nil {
		return Nature{}, fmt.Errorf("error with nature: %v", err)
	}
	return nature, nil
}

func GetGrowthRate() (GrowthRate, error) {
	random := (rand.IntN(4) + 1)
	url := fmt.Sprintf("https://pokeapi.co/api/v2/growth-rate/%v/", random)
	var growthRate = GrowthRate{}
	res, err := http.Get(url) // no default url
	if err != nil {
		return growthRate, fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&growthRate); err != nil {
		return GrowthRate{}, err
	}
	return growthRate, nil
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
		fmt.Println(pokemonEncounter.HPokemon.Name)
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

func (p HPokemon) List() error { // list pokemon info
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
