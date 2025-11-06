package main

import (
	"fmt"
	"log"
	"oldfartscounter/internal/environment"
	"oldfartscounter/internal/notifier"
	"oldfartscounter/internal/teambuilder"
	"oldfartscounter/internal/telegram"
	"oldfartscounter/internal/tui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var SorryBro string

func main() {
	// Создаем репозиторий игроков
	repo := teambuilder.NewPlayerRepository()

	// Создаем notifier для Telegram
	telegramFormatter := telegram.NewTeamTableFormatter()
	notifiers := []notifier.Notifier{
		telegram.NewNotifier(apiHandler(), telegramFormatter),
	}

	// Создаем модель TUI
	model := tui.NewModel(repo, notifiers)
	if SorryBro != "" {
		model.SetSorryBro(SorryBro)
	}

	// Запускаем Bubble Tea приложение
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Ошибка запуска приложения: %v\n", err)
		os.Exit(1)
	}

	log.Println("Спасибо за использование Old Farts Team Builder!")
}

func apiHandler() *telegram.DefaultAPIHandler {
	bot := telegram.NewBotFromEnv()
	chatId := environment.GetVariable("TELEGRAM_CHAT_ID", telegram.ChatID)
	return telegram.NewDefaultAPIHandler(bot, chatId)
}
