package tui

import (
	"fmt"
	"oldfartscounter/internal/teambuilder"
	"oldfartscounter/internal/tui/styles"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Цель конкурса "Сорян, Братан" на текущий месяц
const sorryBroTarget = "maslina420"

func updateSettings(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	selectedPlayers := m.getSelectedPlayersList()

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "up", "k", "л":
			// Убрана навигация - только одна настройка (количество команд)
		case keyDown, "j", "о":
			// Убрана навигация - только одна настройка (количество команд)
		case "left", "h", "р":
			// Переключение количества команд
			if m.numTeams == 4 {
				m.numTeams = 2
			}
		case "right", "l", "д":
			// Переключение количества команд
			if m.numTeams == 2 {
				m.numTeams = 4
			}
		case "space", keyEnter:
			// Переключение количества команд
			if m.numTeams == 2 {
				m.numTeams = 4
			} else {
				m.numTeams = 2
			}
		case keyTab:
			// Проверяем есть ли целевой игрок в списке выбранных
			sorryBroName := sorryBroTarget
			targetFound := false
			for _, playerIdx := range selectedPlayers {
				if m.allPlayers[playerIdx].NickName == sorryBroTarget {
					targetFound = true
					break
				}
			}

			// Устанавливаем SorryBro только если цель найдена
			if targetFound {
				m.sorryBro = &sorryBroName
			} else {
				m.sorryBro = nil
			}

			// Генерируем команды и переходим к результатам
			m.generateTeams()
			m.currentScreen = ScreenResults
			m.cursor = 0
		case keyEsc:
			// Возврат к constraints
			m.currentScreen = ScreenConstraints
			m.cursor = 0
		}
	}
	return m, nil
}

func viewSettings(m Model) string {
	var b strings.Builder

	selectedPlayers := m.getSelectedPlayersList()

	// Заголовок
	title := styles.TitleStyle.Render("Настройки генерации команд")
	subtitle := styles.SubtitleStyle.Render("Настройте параметры перед генерацией")

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorGrid).
		Padding(1, 2).
		Width(80).
		Render(title + "\n" + subtitle))
	b.WriteString("\n\n")

	// Настройка количества команд
	numTeamsLabel := "Количество команд:"
	numTeamsValue := ""
	if m.numTeams == 2 {
		numTeamsValue = "● 2 команды  ○ 4 команды"
	} else {
		numTeamsValue = "○ 2 команды  ● 4 команды"
	}

	numTeamsBox := fmt.Sprintf("► %s %s",
		styles.SelectedItemStyle.Render(numTeamsLabel),
		numTeamsValue,
	)

	b.WriteString(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorAccent).
		Padding(1, 2).
		Width(80).
		Render(numTeamsBox))
	b.WriteString("\n\n")

	// Информация о SorryBro
	sorryBroLabel := "🎯 Цель конкурса \"Сорян, Братан\" (ноябрь):"

	// Проверяем есть ли целевой игрок в списке
	targetFound := false
	for _, playerIdx := range selectedPlayers {
		if m.allPlayers[playerIdx].NickName == sorryBroTarget {
			targetFound = true
			break
		}
	}

	var sorryBroBox string
	if targetFound {
		sorryBroBox = fmt.Sprintf("%s\n\n    ✅ %s",
			sorryBroLabel,
			styles.AccentTextStyle.Render(sorryBroTarget),
		)
	} else {
		sorryBroBox = fmt.Sprintf("%s\n\n    ❌ %s (не выбран)",
			sorryBroLabel,
			styles.SubtitleStyle.Render(sorryBroTarget),
		)
	}

	b.WriteString(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorGrid).
		Padding(1, 2).
		Width(80).
		Render(sorryBroBox))
	b.WriteString("\n\n")

	// Помощь
	helpText := "←/→ или Space/Enter: переключить количество команд • Tab: сгенерировать • Esc: назад"
	help := styles.HelpStyle.Render(helpText)
	b.WriteString(help)

	return b.String()
}

// generateTeams генерирует команды на основе текущей конфигурации
func (m *Model) generateTeams() {
	selectedPlayers := m.getSelectedPlayersList()

	// Формируем конфигурацию
	players := make([]teambuilder.TeamPlayer, 0, len(selectedPlayers))
	for _, idx := range selectedPlayers {
		players = append(players, m.allPlayers[idx])
	}

	m.config = &teambuilder.TeamConfiguration{
		Players:     players,
		Constraints: m.constraints,
		NumTeams:    m.numTeams,
		SorryBro:    m.sorryBro,
	}

	// Генерируем команды
	m.generatedTeams = m.teamBuilder.Build(m.config)
}
