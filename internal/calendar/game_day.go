package calendar

import "time"

type GameDayType string

const (
	GameDay  GameDayType = "game"
	ThinkDay GameDayType = "think"
)

type GameDayOption struct {
	Type GameDayType
	Date time.Time
}
