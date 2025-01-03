package playerrating

import (
	"math"
	"oldfartscounter/internal/player"
	"sort"
	"time"
)

var ActualityHoursRange = 14 * 24
var DegradationCoefficient = -0.1

type Calculator struct {
}

func (c *Calculator) ApplyForRange(p player.Player, timeRange time.Time, defaultValue float64, hardcodeValue ...float64) float64 {
	if hardcodeValue != nil {
		return hardcodeValue[0]
	}

	recentGames := c.recentGames(p, timeRange)

	finalScore := 0.0
	weights := []float64{0.5, 0.3, 0.2}
	usedWeights := 0

	finalScore = c.calculateRecentGames(recentGames, usedWeights, weights, finalScore)

	if len(recentGames) < len(weights) {
		finalScore = c.calculateOlderGames(p, timeRange, finalScore)
	}

	if len(p.Games) == 0 {
		finalScore = defaultValue
	}

	return finalScore
}

func (c *Calculator) recentGames(p player.Player, timeRange time.Time) []*player.PlayedGame {
	recentGames := make([]*player.PlayedGame, 0)

	for _, game := range p.Games {
		if game.Date.After(timeRange) {
			recentGames = append(recentGames, game)
		}
	}

	sort.Slice(recentGames, func(i, j int) bool {
		return recentGames[i].Date.After(recentGames[j].Date)
	})
	return recentGames
}

func (c *Calculator) calculateRecentGames(recentGames []*player.PlayedGame, usedWeights int, weights []float64, finalScore float64) float64 {
	for _, recentGame := range recentGames {
		if usedWeights >= len(weights) {
			break
		}

		if time.Since(recentGame.Date).Hours() <= float64(ActualityHoursRange) {
			finalScore += recentGame.Score() * weights[0]
			usedWeights++
		} else {
			finalScore += recentGame.Score() * weights[usedWeights]
			usedWeights++
		}
	}
	return finalScore
}

func (c *Calculator) calculateOlderGames(p player.Player, timeRange time.Time, finalScore float64) float64 {
	oldGames := make([]*player.PlayedGame, 0)
	for _, game := range p.Games {
		if game.Date.Before(timeRange) {
			oldGames = append(oldGames, game)
		}
	}

	for _, oldGame := range oldGames {
		weeksSinceOldGame := time.Since(oldGame.Date).Hours() / (7 * 24)
		degradationFactor := math.Exp(DegradationCoefficient * weeksSinceOldGame)
		finalScore += oldGame.Score() * degradationFactor
	}
	return finalScore
}
