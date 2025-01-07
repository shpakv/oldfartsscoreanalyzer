package notifier

import "oldfartscounter/internal/teambuilder"

type Notify interface {
	NotifyOldFarts(team1, team2 teambuilder.Team) error
}
