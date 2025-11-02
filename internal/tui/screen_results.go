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
		case "t", "е": // е - это t на русской раскладке
			// Отправить в Telegram (немедленно)
			sorryBroName := ""
			if m.sorryBro != nil {
				sorryBroName = *m.sorryBro
			}

			// Вызываем telegram notifier
			if len(m.notifiers) > 0 {
				if err := m.notifiers[0].Notify(m.generatedTeams, sorryBroName); err != nil {
					m.errorMsg = fmt.Sprintf("Ошибка отправки в Telegram: %v", err)
					return m, nil
				}
			}

			m.errorMsg = "✅ Результаты отправлены в чат OldFarts"
		case "r", "к": // к - это r на русской раскладке
			// Перегенерировать команды
			m.generateTeams()
		case keyEsc:
			// Возврат в меню
			m.currentScreen = ScreenMenu
			m.cursor = 0
		case "q", "й", keyCtrlC: // й - это q на русской раскладке
			return m, tea.Quit
		}
	}
	return m, nil
}

func viewResults(m Model) string {
	var b strings.Builder

	// Заголовок
	title := styles.TitleStyle.Render("Сгенерированные команды")
	numTeamsText := fmt.Sprintf("%d команды", m.numTeams)
	if m.numTeams == 4 {
		numTeamsText = "4 команды"
	}
	subtitle := styles.SubtitleStyle.Render(numTeamsText)

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorGrid).
		Padding(1, 2).
		Width(100).
		Render(title + "\n" + subtitle))
	b.WriteString("\n\n")

	// Визуализация команд
	if m.numTeams == 2 {
		b.WriteString(renderTwoTeams(m))
	} else {
		b.WriteString(renderFourTeams(m))
	}

	b.WriteString("\n\n")

	// Статистика баланса
	b.WriteString(renderBalance(m))
	b.WriteString("\n\n")

	// Ошибка, если есть
	if m.errorMsg != "" {
		errorBox := styles.ErrorStyle.Render("⚠ " + m.errorMsg)
		b.WriteString(errorBox)
		b.WriteString("\n\n")
	}

	// Помощь
	help := styles.HelpStyle.Render("T: отправить в чат OldFarts • R: перегенерировать • Esc: в меню • Q: выход")
	b.WriteString(help)

	return b.String()
}

func renderTwoTeams(m Model) string {
	if len(m.generatedTeams) < 2 {
		return styles.ErrorStyle.Render("Ошибка генерации команд")
	}

	team1 := m.generatedTeams[0]
	team2 := m.generatedTeams[1]

	// Определяем максимальный рейтинг для градиента
	maxRating := 0.0
	for _, team := range m.generatedTeams {
		for _, player := range team {
			if player.Score > maxRating {
				maxRating = player.Score
			}
		}
	}

	// Рендерим команду 1
	team1Box := renderTeam(team1, "Команда 1", maxRating)

	// Рендерим команду 2
	team2Box := renderTeam(team2, "Команда 2", maxRating)

	// Размещаем команды рядом
	return lipgloss.JoinHorizontal(lipgloss.Top, team1Box, "  ", team2Box)
}

func renderFourTeams(m Model) string {
	if len(m.generatedTeams) < 4 {
		return styles.ErrorStyle.Render("Ошибка генерации команд")
	}

	// Определяем максимальный рейтинг для градиента
	maxRating := 0.0
	for _, team := range m.generatedTeams {
		for _, player := range team {
			if player.Score > maxRating {
				maxRating = player.Score
			}
		}
	}

	// Рендерим команды
	team1Box := renderTeam(m.generatedTeams[0], "Команда 1", maxRating)
	team2Box := renderTeam(m.generatedTeams[1], "Команда 2", maxRating)
	team3Box := renderTeam(m.generatedTeams[2], "Команда 3", maxRating)
	team4Box := renderTeam(m.generatedTeams[3], "Команда 4", maxRating)

	// Размещаем команды в 2 ряда по 2
	row1 := lipgloss.JoinHorizontal(lipgloss.Top, team1Box, "  ", team2Box)
	row2 := lipgloss.JoinHorizontal(lipgloss.Top, team3Box, "  ", team4Box)

	return lipgloss.JoinVertical(lipgloss.Left, row1, "\n", row2)
}

func renderTeam(team teambuilder.Team, teamName string, maxRating float64) string {
	var b strings.Builder

	// Заголовок команды
	b.WriteString(styles.AccentTextStyle.Render(teamName))
	b.WriteString("\n\n")

	// Игроки
	for _, player := range team {
		progressBar := styles.RenderProgressBar(player.Score, maxRating, 10)
		line := fmt.Sprintf("%-25s %s %6.0f",
			player.NickName,
			progressBar,
			player.Score,
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
		Width(45).
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
