package notifier

import (
	"fmt"
	"oldfartscounter/internal/teambuilder"
	"oldfartscounter/internal/teamtable"
)

type Notifier interface {
	Notify(team1, team2 teambuilder.Team, sorryBro string) error
}

type consoleNotifier struct {
	formatter teamtable.Formatter
}

func NewConsoleNotifier(formatter teamtable.Formatter) Notifier {
	return &consoleNotifier{formatter: formatter}
}

func (c *consoleNotifier) Notify(team1, team2 teambuilder.Team, sorryBro string) error {
	teamTable := teamtable.NewTeamTable(team1, team2, sorryBro)
	message := c.formatter.Format(teamTable)
	fmt.Println(message)
	return nil
}
