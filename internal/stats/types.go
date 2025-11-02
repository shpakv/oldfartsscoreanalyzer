package stats

import "oldfartscounter/internal/logparser"

// StatsData содержит обработанные статистики
type StatsData struct {
	Players            []Player
	Weapons            []string
	KillMatrix         KillMatrix
	WeaponData         WeaponData
	FlashData          FlashData
	DefuseData         DefuseData
	DateRange          string  // Период данных в формате "DD-MM-YYYY - DD-MM-YYYY"
	HighlightedPlayer  string  // Игрок для золотой подсветки в табе "Сорян, Братан"
	MinRoundsForRating float64 // Минимальное количество раундов для достоверного рейтинга (K)
	AverageMu          float64 // Средний EPI всех игроков (μ) - рассчитывается из реальных данных
	KillEvents         []logparser.KillEvent
	FlashEvents        []logparser.FlashEvent
	DefuseEvents       []logparser.DefuseEvent
	RoundStats         []logparser.RoundStats // Статистика раундов
	PlayerRatings      []PlayerRating         // Агрегированные рейтинги игроков
	// Агрегированные данные по датам для оптимизации
	DailyKills  map[string][]logparser.KillEvent   // дата -> события
	DailyFlash  map[string][]logparser.FlashEvent  // дата -> события
	DailyDefuse map[string][]logparser.DefuseEvent // дата -> события
	DailyRounds map[string][]logparser.RoundStats  // дата -> раунды
}

// Player представляет игрока
type Player struct {
	Key   string
	Title string
}

// KillMatrix содержит матрицу убийств
type KillMatrix struct {
	Matrix [][]int
	Max    int
}

// WeaponData содержит данные по оружию
type WeaponData struct {
	KillerWeaponMatrix [][]int // Players × Weapons
	VictimWeaponMatrix [][]int // Players × Weapons
	WeaponKillsMatrix  [][]int // Weapons × Players (транспонированная)
	KillerMax          int
	VictimMax          int
}

// FlashData содержит данные по флешкам
type FlashData struct {
	CountMatrix   [][]int
	SecondsMatrix [][]float64
	CountMax      int
	SecondsMax    float64
}

// DefuseData содержит данные по дефьюзу
type DefuseData struct {
	Attempts          []int // общее количество попыток дефьюза по игрокам
	WithKit           []int // количество попыток с китом по игрокам
	WithoutKit        []int // количество попыток без кита по игрокам
	SuccessWithKit    []int // успешные дефьюзы с китом
	SuccessWithoutKit []int // успешные дефьюзы без кита
	Abandoned         []int // брошенные дефьюзы
	Failed            []int // не успел (бомба взорвалась в процессе)
	TotalMax          int   // максимальное количество попыток среди всех игроков
}

// PlayerRating содержит агрегированный рейтинг игрока
type PlayerRating struct {
	AccountID    int64   // Steam Account ID
	Name         string  // Имя игрока
	RoundsPlayed int     // Количество сыгранных раундов
	TotalEPI     float64 // Сумма EPI по всем раундам
	AverageEPI   float64 // Простое среднее EPI
	BayesianEPI  float64 // Байесовский рейтинг с регуляризацией
	TotalDamage  int     // Общий урон
	TotalKills   int     // Общие убийства
	TotalDeaths  int     // Общие смерти
	TotalAssists int     // Общие ассисты
	WinRounds    int     // Выигранные раунды
	LastPlayed   string  // Дата последней игры
}
