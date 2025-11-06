package tui

import (
	"fmt"
	"oldfartscounter/internal/teambuilder"
	"oldfartscounter/internal/tui/styles"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func updateConstraints(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if m.editingConstraint {
			return updateConstraintEditor(m, keyMsg)
		}

		switch keyMsg.String() {
		case "up", "k", "л":
			if m.cursor > 0 {
				m.cursor--
			}
		case keyDown, "j", "о":
			if m.cursor < len(m.constraints) {
				m.cursor++
			}
		case "n", "т": // т - это n на русской раскладке
			// Добавить новое ограничение
			m.editingConstraint = true
			m.editingConstraintNew = true
			m.constraintPlayer1Idx = 0
			// Инициализируем второго игрока как 1, чтобы избежать дублирования
			selectedPlayers := m.getSelectedPlayersList()
			if len(selectedPlayers) > 1 {
				m.constraintPlayer2Idx = 1
			} else {
				m.constraintPlayer2Idx = 0
			}
			m.constraintType = teambuilder.ConstraintTogether
			m.constraintFieldFocus = 0
		case "delete", "x", "ч": // ч - это x на русской раскладке
			// Удалить выбранное ограничение
			if m.cursor < len(m.constraints) {
				m.constraints = append(m.constraints[:m.cursor], m.constraints[m.cursor+1:]...)
				if m.cursor >= len(m.constraints) && m.cursor > 0 {
					m.cursor--
				}
			}
		case keyTab, keyEnter:
			// Генерируем команды и переходим к результатам
			m.generateTeams()
			m.currentScreen = ScreenResults
			m.cursor = 0
		case keyEsc:
			// Возврат к выбору игроков
			m.currentScreen = ScreenPlayers
			m.cursor = 0
		}
	}
	return m, nil
}

func updateConstraintEditor(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	selectedPlayers := m.getSelectedPlayersList()

	// Проверка границ индексов
	if m.constraintPlayer1Idx >= len(selectedPlayers) {
		m.constraintPlayer1Idx = 0
	}
	if m.constraintPlayer2Idx >= len(selectedPlayers) {
		m.constraintPlayer2Idx = len(selectedPlayers) - 1
	}
	if m.constraintPlayer1Idx < 0 {
		m.constraintPlayer1Idx = 0
	}
	if m.constraintPlayer2Idx < 0 {
		m.constraintPlayer2Idx = 0
	}

	switch msg.String() {
	case "up", "k", "л":
		// Навигация вверх в зависимости от активного поля
		switch m.constraintFieldFocus {
		case 0: // Игрок 1
			if m.constraintPlayer1Idx > 0 {
				m.constraintPlayer1Idx--
			}
		case 1: // Игрок 2
			if m.constraintPlayer2Idx > 0 {
				m.constraintPlayer2Idx--
			}
		case 2: // Тип - переключаем
			if m.constraintType == teambuilder.ConstraintSeparate {
				m.constraintType = teambuilder.ConstraintTogether
			}
		}
	case keyDown, "j", "о":
		// Навигация вниз в зависимости от активного поля
		switch m.constraintFieldFocus {
		case 0: // Игрок 1
			if m.constraintPlayer1Idx < len(selectedPlayers)-1 {
				m.constraintPlayer1Idx++
			}
		case 1: // Игрок 2
			if m.constraintPlayer2Idx < len(selectedPlayers)-1 {
				m.constraintPlayer2Idx++
			}
		case 2: // Тип - переключаем
			if m.constraintType == teambuilder.ConstraintTogether {
				m.constraintType = teambuilder.ConstraintSeparate
			}
		}
	case keyTab:
		// Переключение между полями вперед
		m.constraintFieldFocus = (m.constraintFieldFocus + 1) % 3
	case "left", "h", "р":
		// Переключение между полями назад
		m.constraintFieldFocus = (m.constraintFieldFocus + 2) % 3
	case "right", "l", "д":
		// Переключение между полями вперед
		m.constraintFieldFocus = (m.constraintFieldFocus + 1) % 3
	case "space":
		// Переключение типа constraint
		if m.constraintType == teambuilder.ConstraintTogether {
			m.constraintType = teambuilder.ConstraintSeparate
		} else {
			m.constraintType = teambuilder.ConstraintTogether
		}
	case keyEnter:
		// Сохранить constraint
		if len(selectedPlayers) >= 2 {
			player1 := m.allPlayers[selectedPlayers[m.constraintPlayer1Idx]].NickName
			player2 := m.allPlayers[selectedPlayers[m.constraintPlayer2Idx]].NickName

			if player1 != player2 {
				newConstraint := teambuilder.Constraint{
					Type:    m.constraintType,
					Player1: player1,
					Player2: player2,
				}

				if m.editingConstraintNew {
					m.constraints = append(m.constraints, newConstraint)
				}

				m.editingConstraint = false
				m.editingConstraintNew = false
				m.errorMsg = ""
			} else {
				m.errorMsg = "Нельзя выбрать одного и того же игрока"
			}
		}
	case keyEsc:
		// Отмена редактирования
		m.editingConstraint = false
		m.editingConstraintNew = false
		m.errorMsg = ""
	}

	return m, nil
}

func viewConstraints(m Model) string {
	var b strings.Builder

	// Если редактируем constraint
	if m.editingConstraint {
		return viewConstraintEditor(m)
	}

	// Заголовок
	title := styles.TitleStyle.Render("Управление ограничениями (Constraints)")
	subtitle := styles.SubtitleStyle.Render(fmt.Sprintf("Всего ограничений: %d", len(m.constraints)))

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorGrid).
		Padding(1, 2).
		Width(80).
		Render(title + "\n" + subtitle))
	b.WriteString("\n\n")

	// Список constraints
	constraintsList := ""
	if len(m.constraints) == 0 {
		constraintsList = styles.SubtitleStyle.Render("  Ограничений пока нет. Нажмите 'N' чтобы добавить.")
	} else {
		for i, constraint := range m.constraints {
			cursor := " "
			itemStyle := styles.UnselectedItemStyle
			if i == m.cursor {
				cursor = keyCursor
				itemStyle = styles.SelectedItemStyle
			}

			icon := "🤝"
			typeText := "вместе"
			if constraint.Type == teambuilder.ConstraintSeparate {
				icon = "💔"
				typeText = "раздельно"
			}

			line := fmt.Sprintf("%s %s %-20s ↔ %-20s (%s)",
				cursor,
				icon,
				itemStyle.Render(constraint.Player1),
				itemStyle.Render(constraint.Player2),
				typeText,
			)
			constraintsList += line + "\n"
		}
	}

	b.WriteString(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorGrid).
		Padding(1, 2).
		Width(80).
		Height(12).
		Render(constraintsList))
	b.WriteString("\n\n")

	// Помощь
	help := styles.HelpStyle.Render("↑/↓: навигация • N: добавить • X: удалить • Tab: далее • Esc: назад")
	b.WriteString(help)

	return b.String()
}

func viewConstraintEditor(m Model) string {
	var b strings.Builder

	selectedPlayers := m.getSelectedPlayersList()
	if len(selectedPlayers) < 2 {
		return styles.ErrorStyle.Render("Недостаточно игроков для создания ограничения")
	}

	// Заголовок
	title := styles.TitleStyle.Render("Добавление ограничения")
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorGrid).
		Padding(1, 2).
		Width(80).
		Render(title))
	b.WriteString("\n\n")

	// Выбор игрока 1
	player1Label := "[1] Игрок 1:"
	player1Active := ""
	if m.constraintFieldFocus == 0 {
		player1Label = styles.AccentTextStyle.Render("► [1] Игрок 1:")
	}

	player1Name := m.allPlayers[selectedPlayers[m.constraintPlayer1Idx]].NickName

	player1Box := fmt.Sprintf("%s %s%s", player1Label, player1Name, player1Active)
	b.WriteString(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(func() lipgloss.Color {
			if m.constraintFieldFocus == 0 {
				return styles.ColorAccent
			}
			return styles.ColorGrid
		}()).
		Padding(0, 1).
		Width(80).
		Render(player1Box))
	b.WriteString("\n")

	// Выбор игрока 2
	player2Label := "[2] Игрок 2:"
	player2Active := ""
	if m.constraintFieldFocus == 1 {
		player2Label = styles.AccentTextStyle.Render("► [2] Игрок 2:")
	}

	player2Name := m.allPlayers[selectedPlayers[m.constraintPlayer2Idx]].NickName

	player2Box := fmt.Sprintf("%s %s%s", player2Label, player2Name, player2Active)
	b.WriteString(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(func() lipgloss.Color {
			if m.constraintFieldFocus == 1 {
				return styles.ColorAccent
			}
			return styles.ColorGrid
		}()).
		Padding(0, 1).
		Width(80).
		Render(player2Box))
	b.WriteString("\n")

	// Выбор типа
	typeLabel := "[3] Тип:"
	if m.constraintFieldFocus == 2 {
		typeLabel = styles.AccentTextStyle.Render("► [3] Тип:")
	}

	typeText := ""
	if m.constraintType == teambuilder.ConstraintTogether {
		typeText = "[X] Вместе [ ] Раздельно"
	} else {
		typeText = "[ ] Вместе [X] Раздельно"
	}

	typeBox := fmt.Sprintf("%s %s", typeLabel, typeText)
	b.WriteString(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(func() lipgloss.Color {
			if m.constraintFieldFocus == 2 {
				return styles.ColorAccent
			}
			return styles.ColorGrid
		}()).
		Padding(0, 1).
		Width(80).
		Render(typeBox))
	b.WriteString("\n\n")

	// Ошибка, если есть
	if m.errorMsg != "" {
		errorBox := styles.ErrorStyle.Render("⚠ " + m.errorMsg)
		b.WriteString(errorBox)
		b.WriteString("\n\n")
	}

	// Помощь
	helpText := ""
	switch m.constraintFieldFocus {
	case 0, 1:
		helpText = "↑/↓: выбор игрока • Tab/←/→: переключить поле • Enter: сохранить • Esc: отмена"
	case 2:
		helpText = "↑/↓ или Space: изменить тип • Tab/←/→: переключить поле • Enter: сохранить • Esc: отмена"
	}
	help := styles.HelpStyle.Render(helpText)
	b.WriteString(help)

	return b.String()
}
