package player

import (
	"oldfartscounter/internal/maps"
	"time"
)

type (
	Player struct {
		Person
		Games []*PlayedGame
	}

	Person struct {
		FirstName string
		LastName  string
	}

	PlayedGame struct {
		Date time.Time
		Maps []*PlayedMap
	}
	PlayedMap struct {
		Map         maps.Map
		Kills       int
		Deaths      int
		Assists     int
		TotalDamage int
		Score       float64
	}
)

func (g *PlayedGame) Score() float64 {
	totalScore := 0.0
	if len(g.Maps) > 0 {
		for _, playedMap := range g.Maps {
			totalScore += playedMap.Score
		}
		totalScore = totalScore / float64(len(g.Maps))
	}
	return totalScore
}
