package telegram

import "oldfartscounter/internal/environment"

const ChatID = "-1002150403113"
const BotID = "-veryveryoldfartbot"
const BotName = "OldFartsBot"
const BotToken = "7890733898:AAF1Vak1FwUQm-aU1OCrm6svxSXVS9gnago"

type Bot struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

func NewBotFromEnv() *Bot {
	return &Bot{
		Id:    environment.GetVariable("TELEGRAM_BOT_ID", BotID),
		Name:  environment.GetVariable("TELEGRAM_BOT_NAME", BotName),
		Token: environment.GetVariable("TELEGRAM_BOT_TOKEN", BotToken),
	}
}
