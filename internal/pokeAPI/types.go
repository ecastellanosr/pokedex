package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HPokeType struct {
	Name            string        `json:"name"`
	ID              int           `json:"id"`
	DamageRelations typeRelations `json:"damage_relations"`
}

type typeRelations struct {
	NoDamageTo       []NamedAPIResource `json:"no_damage_to"`
	HalfDamageTo     []NamedAPIResource `json:"half_damage_to"`
	DoubleDamageTo   []NamedAPIResource `json:"double_damage_to"`
	NoDamageFrom     []NamedAPIResource `json:"no_damage_from"`
	HalfDamageFrom   []NamedAPIResource `json:"half_damage_from"`
	DoubleDamageFrom []NamedAPIResource `json:"double_damage_from"`
}

type PokeType struct {
	Name            string
	ID              int
	DamageRelations []typeRelation
}

type typeRelation struct {
	name string
}

func GetType(typeName string) (HPokeType, error) { //get pokemon struct
	var pokeType = HPokeType{}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/type/%v/", typeName)
	res, err := http.Get(url) // no default url
	if err != nil {
		return pokeType, fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&pokeType); err != nil {
		return HPokeType{}, err
	}
	return pokeType, nil
}

func changetoTypes(Types []types) []string {
	types := []string{}
	for _, poketype := range Types {
		types = append(types, poketype.PokemonType.Name)
	}
	return types
}

func TypeModifier(defender Pokemon, move PokemonMove) ([]float64, error) {
	typeModifier := []float64{}
	moveType := move.PokeType.Name
	pokeType, err := GetType(moveType)
	if err != nil {
		return []float64{1.0, 1.0}, err
	}
	for _, defendertype := range defender.Types {
		for _, DoubleType := range pokeType.DamageRelations.DoubleDamageTo {
			if DoubleType.Name == defendertype {
				typeModifier = append(typeModifier, 2.0)
			}
		}
		for _, halfDamage := range pokeType.DamageRelations.HalfDamageTo {
			if halfDamage.Name == defendertype {
				typeModifier = append(typeModifier, 0.5)
			}
		}
		for _, DoubleType := range pokeType.DamageRelations.NoDamageTo {
			if DoubleType.Name == defendertype {
				typeModifier = append(typeModifier, 0.0)
			}
		}
	}
	if len(typeModifier) < 2 {
		typeModifier = append(typeModifier, 1.0)
		typeModifier = append(typeModifier, 1.0)
	}
	return typeModifier, nil
}
