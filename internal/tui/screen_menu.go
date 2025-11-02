package tui

import (
	"fmt"
	"oldfartscounter/internal/teambuilder"
	"oldfartscounter/internal/tui/styles"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func updateMenu(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "up", "k":
			if m.menuCursor > 0 {
				m.menuCursor--
			}
		case "down", "j":
			if m.menuCursor < 1 {
				m.menuCursor++
			}
		case keyEnter:
			switch m.menuCursor {
			case 0: // Создать команды
				m.currentScreen = ScreenPlayers
				m.selectedPlayers = make(map[int]bool)
				m.constraints = []teambuilder.Constraint{}
				m.numTeams = 2
				m.sorryBro = nil
				m.errorMsg = ""
			case 1: // Выход
				return m, tea.Quit
			}
		case "q", "й", keyCtrlC: // й - это q на русской раскладке
			return m, tea.Quit
		case keyEsc:
			m.errorMsg = ""
		}
	}
	return m, nil
}

func viewMenu(m Model) string {
	var b strings.Builder

	// Заголовок
	title := styles.TitleStyle.Render("🎮 Old Farts Team Builder")
	subtitle := styles.SubtitleStyle.Render("Балансировщик команд для CS2")

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorGrid).
		Padding(1, 2).
		Width(60).
		Render(title + "\n" + subtitle))
	b.WriteString("\n\n")

	// Меню
	menuItems := []string{
		"Создать команды",
		"Выход",
	}

	menuBox := ""
	for i, item := range menuItems {
		cursor := " "
		if m.menuCursor == i {
			cursor = keyCursor
			item = styles.SelectedItemStyle.Render(item)
		} else {
			item = styles.UnselectedItemStyle.Render(item)
		}
		menuBox += fmt.Sprintf("%s %s\n", cursor, item)
	}

	b.WriteString(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorGrid).
		Padding(1, 2).
		Width(60).
		Render(menuBox))
	b.WriteString("\n\n")

	// Ошибка, если есть
	if m.errorMsg != "" {
		errorBox := styles.ErrorStyle.Render("⚠ " + m.errorMsg)
		b.WriteString(lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorError).
			Padding(1, 2).
			Width(60).
			Render(errorBox))
		b.WriteString("\n\n")
	}

	// Помощь
	help := styles.HelpStyle.Render("↑/↓: навигация • Enter: выбрать • Q: выход")
	b.WriteString(help)

	return b.String()
}
