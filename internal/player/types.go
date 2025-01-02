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
	}
)
