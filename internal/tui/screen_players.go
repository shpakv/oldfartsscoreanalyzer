package tui

import (
	"fmt"
	"oldfartscounter/internal/tui/styles"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func updatePlayers(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		// Режим поиска - обрабатываем ввод текста
		if m.searchMode {
			switch keyMsg.String() {
			case keyEsc:
				// Выход из режима поиска
				m.searchMode = false
			case keyEnter:
				// Выход из режима поиска
				m.searchMode = false
			case "backspace":
				// Удаление символа из поиска
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
					m.cursor = 0
				}
			case keyCtrlC:
				return m, tea.Quit
			default:
				// Добавление символа в поиск (только печатаемые символы)
				if len(keyMsg.String()) == 1 && keyMsg.String()[0] >= 32 && keyMsg.String()[0] <= 126 {
					m.searchQuery += keyMsg.String()
					m.cursor = 0
				}
			}
			return m, nil
		}

		// Режим навигации
		switch keyMsg.String() {
		case "up", "k", "л": // л - это k на русской раскладке
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j", "о": // о - это j на русской раскладке
			filteredPlayers := m.getFilteredPlayers()
			if m.cursor < len(filteredPlayers)-1 {
				m.cursor++
			}
		case keyEnter:
			// Переключаем выбор игрока
			filteredPlayers := m.getFilteredPlayers()
			if m.cursor < len(filteredPlayers) {
				actualIdx := filteredPlayers[m.cursor].index
				m.selectedPlayers[actualIdx] = !m.selectedPlayers[actualIdx]
			}
		case "a", "ф": // ф - это a на русской раскладке
			// Выбрать всех
			filteredPlayers := m.getFilteredPlayers()
			for _, fp := range filteredPlayers {
				m.selectedPlayers[fp.index] = true
			}
		case "d", "в": // в - это d на русской раскладке
			// Отменить всех
			filteredPlayers := m.getFilteredPlayers()
			for _, fp := range filteredPlayers {
				m.selectedPlayers[fp.index] = false
			}
		case "/", ".", "ю": // / и . на разных раскладках, ю - для русской
			// Включить режим поиска
			m.searchMode = true
		case "f", "а": // а - это f на русской раскладке
			// Включить режим поиска
			m.searchMode = true
		case "c": // Английская C
			// Очистить поиск
			m.searchQuery = ""
			m.cursor = 0
		case "с": // Русская С
			// Очистить поиск
			m.searchQuery = ""
			m.cursor = 0
		case keyTab:
			// Переход к следующему экрану (constraints)
			if len(m.getSelectedPlayersList()) > 0 {
				m.currentScreen = ScreenConstraints
				m.cursor = 0
				m.searchQuery = ""
				m.searchMode = false
			} else {
				m.errorMsg = "Выберите хотя бы одного игрока"
			}
		case keyEsc:
			// Возврат в меню
			m.currentScreen = ScreenMenu
			m.cursor = 0
			m.searchQuery = ""
			m.searchMode = false
		}
	}
	return m, nil
}

func viewPlayers(m Model) string {
	var b strings.Builder

	// Заголовок
	title := styles.TitleStyle.Render("Выбор игроков")

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorGrid).
		Padding(1, 2).
		Width(80).
		Render(title))
	b.WriteString("\n\n")

	// Поле поиска
	searchBox := ""
	searchBorderColor := styles.ColorGrid
	if m.searchMode {
		searchBox = fmt.Sprintf("🔍 Поиск: %s█", m.searchQuery)
		searchBorderColor = styles.ColorAccent
	} else {
		if m.searchQuery != "" {
			searchBox = fmt.Sprintf("🔍 Фильтр: %s (нажмите / или F для редактирования, C для очистки)", m.searchQuery)
		} else {
			searchBox = "🔍 Поиск: (нажмите / или F для начала поиска)"
		}
	}

	b.WriteString(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(searchBorderColor).
		Padding(0, 1).
		Width(80).
		Render(searchBox))
	b.WriteString("\n\n")

	// Счетчик выбранных игроков
	selectedCount := len(m.getSelectedPlayersList())
	filteredCount := len(m.getFilteredPlayers())
	totalCount := len(m.allPlayers)

	counterText := ""
	if m.searchQuery != "" {
		counterText = fmt.Sprintf("📊 Выбрано: %d | Показано: %d/%d", selectedCount, filteredCount, totalCount)
	} else {
		counterText = fmt.Sprintf("📊 Выбрано: %d/%d", selectedCount, totalCount)
	}

	b.WriteString(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorGrid).
		Padding(0, 1).
		Width(80).
		Render(counterText))
	b.WriteString("\n\n")

	// Список игроков
	filteredPlayers := m.getFilteredPlayers()
	playersList := ""

	// Определяем максимальный рейтинг для градиента
	maxRating := 0.0
	for _, player := range m.allPlayers {
		if player.Score > maxRating {
			maxRating = player.Score
		}
	}

	// Улучшенная логика скроллинга
	const visibleLines = 15
	visibleStart := 0
	visibleEnd := 0

	if len(filteredPlayers) > visibleLines {
		// Скроллим только когда курсор приближается к нижней границе
		if m.cursor >= visibleLines-3 {
			// Курсор в нижней части - начинаем скроллить
			visibleStart = m.cursor - visibleLines + 4
			visibleEnd = m.cursor + 4

			// Не выходим за границы
			if visibleEnd > len(filteredPlayers) {
				visibleEnd = len(filteredPlayers)
				visibleStart = visibleEnd - visibleLines
			}
		} else {
			// Курсор в верхней части - показываем с начала
			visibleEnd = visibleLines
		}
	} else {
		// Если игроков меньше чем видимых линий
		visibleEnd = len(filteredPlayers)
	}

	// Индикатор скроллинга сверху
	if visibleStart > 0 {
		playersList += styles.SubtitleStyle.Render(fmt.Sprintf("  ▲ Еще %d игроков выше...\n", visibleStart))
	}

	for i := visibleStart; i < visibleEnd; i++ {
		fp := filteredPlayers[i]
		player := m.allPlayers[fp.index]

		// Чекбокс
		checkbox := "☐"
		checkStyle := styles.CheckboxUncheckedStyle
		if m.selectedPlayers[fp.index] {
			checkbox = "☑"
			checkStyle = styles.CheckboxCheckedStyle
		}

		// Курсор
		cursor := " "
		itemStyle := styles.UnselectedItemStyle
		if i == m.cursor {
			cursor = keyCursor
			itemStyle = styles.SelectedItemStyle
		}

		// Прогресс-бар рейтинга
		progressBar := styles.RenderProgressBar(player.Score, maxRating, 10)

		// Форматирование строки
		line := fmt.Sprintf("%s %s %-25s %s %6.0f",
			cursor,
			checkStyle.Render(checkbox),
			itemStyle.Render(player.NickName),
			progressBar,
			player.Score,
		)
		playersList += line + "\n"
	}

	// Индикатор скроллинга снизу
	if visibleEnd < len(filteredPlayers) {
		playersList += styles.SubtitleStyle.Render(fmt.Sprintf("  ▼ Еще %d игроков ниже...\n", len(filteredPlayers)-visibleEnd))
	}

	b.WriteString(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorGrid).
		Padding(1, 2).
		Width(80).
		Height(17).
		Render(playersList))
	b.WriteString("\n\n")

	// Ошибка, если есть
	if m.errorMsg != "" {
		errorBox := styles.ErrorStyle.Render("⚠ " + m.errorMsg)
		b.WriteString(errorBox)
		b.WriteString("\n\n")
	}

	// Помощь
	helpText := ""
	if m.searchMode {
		helpText = "Режим поиска: печатайте для поиска • Backspace: удалить • Enter/Esc: выход из поиска"
	} else {
		helpText = "↑/↓,K/J: навигация • Enter: выбрать • A: все • D: отменить • /,F: поиск • C: очистить • Tab: далее • Esc: назад"
	}
	help := styles.HelpStyle.Render(helpText)
	b.WriteString(help)

	return b.String()
}

// filteredPlayer содержит игрока и его оригинальный индекс
type filteredPlayer struct {
	index int
}

// getFilteredPlayers возвращает отфильтрованных игроков по поисковому запросу
func (m Model) getFilteredPlayers() []filteredPlayer {
	if m.searchQuery == "" {
		result := make([]filteredPlayer, len(m.allPlayers))
		for i := range m.allPlayers {
			result[i] = filteredPlayer{index: i}
		}
		return result
	}

	var filtered []filteredPlayer
	query := strings.ToLower(m.searchQuery)
	for i, player := range m.allPlayers {
		if strings.Contains(strings.ToLower(player.NickName), query) {
			filtered = append(filtered, filteredPlayer{index: i})
		}
	}
	return filtered
}

// getSelectedPlayersList возвращает список выбранных игроков
func (m Model) getSelectedPlayersList() []int {
	var selected []int
	for idx, isSelected := range m.selectedPlayers {
		if isSelected {
			selected = append(selected, idx)
		}
	}
	// ВАЖНО: Сортируем, чтобы порядок был стабильным!
	// Map в Go не гарантирует порядок итерации
	sort.Ints(selected)
	return selected
}
