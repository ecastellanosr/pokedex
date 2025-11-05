package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type pokeMapInfo struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []pokeLocation `json:"results"`
}

type pokeLocation struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func getPokeMap(url string) (pokeMapInfo, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	}
	var pokeMap = pokeMapInfo{}
	res, err := http.Get(url)
	if err != nil {
		return pokeMap, fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&pokeMap); err != nil {
		return pokeMapInfo{}, err
	}
	return pokeMap, nil
}
