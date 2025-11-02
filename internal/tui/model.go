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
	ScreenSettings
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

	// Индекс выбранного constraint для редактирования
	selectedConstraint int

	// Настройки
	numTeams      int     // 2 или 4
	sorryBro      *string // Игрок, который остается за бортом
	sorryBroIndex int     // Индекс для выбора sorryBro

	// Результаты генерации
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

	// Путь к конфигурационному файлу
	configPath string

	// Notifier'ы для отправки результатов
	notifiers []notifier.Notifier

	// Ошибка для отображения
	errorMsg string
}

// NewModel создает новую модель приложения
func NewModel(repo teambuilder.PlayerRepository, configPath string, notifiers []notifier.Notifier) Model {
	// Загружаем всех доступных игроков из репозитория
	allPlayers := loadAllPlayers(repo)

	return Model{
		currentScreen:   ScreenMenu,
		config:          &teambuilder.TeamConfiguration{},
		allPlayers:      allPlayers,
		selectedPlayers: make(map[int]bool),
		constraints:     []teambuilder.Constraint{},
		numTeams:        2,
		sorryBro:        nil,
		cursor:          0,
		menuCursor:      0,
		playerRepo:      repo,
		teamBuilder:     teambuilder.NewTeamBuilder(repo),
		configPath:      configPath,
		notifiers:       notifiers,
	}
}

// loadAllPlayers загружает всех игроков из репозитория
func loadAllPlayers(repo teambuilder.PlayerRepository) []teambuilder.TeamPlayer {
	// Получаем всех игроков из репозитория
	allPlayers := repo.GetAll()

	players := make([]teambuilder.TeamPlayer, 0, len(allPlayers))
	for _, player := range allPlayers {
		players = append(players, teambuilder.TeamPlayer{
			NickName: player.NickName,
			Score:    player.Score,
		})
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
		case "ctrl+c", "q":
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
	case ScreenSettings:
		return updateSettings(m, msg)
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
	case ScreenSettings:
		return viewSettings(m)
	case ScreenResults:
		return viewResults(m)
	}
	return ""
}
