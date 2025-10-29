package main

import (
	"flag"
	"log"
	"oldfartscounter/internal/environment"
	"oldfartscounter/internal/notifier"
	"oldfartscounter/internal/teambuilder"
	"oldfartscounter/internal/telegram"
	"os"

	"github.com/yosuke-furukawa/json5/encoding/json5"
)

var SorryBro = ""

func main() {
	c := config()
	f := telegram.NewTeamTableFormatter()
	notifiers := []notifier.Notifier{
		notifier.NewConsoleNotifier(f),
		//telegram.NewNotifier(apiHandler(), f),
	}
	repo := teambuilder.NewPlayerRepository()
	teamBuilder := teambuilder.NewTeamBuilder(repo)

	// Проверяем количество команд
	numTeams := c.NumTeams
	if numTeams != 2 && numTeams != 4 {
		numTeams = 2 // Default to 2 teams
	}

	if numTeams == 4 {
		// Используем новый метод для 4 команд
		teams := teamBuilder.BuildMultiple(c)
		for _, n := range notifiers {
			err := n.NotifyMultiple(teams, SorryBro)
			if err != nil {
				log.Fatalf("Failed to notify old farts: %v", err)
			}
		}
	} else {
		// Используем старый метод для 2 команд
		team1, team2 := teamBuilder.Build(c)
		for _, n := range notifiers {
			err := n.Notify(team1, team2, SorryBro)
			if err != nil {
				log.Fatalf("Failed to notify old farts: %v", err)
			}
		}
	}
}

func config() *teambuilder.TeamConfiguration {
	filePath := flag.String("c", "bin/config.json5", "Path to the config.json file")
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
	if err = json5.Unmarshal(content, &c); err != nil {
		log.Fatalf("Invalid JSON format: %v", err)
	}
	if c.SorryBro == nil {
		c.SorryBro = &SorryBro
	}
	return &c
}

func apiHandler() *telegram.DefaultAPIHandler {
	bot := telegram.NewBotFromEnv()
	chatId := environment.GetVariable("TELEGRAM_CHAT_ID", telegram.ChatID)
	return telegram.NewDefaultAPIHandler(bot, chatId)
}
