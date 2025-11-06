package tui

import (
	"oldfartscounter/internal/notifier"
	"oldfartscounter/internal/teambuilder"

	tea "github.com/charmbracelet/bubbletea"
)

// Screen представляет текущий экран приложения
type Screen int

const (
	ScreenMenu Screen = iota
	ScreenPlayers
	ScreenConstraints
	ScreenResults
)

// Model основная модель приложения
type Model struct {
	// Текущий экран
	currentScreen Screen

	// Конфигурация
	config *teambuilder.TeamConfiguration

	// Список всех доступных игроков (из репозитория)
	allPlayers []teambuilder.TeamPlayer

	// Выбранные игроки (индексы в allPlayers)
	selectedPlayers map[int]bool

	// Текущий курсор для навигации
	cursor int

	// Поисковый запрос для экрана игроков
	searchQuery string
	// Режим поиска (true = печатаем, false = навигация)
	searchMode bool

	// Constraints
	constraints []teambuilder.Constraint

	// SorryBro - игрок, который остается за бортом
	sorryBro *string

	// Результаты генерации (всегда 2 команды)
	generatedTeams []teambuilder.Team

	// Размеры окна
	width  int
	height int

	// Флаг режима редактирования constraint
	editingConstraint    bool
	editingConstraintNew bool // true если создаем новый constraint
	constraintPlayer1Idx int
	constraintPlayer2Idx int
	constraintType       teambuilder.ConstrainType
	constraintFieldFocus int // 0 - player1, 1 - player2, 2 - type

	// Подэкраны для навигации в меню
	menuCursor int

	// Репозиторий игроков
	playerRepo teambuilder.PlayerRepository

	// TeamBuilder
	teamBuilder *teambuilder.TeamBuilder

	// Notifier'ы для отправки результатов
	notifiers []notifier.Notifier

	// Ошибка для отображения
	errorMsg string
	// Успешное сообщение для отображения
	successMsg string
}

// NewModel создает новую модель приложения
func NewModel(repo teambuilder.PlayerRepository, notifiers []notifier.Notifier) Model {
	// Загружаем всех доступных игроков из репозитория
	allPlayers := loadAllPlayers(repo)

	return Model{
		currentScreen:   ScreenMenu,
		config:          &teambuilder.TeamConfiguration{},
		allPlayers:      allPlayers,
		selectedPlayers: make(map[int]bool),
		constraints:     []teambuilder.Constraint{},
		sorryBro:        nil,
		cursor:          0,
		menuCursor:      0,
		playerRepo:      repo,
		teamBuilder:     teambuilder.NewTeamBuilder(repo),
		notifiers:       notifiers,
	}
}

// loadAllPlayers загружает всех игроков из репозитория
func loadAllPlayers(repo teambuilder.PlayerRepository) []teambuilder.TeamPlayer {
	// Получаем всех игроков из репозитория
	allPlayers := repo.GetAll()

	players := make([]teambuilder.TeamPlayer, 0, len(allPlayers))
	for _, player := range allPlayers {
		players = append(players, teambuilder.TeamPlayer(player))
	}

	return players
}

// Init инициализирует приложение
func (m Model) Init() tea.Cmd {
	return nil
}

// Update обрабатывает сообщения
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case keyCtrlC, "q":
			if m.currentScreen == ScreenMenu {
				return m, tea.Quit
			}
		}
	}

	// Делегируем обработку текущему экрану
	switch m.currentScreen {
	case ScreenMenu:
		return updateMenu(m, msg)
	case ScreenPlayers:
		return updatePlayers(m, msg)
	case ScreenConstraints:
		return updateConstraints(m, msg)
	case ScreenResults:
		return updateResults(m, msg)
	}

	return m, nil
}

// View рендерит UI
func (m Model) View() string {
	switch m.currentScreen {
	case ScreenMenu:
		return viewMenu(m)
	case ScreenPlayers:
		return viewPlayers(m)
	case ScreenConstraints:
		return viewConstraints(m)
	case ScreenResults:
		return viewResults(m)
	}
	return ""
}

// generateTeams генерирует 2 команды на основе текущей конфигурации
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
		NumTeams:    2, // Всегда 2 команды
		SorryBro:    m.sorryBro,
	}

	// Генерируем команды
	m.generatedTeams = m.teamBuilder.Build(m.config)
}

// getSelectedPlayersList возвращает список выбранных игроков (нужен для generateTeams)
func (m Model) getSelectedPlayersList() []int {
	var selected []int
	for idx, isSelected := range m.selectedPlayers {
		if isSelected {
			selected = append(selected, idx)
		}
	}
	// Сортируем для стабильного порядка
	for i := 0; i < len(selected)-1; i++ {
		for j := i + 1; j < len(selected); j++ {
			if selected[i] > selected[j] {
				selected[i], selected[j] = selected[j], selected[i]
			}
		}
	}
	return selected
}

// calculateAverageMu возвращает μ для расчета категорий рейтинга
// μ берется из репозитория, который получает его из реальных логов или БД
func (m Model) calculateAverageMu() float64 {
	return m.playerRepo.GetAverageMu()
}

// SetSorryBro устанавливает SorryBro (например, из переменной окружения)
func (m *Model) SetSorryBro(nickname string) {
	m.sorryBro = &nickname
}
