package styles

import "github.com/charmbracelet/lipgloss"

// Цветовая схема из HTML (cs2_stats.html)
var (
	// Основные цвета
	ColorBg     = lipgloss.Color("#1e1e1e")
	ColorPanel  = lipgloss.Color("#2b2b2b")
	ColorPanel2 = lipgloss.Color("#232323")
	ColorText   = lipgloss.Color("#c8c8c8")
	ColorMuted  = lipgloss.Color("#9aa0a6")
	ColorAccent = lipgloss.Color("#7c5cff")
	ColorGrid   = lipgloss.Color("#3a3a3a")
	ColorSticky = lipgloss.Color("#242424")

	// Дополнительные цвета
	ColorSuccess = lipgloss.Color("#22c55e")
	ColorWarning = lipgloss.Color("#f59e0b")
	ColorError   = lipgloss.Color("#ef4444")
	ColorInfo    = lipgloss.Color("#0ea5e9")
)

// Стили для основных элементов
var (
	// Стиль заголовка
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Bold(true).
			Padding(0, 1)

	// Стиль подзаголовка
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Italic(true)

	// Стиль панели
	PanelStyle = lipgloss.NewStyle().
			Background(ColorPanel).
			Foreground(ColorText).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorGrid).
			Padding(1, 2)

	// Стиль активного элемента
	ActiveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fff")).
			Background(lipgloss.Color("#2a2440")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorAccent).
			Padding(0, 1)

	// Стиль неактивного элемента
	InactiveStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Background(lipgloss.Color("transparent")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorGrid).
			Padding(0, 1)

	// Стиль помощи (футер)
	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Italic(true).
			Padding(0, 1)

	// Стиль для выбранного пункта списка
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(ColorAccent).
				Bold(true).
				PaddingLeft(1)

	// Стиль для невыбранного пункта списка
	UnselectedItemStyle = lipgloss.NewStyle().
				Foreground(ColorText).
				PaddingLeft(1)

	// Стиль для чекбокса (выбран)
	CheckboxCheckedStyle = lipgloss.NewStyle().
				Foreground(ColorAccent).
				Bold(true)

	// Стиль для чекбокса (не выбран)
	CheckboxUncheckedStyle = lipgloss.NewStyle().
				Foreground(ColorMuted)

	// Стиль для разделителя
	SeparatorStyle = lipgloss.NewStyle().
			Foreground(ColorGrid)

	// Стиль для статистики (счет команд)
	StatStyle = lipgloss.NewStyle().
			Foreground(ColorInfo).
			Bold(true)

	// Стиль для акцентного текста
	AccentTextStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	// Стиль для ошибок
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	// Стиль для успеха
	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	// Стиль для предупреждений
	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorWarning).
			Bold(true)
)

// GetRatingColor возвращает цвет в зависимости от рейтинга (градиент как в HTML)
func GetRatingColor(rating float64, maxRating float64) lipgloss.Color {
	// Нормализуем рейтинг от 0 до 1
	normalized := rating / maxRating
	if normalized > 1.0 {
		normalized = 1.0
	}

	// Градиент от темно-фиолетового до красного (как в HTML)
	// #2d1b4e -> #4b69ff -> #0ea5e9 -> #22c55e -> #cfb53b -> #f59e0b -> #ef4444
	switch {
	case normalized < 0.16:
		return lipgloss.Color("#2d1b4e") // Темно-фиолетовый
	case normalized < 0.33:
		return lipgloss.Color("#4b69ff") // Синий
	case normalized < 0.50:
		return lipgloss.Color("#0ea5e9") // Голубой
	case normalized < 0.66:
		return lipgloss.Color("#22c55e") // Зеленый
	case normalized < 0.83:
		return lipgloss.Color("#cfb53b") // Желтый
	case normalized < 0.95:
		return lipgloss.Color("#f59e0b") // Оранжевый
	default:
		return lipgloss.Color("#ef4444") // Красный
	}
}

// RenderProgressBar рендерит прогресс-бар для рейтинга
func RenderProgressBar(rating float64, maxRating float64, width int) string {
	normalized := rating / maxRating
	if normalized > 1.0 {
		normalized = 1.0
	}

	filled := int(normalized * float64(width))
	empty := width - filled

	color := GetRatingColor(rating, maxRating)
	barStyle := lipgloss.NewStyle().Foreground(color)

	bar := ""
	for i := 0; i < filled; i++ {
		bar += "█"
	}
	for i := 0; i < empty; i++ {
		bar += "░"
	}

	return barStyle.Render(bar)
}
