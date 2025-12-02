package main

import (
	"bytes"
	"encoding/gob"
	"fmt"

	pokeapi "github.com/ecastellanosr/pokedex/internal/pokeAPI"
	pokebattle "github.com/ecastellanosr/pokedex/internal/pokeBattle"
	"github.com/ecastellanosr/pokedex/internal/pokecache"
)

// config has the next and previous maps, cache, current extra arg for commands and pokedex
type config struct {
	hasStarter    bool
	next          string
	previous      string
	currentRegion string
	arg           string
	starters      [3]string
	cache         *pokecache.Cache
	pokedex       map[string]pokeapi.HPokemon
	playerinfo    *pokebattle.PlayerInfo
	areaPokemon   map[string]bool
}

func cacheShowList(cfg *config, url string, cacheType string) (bool, error) {
	cache, err := getcache(cfg, url)      //check if there is cache
	_, ok := err.(pokecache.NoCacheError) // if there is no cache return false
	if ok {
		return false, nil
	} else if err != nil {
		return false, err
	}

	err = decodeAndList(cache, cacheType) // decode cache and list it
	if err != nil {
		return false, err
	}
	return true, nil
}

func getcache(cfg *config, name string) ([]byte, error) { //get cache
	item, ok := cfg.cache.Get(name)
	if !ok {
		return []byte{}, pokecache.NoCacheError{
			Url: name,
		}
	}
	return item, nil
}

func decodeAndList(cacheVal []byte, cacheType string) error { //decode cache and list it
	var list pokeapi.List
	buf := bytes.NewBuffer(cacheVal)
	dec := gob.NewDecoder(buf)
	if cacheType == "map" {
		var pokeM pokeapi.PokeMapInfo
		err := dec.Decode(&pokeM) // decode into map
		if err != nil {
			return fmt.Errorf("error while decoding into PokeMapInfo: %v", err)
		}
		list = pokeM

	} else if cacheType == "explore" {
		var pokeArea pokeapi.PokeArea
		err := dec.Decode(&pokeArea) // decode into area
		if err != nil {
			return err
		}
		list = pokeArea
	} else if cacheType == "argmap" {
		var pokeLocation pokeapi.PokeArea
		err := dec.Decode(&pokeLocation) // decode into area
		if err != nil {
			return err
		}
		list = pokeLocation
	}
	err := list.List() // list it
	if err != nil {
		return err
	}
	return nil
}

func cfgUpdate(cfg *config, pokeStruct pokeapi.List, url string) error { //update config
	cfg.arg = ""                 //reset argument
	var data bytes.Buffer        //
	enc := gob.NewEncoder(&data) // encode new struct
	if err := enc.Encode(pokeStruct); err != nil {
		return err
	}
	cfg.cache.Add(url, data.Bytes())                // add struct to cache
	pokeMap, ok := pokeStruct.(pokeapi.PokeMapInfo) // if map update next and previous
	if ok {
		cfg.previous = pokeMap.Previous
		cfg.next = pokeMap.Next
		return nil
	} else {
		return nil
	}
}

func pokedexUpdate(cfg *config, pokemon pokeapi.HPokemon) error { // add pokemon to pokedex
	cfg.pokedex[pokemon.Name] = pokemon
	return nil
}

func teamUpdate(cfg *config, pokemon *pokeapi.Pokemon) {
	if len(cfg.playerinfo.Team) != 6 {
		cfg.playerinfo.Team = append(cfg.playerinfo.Team, pokemon) // change hpokemon to pokemon
	}
}
