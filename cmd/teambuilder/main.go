package main

import (
	"flag"
	"log"
	"oldfartscounter/internal/notifier"
	"oldfartscounter/internal/teambuilder"
	"oldfartscounter/internal/telegram"
	"os"

	"github.com/yosuke-furukawa/json5/encoding/json5"
)

var SorryBro = "Mr. Titspervert"

func main() {
	c := config()
	f := telegram.NewTeamTableFormatter()
	notifiers := []notifier.Notifier{
		notifier.NewConsoleNotifier(f),
		// telegram.NewNotifier(apiHandler(), f),
	}
	repo := teambuilder.NewPlayerRepository()
	teamBuilder := teambuilder.NewTeamBuilder(repo)

	teams := teamBuilder.Build(c)

	for _, n := range notifiers {
		err := n.Notify(teams, SorryBro)
		if err != nil {
			log.Fatalf("Failed to notify old farts: %v", err)
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

// apiHandler создает API handler для Telegram (закомментировано, но оставлено для будущего использования)
// func apiHandler() *telegram.DefaultAPIHandler {
// 	bot := telegram.NewBotFromEnv()
// 	chatId := environment.GetVariable("TELEGRAM_CHAT_ID", telegram.ChatID)
// 	return telegram.NewDefaultAPIHandler(bot, chatId)
// }
