package main

import (
	"encoding/json"
	"flag"
	"log"
	"oldfartscounter/internal/environment"
	"oldfartscounter/internal/notifier"
	"oldfartscounter/internal/teambuilder"
	"oldfartscounter/internal/telegram"
	"os"
)

func main() {
	c := config()
	f := telegram.NewTeamTableFormatter()
	notifiers := []notifier.Notifier{
		notifier.NewConsoleNotifier(f),
		telegram.NewNotifier(apiHandler(), f),
	}
	teamBuilder := &teambuilder.TeamBuilder{}
	team1, team2 := teamBuilder.Build(c)
	for _, n := range notifiers {
		err := n.Notify(team1, team2)
		if err != nil {
			log.Fatalf("Failed to notify old farts: %v", err)
		}
	}
}

func config() *teambuilder.TeamConfiguration {
	filePath := flag.String("c", "config.json5", "Path to the config.json file")
	flag.Parse()
	if *filePath == "" {
		log.Fatal("Please provide a file path using the -c flag")
	}
	if _, err := os.Stat(*filePath); os.IsNotExist(err) {
		log.Fatalf("File does not exist: %s", *filePath)
	}
	content, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var c teambuilder.TeamConfiguration
	if err = json.Unmarshal(content, &c); err != nil {
		log.Fatalf("Invalid JSON format: %v", err)
	}
	return &c
}

func apiHandler() *telegram.DefaultAPIHandler {
	bot := telegram.NewBotFromEnv()
	chatId := environment.GetVariable("TELEGRAM_CHAT_ID")
	return telegram.NewDefaultAPIHandler(bot, chatId)
}
