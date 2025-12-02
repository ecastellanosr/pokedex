package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PokemonMove struct {
	ID            int              `json:"id"`
	Name          string           `json:"name"`
	Accuracy      int              `json:"accuracy"`
	Effect_chance int              `json:"effect_chance"`
	PP            int              `json:"pp"` //Power points, how many timesthe move can be used
	Priority      int              `json:"priority"`
	Power         int              `json:"power"`
	Damage_class  NamedAPIResource `json:"damage_class"`
	Stat_changes  []MoveStatChange `json:"stat_changes"`
	Target        NamedAPIResource `json:"target"`
	PokeType      NamedAPIResource `json:"type"`
}

type MoveStatChange struct {
	Change int              `json:"change"`
	Stat   NamedAPIResource `json:"stat"`
}

func GetMove(moveName string) (PokemonMove, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/move/%v/", moveName)
	var move = PokemonMove{}
	res, err := http.Get(url) // no default url
	if err != nil {
		return move, fmt.Errorf("error creating request: %w", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&move); err != nil {
		return PokemonMove{}, err
	}
	return move, nil
}

func removeMoveExtraDetails(Moves *[]Move) *[]Move {
	for i, move := range *Moves {
		for _, detail := range move.Details {
			if detail.VersionGroup.Name == "x-y" {
				(*Moves)[i].Details = []MoveDetails{
					detail,
				}
				break
			}
		}
	}
	return Moves
}

func changetoMoves(Moves []Move, level int) map[int]NamedAPIResource {
	removeMoveExtraDetails(&Moves)
	movesMap := map[int]NamedAPIResource{}
	moveSlot := 3 //Slot where the move is in
	for _, move := range Moves {
		for _, detail := range move.Details {
			if detail.AtLevel < level+1 && detail.MoveLearnedMethod.Name == "level-up" {
				if detail.Order != 0 {
					movesMap[detail.Order] = move.Move
				} else {
					movesMap[moveSlot] = move.Move
					moveSlot++
					if moveSlot == 5 {
						moveSlot = 3
					}
				}
			}
		}
	}
	return movesMap
}
