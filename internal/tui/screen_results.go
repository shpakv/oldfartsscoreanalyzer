package tui

import (
	"fmt"
	"oldfartscounter/internal/teambuilder"
	"oldfartscounter/internal/tui/styles"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func updateResults(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "s", "ы": // ы - это s на русской раскладке
			// Выбор SorryBro
			selectedPlayers := m.getSelectedPlayersList()
			if len(selectedPlayers) > 0 {
				// Переключаемся в режим выбора SorryBro
				m.cursor = 0
				m.searchMode = true // Используем searchMode как флаг для режима выбора SorryBro
			}
		case "up", "k", "л": // навигация в режиме выбора SorryBro
			if m.searchMode { // searchMode используется для режима SorryBro
				if m.cursor > 0 {
					m.cursor--
				}
			}
		case "down", "j", "о":
			if m.searchMode {
				selectedPlayers := m.getSelectedPlayersList()
				if m.cursor < len(selectedPlayers) {
					m.cursor++
				}
			}
		case keyEnter:
			if m.searchMode {
				// Подтверждаем выбор SorryBro
				selectedPlayers := m.getSelectedPlayersList()
				if m.cursor < len(selectedPlayers) {
					playerName := m.allPlayers[selectedPlayers[m.cursor]].NickName
					m.sorryBro = &playerName
					m.generateTeams() // Перегенерируем с новым SorryBro
				} else if m.cursor == len(selectedPlayers) {
					// Выбрали "Нет SorryBro"
					m.sorryBro = nil
					m.generateTeams()
				}
				m.searchMode = false
				m.cursor = 0
			}
		case "t", "е": // е - это t на русской раскладке
			if !m.searchMode {
				// Отправить в Telegram (немедленно)
				sorryBroName := ""
				if m.sorryBro != nil {
					sorryBroName = *m.sorryBro
				}

				// Очищаем предыдущие сообщения
				m.errorMsg = ""
				m.successMsg = ""

				// Вызываем telegram notifier
				if len(m.notifiers) > 0 {
					if err := m.notifiers[0].Notify(m.generatedTeams, sorryBroName); err != nil {
						m.errorMsg = fmt.Sprintf("Ошибка отправки в Telegram: %v", err)
						return m, nil
					}
				}

				m.successMsg = "Результаты отправлены в чат OldFarts"
			}
		case "r", "к": // к - это r на русской раскладке
			if !m.searchMode {
				// Очищаем сообщения
				m.errorMsg = ""
				m.successMsg = ""
				// Перегенерировать команды
				m.generateTeams()
			}
		case keyEsc:
			if m.searchMode {
				// Выход из режима выбора SorryBro
				m.searchMode = false
				m.cursor = 0
			} else {
				// Возврат в меню
				m.currentScreen = ScreenMenu
				m.cursor = 0
			}
		case "q", "й", keyCtrlC: // й - это q на русской раскладке
			return m, tea.Quit
		}
	}
	return m, nil
}

func viewResults(m Model) string {
	var b strings.Builder

	// Если в режиме выбора SorryBro - показываем UI выбора
	if m.searchMode {
		return viewSorryBroSelector(m)
	}

	// Заголовок
	title := styles.TitleStyle.Render("Сгенерированные команды")
	subtitle := ""
	if m.sorryBro != nil {
		subtitle = styles.SubtitleStyle.Render(fmt.Sprintf("SorryBro: %s", *m.sorryBro))
	} else {
		subtitle = styles.SubtitleStyle.Render("SorryBro: не выбран")
	}

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorGrid).
		Padding(1, 2).
		Width(130).
		Render(title + "\n" + subtitle))
	b.WriteString("\n\n")

	// Визуализация команд (всегда 2)
	b.WriteString(renderTwoTeams(m))

	b.WriteString("\n\n")

	// Статистика баланса
	b.WriteString(renderBalance(m))
	b.WriteString("\n\n")

	// Сообщения об ошибках и успехе
	if m.errorMsg != "" {
		errorBox := styles.ErrorStyle.Render("⚠ " + m.errorMsg)
		b.WriteString(errorBox)
		b.WriteString("\n\n")
	}
	if m.successMsg != "" {
		successBox := styles.SuccessStyle.Render("✅ " + m.successMsg)
		b.WriteString(successBox)
		b.WriteString("\n\n")
	}

	// Помощь
	help := styles.HelpStyle.Render("S: выбрать SorryBro • T: отправить в чат OldFarts • R: перегенерировать • Esc: в меню • Q: выход")
	b.WriteString(help)

	return b.String()
}

func renderTwoTeams(m Model) string {
	if len(m.generatedTeams) < 2 {
		return styles.ErrorStyle.Render("Ошибка генерации команд")
	}

	team1 := m.generatedTeams[0]
	team2 := m.generatedTeams[1]

	// Определяем максимальный рейтинг для градиента и средний для категорий
	maxRating := 0.0
	for _, team := range m.generatedTeams {
		for _, player := range team {
			if player.Score > maxRating {
				maxRating = player.Score
			}
		}
	}
	averageMu := m.calculateAverageMu()

	// Рендерим команду 1
	team1Box := renderTeam(team1, "Команда 1", maxRating, averageMu)

	// Рендерим команду 2
	team2Box := renderTeam(team2, "Команда 2", maxRating, averageMu)

	// Размещаем команды рядом
	return lipgloss.JoinHorizontal(lipgloss.Top, team1Box, "  ", team2Box)
}

func renderTeam(team teambuilder.Team, teamName string, maxRating float64, averageMu float64) string {
	var b strings.Builder

	// Заголовок команды
	b.WriteString(styles.AccentTextStyle.Render(teamName))
	b.WriteString("\n\n")

	// Игроки
	for _, player := range team {
		progressBar := styles.RenderProgressBar(player.Score, maxRating, 10)

		// Категория рейтинга (используем средний рейтинг)
		category, categoryColor := styles.GetRatingCategory(player.Score, averageMu)
		categoryStyle := lipgloss.NewStyle().Foreground(categoryColor).Bold(true)

		line := fmt.Sprintf("%-25s %s %s",
			player.NickName,
			progressBar,
			categoryStyle.Render(category),
		)
		b.WriteString(line)
		b.WriteString("\n")
	}

	// Суммарный рейтинг
	totalScore := team.Score()
	b.WriteString("\n")
	b.WriteString(styles.StatStyle.Render(fmt.Sprintf("Σ: %.0f", totalScore)))

	// Оборачиваем в панель
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorAccent).
		Padding(1, 2).
		Width(60).
		Render(b.String())
}

func renderBalance(m Model) string {
	if len(m.generatedTeams) == 0 {
		return ""
	}

	// Находим минимальный и максимальный рейтинг команд
	minScore := m.generatedTeams[0].Score()
	maxScore := m.generatedTeams[0].Score()

	for _, team := range m.generatedTeams {
		score := team.Score()
		if score < minScore {
			minScore = score
		}
		if score > maxScore {
			maxScore = score
		}
	}

	diff := maxScore - minScore
	avgScore := 0.0
	for _, team := range m.generatedTeams {
		avgScore += team.Score()
	}
	avgScore /= float64(len(m.generatedTeams))

	diffPercent := 0.0
	if avgScore > 0 {
		diffPercent = (diff / avgScore) * 100
	}

	// Оценка баланса
	var balanceText string
	var balanceStyle lipgloss.Style

	switch {
	case diffPercent < 2.0:
		balanceText = "✓ Отлично!"
		balanceStyle = styles.SuccessStyle
	case diffPercent < 5.0:
		balanceText = "✓ Хорошо"
		balanceStyle = styles.SuccessStyle
	case diffPercent < 10.0:
		balanceText = "⚠ Удовлетворительно"
		balanceStyle = styles.WarningStyle
	default:
		balanceText = "⚠ Плохой баланс"
		balanceStyle = styles.ErrorStyle
	}

	result := fmt.Sprintf("Разница балансов: %.0f (%.2f%%)  %s",
		diff,
		diffPercent,
		balanceStyle.Render(balanceText),
	)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorGrid).
		Padding(1, 2).
		Width(100).
		Render(result)
}

func viewSorryBroSelector(m Model) string {
	var b strings.Builder

	// Заголовок
	title := styles.TitleStyle.Render("Выбор SorryBro")
	subtitle := styles.SubtitleStyle.Render("Игрок, который останется за бортом")

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorAccent).
		Padding(1, 2).
		Width(80).
		Render(title + "\n" + subtitle))
	b.WriteString("\n\n")

	// Список игроков для выбора
	selectedPlayers := m.getSelectedPlayersList()
	playersList := ""

	// Опция "Нет SorryBro"
	cursor := " "
	itemStyle := styles.UnselectedItemStyle
	if m.cursor == len(selectedPlayers) {
		cursor = keyCursor
		itemStyle = styles.SelectedItemStyle
	}
	playersList += fmt.Sprintf("%s %s\n", cursor, itemStyle.Render("Нет SorryBro (все играют)"))
	playersList += "\n"

	// Список игроков
	for i, playerIdx := range selectedPlayers {
		player := m.allPlayers[playerIdx]

		cursor = " "
		itemStyle = styles.UnselectedItemStyle
		if i == m.cursor {
			cursor = keyCursor
			itemStyle = styles.SelectedItemStyle
		}

		// Проверяем, не выбран ли уже этот игрок как SorryBro
		isCurrentSorryBro := ""
		if m.sorryBro != nil && *m.sorryBro == player.NickName {
			isCurrentSorryBro = " ✓ текущий"
		}

		line := fmt.Sprintf("%s %-30s%s",
			cursor,
			itemStyle.Render(player.NickName),
			isCurrentSorryBro,
		)
		playersList += line + "\n"
	}

	b.WriteString(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorGrid).
		Padding(1, 2).
		Width(80).
		Height(20).
		Render(playersList))
	b.WriteString("\n\n")

	// Помощь
	help := styles.HelpStyle.Render("↑/↓,K/J: навигация • Enter: выбрать • Esc: отмена")
	b.WriteString(help)

	return b.String()
}
