package pokeapi

import (
	"fmt"
	"math"
)

func calculateExp(level int, growthRate GrowthRate) int {
	if level == 1 {
		return 0
	}
	quadratic := float64(level * level * level)
	switch growthRate.Name {
	case "slow":
		return int(math.Floor(((5.0 * quadratic) / 4.0)))
	case "medium":
		return level * level * level //rbny
	case "fast":
		return int(math.Floor(((4.0 * quadratic) / 5.0)))
	case "medium-slow": //"\\frac{6x^3}{5} - 15x^2 + 100x - 140"
		return int(math.Floor(((6.0 * quadratic) / 5.0) - float64(15*level*level) + float64(100*level) - 140.0))
	default:
		return level * level * level
	}
}

func calculateWinExp(pokemon *Pokemon, defeatedPokemon *Pokemon) int {
	baseExperience := float64(defeatedPokemon.BaseExperience)
	levelPokemon := float64(pokemon.Level)                                                         //level pokemon
	levelDefeatedPokemon := float64(defeatedPokemon.Level)                                         //level defeated pokemon
	s := 1.0                                                                                       //if want to add sharedexperience
	base := (((2.0 * levelDefeatedPokemon) + 10.0) / (levelDefeatedPokemon + levelPokemon + 10.0)) //base raised to the power in the exp formula
	modifier := 1.0                                                                                //if want to add modifiers
	xp := ((((baseExperience * levelDefeatedPokemon) / 5.0 * (1.0 / s)) * math.Pow(base, 2.5)) + 1) * modifier
	return int(math.Round(xp))
}

func GiveXP(pokemon *Pokemon, defeatedPokemon *Pokemon) {
	xp := calculateWinExp(pokemon, defeatedPokemon)
	pokemon.Exp += xp
	fmt.Printf("%v gained %v Exp Points!\n", pokemon.Name, xp)
	//2 evs start in 90 exp and 3 in 230 exp
	return
}

type statEV struct {
	name     string
	EV       int
	baseStat int
}

func statsCopy(stats map[string]*statM) map[string]statM {
	newMap := make(map[string]statM)

	// Iterate through the original map
	for key, ptrValue := range stats {
		// Dereference the pointer to get the actual value
		value := *ptrValue
		// Assign the value to the new map
		newMap[key] = value
	}
	return newMap
}

func AddEV(pokemon, defeatedpokemon *Pokemon) {
	EVAmount := 1
	if defeatedpokemon.BaseExperience > 90 {
		EVAmount = 2
	} else if defeatedpokemon.BaseExperience > 230 {
		EVAmount = 3
	}
	stats := statsCopy(defeatedpokemon.Stats)
	bestStat := getBestStat(stats)

	lowerLimitShare := bestStat.baseStat - 10
	delete(stats, bestStat.name)
	secondBest := getBestStat(stats)
	delete(stats, secondBest.name)
	thirdBest := getBestStat(stats)
	if EVAmount == 1 || secondBest.baseStat < lowerLimitShare {
		pokemon.Stats[bestStat.name].effort += EVAmount
		return
	}
	if EVAmount == 2 || thirdBest.baseStat < lowerLimitShare {
		if secondBest.baseStat > lowerLimitShare {
			pokemon.Stats[bestStat.name].effort += (EVAmount - 1)
			pokemon.Stats[secondBest.name].effort += 1
			return
		}
	}
	pokemon.Stats[bestStat.name].effort += 1
	pokemon.Stats[secondBest.name].effort += 1
	pokemon.Stats[thirdBest.name].effort += 1
	return
}

func getBestStat(stats map[string]statM) statEV {
	bestStat := statEV{}
	for _, stat := range stats {
		currentStat := statEV{
			name:     stat.Name,
			baseStat: stat.baseStat,
		}
		if currentStat.baseStat > bestStat.baseStat {
			bestStat = currentStat
		}
	}
	return bestStat
}

func LevelUp(pokemon *Pokemon, defeatedPokemon *Pokemon) {
	nextLevel := pokemon.Level + 1
	nextLevelExp := calculateExp(nextLevel, pokemon.growthRate)
	if pokemon.Exp >= nextLevelExp {
		pokemon.Level += 1
		AddEV(pokemon, defeatedPokemon)
		UpdateStats(pokemon)
	}
}
