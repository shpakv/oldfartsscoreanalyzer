package telegram

import (
	"fmt"
	"oldfartscounter/internal/teambuilder"
	"oldfartscounter/internal/teamtable"
)

type Notifier struct {
	handler   API
	formatter teamtable.Formatter
}

func NewNotifier(handler API, formatter teamtable.Formatter) *Notifier {
	return &Notifier{
		handler:   handler,
		formatter: formatter,
	}
}

func (n *Notifier) Notify(team1, team2 teambuilder.Team) error {
	teamTable := teamtable.NewTeamTable(team1, team2)
	message := n.formatter.Format(teamTable)
	fmt.Println(message)

	return n.handler.SendMessage(message)
}
