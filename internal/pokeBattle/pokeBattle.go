package pokebattle

import (
	"bufio"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"strings"

	pokeapi "github.com/ecastellanosr/pokedex/internal/pokeAPI"
)

type PlayerInfo struct {
	Team           []*pokeapi.Pokemon
	FirstPokemonID int
	CurrentArea    pokeapi.PokeArea
	Items          map[string]item
}

type item struct {
	Name     string
	Amount   int
	ItemType string
}

type PokeBattleInfo struct {
	playerInfo     *PlayerInfo
	currentPokemon *pokeapi.Pokemon
	RivalPokemon   *pokeapi.Pokemon
}

func cleanInput(text string) []string { //pasted
	output := strings.ToLower(text)
	stringList := strings.Fields(output)
	return stringList
}

func NewPlayer() *PlayerInfo { //create a new cache
	player := PlayerInfo{
		Team:           []*pokeapi.Pokemon{},
		FirstPokemonID: 0,
		CurrentArea:    pokeapi.PokeArea{},
		Items: map[string]item{
			"pokeball": {
				Name:     "pokeball",
				Amount:   5,
				ItemType: "pokeball",
			},
			"greatball": {
				Name:     "greatball",
				Amount:   2,
				ItemType: "pokeball",
			},
		},
	}

	return &player
}

func PokeBattle(player *PlayerInfo, rivalpokemon *pokeapi.Pokemon) bool {
	fmt.Println("[fight] [item] [run]")
	scanner := bufio.NewScanner(os.Stdin) //scanner for stdIn
	pokemon := player.Team[player.FirstPokemonID]
	pokeBattle := PokeBattleInfo{
		playerInfo:     player,
		currentPokemon: pokemon,
		RivalPokemon:   rivalpokemon,
	}
	currentEvent := "base" //event is the place where the battle is at, it can be "base","choosing","attack","item",""
	for {
		fmt.Print("Pokebattle >")

		if !scanner.Scan() {
			continue
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading input: %s\n", err)
		}
		line := scanner.Text()
		input := cleanInput(line) // user input
		firstLine := input[0]     //command
		switch currentEvent {
		case "base":
			run := baseState(*player, *pokemon, firstLine, &currentEvent)
			if run {
				return false
			}
			continue
		case "abilitySelection":
			dead, err := attack(&pokeBattle, firstLine, &currentEvent) //select pokemon attack
			if err != nil {
				fmt.Println(err)
				abilitySelection(*pokemon, &currentEvent)
			}
			if dead {
				if pokeBattle.RivalPokemon.CurrentHealth <= 0 {
					pokeapi.GiveXP(pokemon, rivalpokemon)
					pokeapi.LevelUp(pokemon, rivalpokemon)
					return false
				} else {
					successful := selectNextPokemon(player, &pokeBattle)
					if successful {
						continue
					}
					fmt.Println("all pokemon have fainted", "Fleeing!")
					return false
				}
			}
		case "itemSelection":
			isPokeball := checkPokeball(firstLine)
			if !isPokeball {
				itemSelection(*player, &currentEvent)
				continue
			}
			rivalSpecies, err := pokeapi.GetSpecies(rivalpokemon.Name)
			if err != nil {
				fmt.Println(err)
				currentEvent = "base"
				continue
			}
			currentEvent = "base"
			if pokeapi.CatchPokemon(*rivalpokemon, rivalSpecies, firstLine) {
				return true
			}
			continue
		default:
			run := baseState(*player, *pokemon, firstLine, &currentEvent)
			if run {
				return false
			}
			continue
		}
	}
}

func selectNextPokemon(player *PlayerInfo, pokeBattle *PokeBattleInfo) bool {
	for _, pokemon := range player.Team {
		if pokemon.CurrentHealth > 0 {
			pokeBattle.currentPokemon = pokemon
			return true
		}
	}
	return false
}

func calculateDamage(attacker, defender pokeapi.Pokemon, move pokeapi.PokemonMove) int { //effect modifiere
	pokemonLevel := float64(attacker.Level)
	power := float64(move.Power)
	attack := float64(attacker.Stats["attack"].Amount)
	defense := float64(defender.Stats["defense"].Amount)
	if move.Damage_class.Name == "special" {
		attack = float64(attacker.Stats["special-attack"].Amount)
		defense = float64(defender.Stats["special-defense"].Amount)
	}
	modifiers, err := pokeapi.TypeModifier(defender, move) //past modifierp[1] values dont matter, there will alway be at least 2 modifiers
	if err != nil {
		fmt.Printf("error while getting type modifiers")
	}
	BaseDamage := (((((2.0 * pokemonLevel) / 5.0) + 2.0) * power * (attack / defense)) / 50.0)
	random := (float64(rand.IntN(16)+85) / 100.0)
	fmt.Printf("- pokemonlevel: %v\n- power: %v\n- attack: %v\n- defense: %v\n- basedamage:%v\n - random: %v\n", pokemonLevel, power, attack, defense, BaseDamage, random)
	Damage := BaseDamage * modifiers[0] * modifiers[1] * random
	return int(math.Floor(Damage))
}

func attack(pokeBattleInfo *PokeBattleInfo, attackname string, currentEvent *string) (bool, error) { //
	userMove, err := pokeapi.GetMove(attackname)
	if err != nil {
		return false, err
	}
	rivalMove, err := pokeapi.GetMove(randomMove(*pokeBattleInfo.RivalPokemon))
	if err != nil {
		return false, err
	}
	rivalDamage := calculateDamage(*pokeBattleInfo.RivalPokemon, *pokeBattleInfo.currentPokemon, rivalMove)
	userDamage := calculateDamage(*pokeBattleInfo.currentPokemon, *pokeBattleInfo.RivalPokemon, userMove)

	if pokeBattleInfo.currentPokemon.Stats["speed"].Amount < pokeBattleInfo.RivalPokemon.Stats["speed"].Amount {
		dead := inflictDamage(pokeBattleInfo.RivalPokemon, pokeBattleInfo.currentPokemon, rivalDamage, rivalMove)
		if dead {
			return true, nil
		}
		dead = inflictDamage(pokeBattleInfo.currentPokemon, pokeBattleInfo.RivalPokemon, userDamage, userMove)
		if dead {
			return true, nil
		}
	} else {
		dead := inflictDamage(pokeBattleInfo.RivalPokemon, pokeBattleInfo.currentPokemon, rivalDamage, rivalMove)
		if dead {
			return true, nil
		}
		dead = inflictDamage(pokeBattleInfo.currentPokemon, pokeBattleInfo.RivalPokemon, userDamage, userMove)
		if dead {
			return true, nil
		}
	}

	*currentEvent = "base"
	return false, nil
}

func inflictDamage(attacker, defender *pokeapi.Pokemon, damage int, move pokeapi.PokemonMove) bool {
	died := false
	fmt.Printf("%v used %v!\n", attacker.Name, move.Name)
	defender.CurrentHealth = defender.CurrentHealth - damage
	if rand.IntN(101) > move.Accuracy {
		fmt.Printf("%v's attack Missed!\n", attacker.Name)
		return died
	}
	fmt.Printf("%v received %v of damage!\n", defender.Name, damage)
	if defender.CurrentHealth <= 0 {
		defender.CurrentHealth = 0
		fmt.Printf("%v has fainted!\n", defender.Name)
		died = true
		return died
	}
	return died
}

func randomMove(pokemon pokeapi.Pokemon) string {
	moveSlot := (rand.IntN(len(pokemon.Moves)) + 1)
	return pokemon.Moves[moveSlot].Name
}

func baseState(player PlayerInfo, pokemon pokeapi.Pokemon, command string, state *string) bool {
	switch command {
	case "fight":
		abilitySelection(pokemon, state)
		return false
	case "item":
		itemSelection(player, state)
		return false
	case "run":
		run()
		return true
	default:
		fmt.Println("select between \n[fight] [item] [run]")
		return false
	}
}

func itemSelection(player PlayerInfo, state *string) {
	for _, item := range player.Items {
		fmt.Printf("[%v]", item.Name)
	}
	fmt.Printf("\n")
	*state = "itemSelection"
}

func abilitySelection(pokemon pokeapi.Pokemon, state *string) {
	for _, move := range pokemon.Moves {
		fmt.Printf("[%v]", move.Name)
	}
	fmt.Printf("\n")
	*state = "abilitySelection"
	return
}

func run() {
	fmt.Println("You have successfully fled!")
}

func checkPokeball(input string) bool {
	if input == "pokeball" || input == "greatball" || input == "masterball" || input == "ultraball" {
		return true
	}
	return false
}
