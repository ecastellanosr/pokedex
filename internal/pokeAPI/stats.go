package pokeapi

import (
	"math"
	"math/rand/v2"
)

type statM struct {
	Name     string
	Amount   int
	effort   int
	baseStat int
	iv       int
}

func changetoStatM(HStats []stat, nature Nature, level int) map[string]*statM {
	stats := map[string]*statM{}
	for _, stat := range HStats {
		IV := rand.IntN(32)
		stats[stat.Stat.Name] = &statM{
			Name:     stat.Stat.Name,
			Amount:   calculateStatAmount(stat.Stat.Name, stat.Effort, stat.BaseStat, IV, level, nature),
			effort:   stat.Effort,
			baseStat: stat.BaseStat,
			iv:       IV,
		}
	}
	return stats
}

func calculateStatAmount(name string, Effort, BaseStat, IV, Level int, nature Nature) int {
	basestat := float64(BaseStat)
	IVfloat := float64(IV)
	effort := float64(Effort)
	levelfloat := float64(Level)
	modifierEquation := math.Floor((((((2.0 * basestat) + IVfloat) + (effort / 4.0)) * levelfloat) / 100.0))
	if name == "hp" {
		amount := (modifierEquation + levelfloat + 10.0)
		return int(math.Round(amount))
	}
	natureModifier := 1.0
	if nature.Decreased_stat.Name == name {
		natureModifier = 0.9
	} else if nature.Increased_stat.Name == name {
		natureModifier = 1.1
	}
	amount := ((modifierEquation + 5.0) * natureModifier)

	return int(math.Round(amount))
}
func UpdateStats(pokemon *Pokemon) {
	for name, stat := range pokemon.Stats {
		pokemon.Stats[name].Amount = calculateStatAmount(stat.Name, stat.effort, stat.baseStat, stat.iv, pokemon.Level, pokemon.nature)
	}
}
