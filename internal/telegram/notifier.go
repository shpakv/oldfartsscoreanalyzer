package telegram

import (
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

func (n *Notifier) Notify(team1, team2 teambuilder.Team, sorryBro string) error {
	teamTable := teamtable.NewTeamTable(team1, team2, sorryBro)
	message := n.formatter.Format(teamTable)

	return n.handler.SendMessage(message)
}

func (n *Notifier) NotifyMultiple(teams []teambuilder.Team, sorryBro string) error {
	teamTable := teamtable.NewTeamTableMultiple(teams, sorryBro)
	message := n.formatter.Format(teamTable)

	return n.handler.SendMessage(message)
}
