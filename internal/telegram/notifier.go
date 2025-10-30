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

func (n *Notifier) Notify(teams []teambuilder.Team, sorryBro string) error {
	teamTable := teamtable.NewTeamTableMultiple(teams, sorryBro)
	message := n.formatter.Format(teamTable)

	return n.handler.SendMessage(message)
}
