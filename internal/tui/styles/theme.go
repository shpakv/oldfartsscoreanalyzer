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

// Константы для категорий рейтинга
const (
	CategoryMonster = "Гиперебака"
	CategoryGood    = "Ебака"
	CategoryAverage = "Пердун"
	CategoryWeak    = "Подпивас"
)

// hslToRGB конвертирует HSL в RGB
func hslToRGB(h, s, l float64) (r, g, b uint8) {
	hue2rgb := func(p, q, t float64) float64 {
		if t < 0 {
			t += 1
		}
		if t > 1 {
			t -= 1
		}
		if t < 1.0/6.0 {
			return p + (q-p)*6*t
		}
		if t < 1.0/2.0 {
			return q
		}
		if t < 2.0/3.0 {
			return p + (q-p)*(2.0/3.0-t)*6
		}
		return p
	}

	var rf, gf, bf float64
	if s == 0 {
		rf, gf, bf = l, l, l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q
		rf = hue2rgb(p, q, h+1.0/3.0)
		gf = hue2rgb(p, q, h)
		bf = hue2rgb(p, q, h-1.0/3.0)
	}

	r = uint8(rf * 255)
	g = uint8(gf * 255)
	b = uint8(bf * 255)
	return r, g, b
}

// GetRatingColor возвращает цвет в зависимости от рейтинга (ТОЧНО как в HTML)
// Формула: h=(220*(1-t))/360, s=0.85, l=0.16+0.39*t
func GetRatingColor(rating float64, maxRating float64) lipgloss.Color {
	if maxRating == 0 {
		return lipgloss.Color("#0b1020")
	}

	// Нормализуем рейтинг от 0 до 1
	t := rating / maxRating
	if t > 1.0 {
		t = 1.0
	}
	if t < 0 {
		t = 0
	}

	// Точно та же формула что и в HTML:
	// h = (220*(1-t))/360  -> от 220° (синий) до 0° (красный)
	// s = 0.85             -> насыщенность 85%
	// l = 0.16 + 0.39*t    -> яркость от 16% до 55%
	h := (220 * (1 - t)) / 360.0
	s := 0.85
	l := 0.16 + 0.39*t

	r, g, b := hslToRGB(h, s, l)
	return lipgloss.Color(formatHex(r, g, b))
}

// formatHex форматирует RGB в hex строку
func formatHex(r, g, b uint8) string {
	return sprintf("#%02x%02x%02x", r, g, b)
}

func sprintf(format string, args ...interface{}) string {
	// Простая реализация для hex формата
	if format == "#%02x%02x%02x" && len(args) == 3 {
		r, _ := args[0].(uint8)
		g, _ := args[1].(uint8)
		b, _ := args[2].(uint8)

		hexChars := "0123456789abcdef"
		result := "#"
		result += string(hexChars[r>>4]) + string(hexChars[r&0x0f])
		result += string(hexChars[g>>4]) + string(hexChars[g&0x0f])
		result += string(hexChars[b>>4]) + string(hexChars[b&0x0f])
		return result
	}
	return ""
}

// GetRatingCategory возвращает категорию рейтинга (как в HTML)
// rating - это байесовский рейтинг (BayesianEPI) из логов или repository.go
// averageMu - средний байесовский рейтинг (из repository.go или логов)
// В TUI averageMu берется из среднего значения всех Score в repository.go
func GetRatingCategory(rating float64, averageMu float64) (category string, color lipgloss.Color) {
	// Используем те же пороги что и в HTML, но применяем к средней по всем игрокам
	// Это позволяет категориям адаптироваться к текущему распределению рейтингов
	mu := averageMu
	if mu == 0 {
		mu = 0.6 // Дефолтное значение, если averageMu не передан
	}

	weak := mu * 0.85    // Подпивас (Mil-Spec)
	average := mu * 1.05 // Пердун (Restricted)
	monster := mu * 1.25 // Гиперебака (Covert)

	switch {
	case rating >= monster:
		return CategoryMonster, lipgloss.Color("#cfb53b") // gold - Covert
	case rating >= average:
		return CategoryGood, lipgloss.Color("#ef4444") // red - Classified
	case rating >= weak:
		return CategoryAverage, lipgloss.Color("#4b69ff") // blue - Restricted
	default:
		return CategoryWeak, lipgloss.Color("#9ca3af") // gray - Mil-Spec
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
