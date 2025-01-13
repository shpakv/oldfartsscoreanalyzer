package telegram

import "oldfartscounter/internal/environment"

type Bot struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

func NewBotFromEnv() *Bot {
	return &Bot{
		Id:    environment.GetVariable("TELEGRAM_BOT_ID"),
		Name:  environment.GetVariable("TELEGRAM_BOT_NAME"),
		Token: environment.GetVariable("TELEGRAM_BOT_TOKEN"),
	}
}
