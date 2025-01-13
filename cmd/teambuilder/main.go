package main

import (
	"encoding/json"
	"flag"
	"log"
	"oldfartscounter/internal/environment"
	"oldfartscounter/internal/teambuilder"
	"oldfartscounter/internal/telegram"
	"os"
)

func main() {
	filePath := flag.String("f", "", "Path to the configuration.json file")
	flag.Parse()
	if *filePath == "" {
		log.Fatal("Please provide a file path using the -f flag")
	}
	if _, err := os.Stat(*filePath); os.IsNotExist(err) {
		log.Fatalf("File does not exist: %s", *filePath)
	}
	content, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var config teambuilder.TeamConfiguration
	if err = json.Unmarshal(content, &config); err != nil {
		log.Fatalf("Invalid JSON format: %v", err)
	}

	bot := telegram.NewBotFromEnv()
	chatId := environment.GetVariable("TELEGRAM_CHAT_ID")
	teamBuilder := &teambuilder.TeamBuilder{}
	teamTableFormatter := telegram.NewTeamTableFormatter()
	apiHandler := telegram.NewDefaultAPIHandler(bot, chatId)
	notifier := telegram.NewNotifier(apiHandler, teamTableFormatter)

	team1, team2 := teamBuilder.Build(&config)
	err = notifier.Notify(team1, team2)
	if err != nil {
		log.Fatalf("Failed to notify old farts: %v", err)
	}
}
