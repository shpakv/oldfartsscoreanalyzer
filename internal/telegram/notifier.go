package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"oldfartscounter/internal/teambuilder"
	"oldfartscounter/internal/teamtable"
)

type Notifier struct {
	Bot       *Bot
	ChatId    string
	Formatter *Formatter
}

func NewNotifier(bot *Bot, chatId string) *Notifier {
	return &Notifier{
		Bot:       bot,
		ChatId:    chatId,
		Formatter: NewFormatter(),
	}
}

func (n *Notifier) NotifyOldFarts(team1, team2 teambuilder.Team) error {
	teamTable := teamtable.NewTeamTable(team1, team2)
	message := n.Formatter.Format(teamTable)
	fmt.Println(message)

	return makeTelegramRequest(n.Bot.Token, n.ChatId, message)
}

func makeTelegramRequest(botToken, chatId, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	body := map[string]string{
		"chat_id":    chatId,
		"text":       message,
		"parse_mode": "Markdown",
	}
	bodyJSON, _ := json.Marshal(body)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyJSON))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed with status code: %d", resp.StatusCode)
	}

	return nil
}
