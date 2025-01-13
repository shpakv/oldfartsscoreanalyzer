package notifier

import "oldfartscounter/internal/teambuilder"

type Notifier interface {
	Notify(team1, team2 teambuilder.Team) error
}
