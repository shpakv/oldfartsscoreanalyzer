package logparser

import (
	"regexp"
	"strings"
)

// Player представляет игрока с ключом группировки и отображаемым именем
type Player struct {
	Key   string // ключ группировки (ник или SteamID)
	Title string // подпись (обычно ник)
}

// KillEvent представляет событие убийства
type KillEvent struct {
	KillerName string
	KillerSID  string
	VictimName string
	VictimSID  string
	Weapon     string
	Date       string // Дата в формате YYYY-MM-DD
}

// FlashEvent представляет событие ослепления
type FlashEvent struct {
	FlasherName string
	FlasherSID  string
	VictimName  string
	VictimSID   string
	Duration    float64
	Date        string // Дата в формате YYYY-MM-DD
}

// DefuseEvent представляет событие дефьюза бомбы
type DefuseEvent struct {
	PlayerName string
	PlayerSID  string
	WithKit    bool   // true если с дефьюз-китом, false если без кита
	EventType  string // "begin", "success", "abandoned", "failed"
	Date       string // Дата в формате YYYY-MM-DD
}

// RoundStats представляет статистику раунда из JSON_BEGIN блока
type RoundStats struct {
	Date        string        // Дата в формате YYYY-MM-DD
	Time        string        // Время в формате HH:MM:SS
	RoundNumber int           // Номер раунда
	ScoreT      int           // Счёт террористов
	ScoreCT     int           // Счёт контр-террористов
	Map         string        // Название карты
	Server      string        // Название сервера
	Players     []PlayerStats // Статистика игроков
	Winner      int           // Победитель раунда: 2=T, 3=CT, 0=неизвестно/ничья
}

// PlayerStats представляет статистику одного игрока в раунде
type PlayerStats struct {
	AccountID int64   // Steam Account ID
	Team      int     // 2=T, 3=CT
	Money     int     // Деньги
	Kills     int     // Убийства
	Deaths    int     // Смерти
	Assists   int     // Ассисты
	Damage    int     // Урон
	HSP       float64 // Headshot percentage
	KDR       float64 // Kill/Death Ratio
	ADR       float64 // Average Damage per Round
	MVP       int     // MVP звёзды
	EF        int     // Entry Frags
	UD        int     // Utility Damage
	ThreeK    int     // 3 kills
	FourK     int     // 4 kills
	FiveK     int     // 5 kills (ace)
	ClutchK   int     // Clutch kills
	FirstK    int     // First kill
	PistolK   int     // Pistol kills
	SniperK   int     // Sniper kills
	BlindK    int     // Blind kills
	BombK     int     // Bomb defuses/plants
	FireDmg   int     // Fire damage
	UniqueK   int     // Unique kills
	Dinks     int     // Headshot hits
	ChickenK  int     // Chicken kills
	Rating    float64 // EPI Rating (рассчитывается после парсинга)
}

// LogRegexps содержит регулярные выражения для парсинга логов
type LogRegexps struct {
	KillPattern            *regexp.Regexp
	FlashPattern           *regexp.Regexp
	DefuseBeginPattern     *regexp.Regexp
	DefuseSuccessPattern   *regexp.Regexp
	DefuseAbandonedPattern *regexp.Regexp
	BombExplodedPattern    *regexp.Regexp
	MatchStartPattern      *regexp.Regexp
	MatchStatusPattern     *regexp.Regexp
	GameOverPattern        *regexp.Regexp
	CTWinPattern           *regexp.Regexp
	TerroristWinPattern    *regexp.Regexp
}

// NewLogRegexps создает новые регулярные выражения для парсинга CS2 логов
func NewLogRegexps() *LogRegexps {
	killRe := regexp.MustCompile(
		`^L\s+\d{2}/\d{2}/\d{4}\s+-\s+\d{2}:\d{2}:\d{2}:\s"` +
			`([^"<]+)<\d+><([^>]+)><[^>]*>"\s+\[[^\]]+\]\s+killed\s+` + // killerName, killerSID
			`"([^"<]+)<\d+><([^>]+)><[^>]*>"\s+\[[^\]]+\]\s+with\s+"([^"]+)"$`) // victimName, victimSID, weapon

	flashRe := regexp.MustCompile(
		`^L\s+\d{2}/\d{2}/\d{4}\s+-\s+\d{2}:\d{2}:\d{2}:\s"` +
			`([^"<]+)<\d+><([^>]+)><[^>]*>"\s+blinded\s+for\s+([0-9.]+)\s+by\s+` + // victimName, victimSID, duration
			`"([^"<]+)<\d+><([^>]+)><[^>]*>"\s+from\s+flashbang\s+entindex\s+\d+\s*$`) // flasherName, flasherSID

	// Пример: "povidlo boy<4><[U:1:44922694]><CT>" triggered "Begin_Bomb_Defuse_With_Kit"
	defuseBeginRe := regexp.MustCompile(
		`^L\s+\d{2}/\d{2}/\d{4}\s+-\s+\d{2}:\d{2}:\d{2}:\s"` +
			`([^"<]+)<\d+><([^>]+)><[^>]*>"\s+triggered\s+"Begin_Bomb_Defuse_(With|Without)_Kit"$`) // playerName, playerSID, withKit

	// Пример: Team "CT" triggered "SFUI_Notice_Bomb_Defused" (CT "8") (T "5")
	defuseSuccessRe := regexp.MustCompile(
		`^L\s+\d{2}/\d{2}/\d{4}\s+-\s+\d{2}:\d{2}:\d{2}:\s+Team\s+"CT"\s+triggered\s+"SFUI_Notice_Bomb_Defused"`)

	// Пример: "player<id><steamid><CT>" stopped defusing the bomb
	defuseAbandonedRe := regexp.MustCompile(
		`^L\s+\d{2}/\d{2}/\d{4}\s+-\s+\d{2}:\d{2}:\d{2}:\s"` +
			`([^"<]+)<\d+><([^>]+)><[^>]*>"\s+stopped\s+defusing\s+the\s+bomb`)

	// Пример: Team "TERRORIST" triggered "SFUI_Notice_Target_Bombed" (CT "8") (T "5")
	bombExplodedRe := regexp.MustCompile(
		`^L\s+\d{2}/\d{2}/\d{4}\s+-\s+\d{2}:\d{2}:\d{2}:\s+Team\s+"TERRORIST"\s+triggered\s+"SFUI_Notice_Target_Bombed"`)

	// Пример: World triggered "Match_Start" on "cs_office"
	matchStartRe := regexp.MustCompile(
		`^L\s+\d{2}/\d{2}/\d{4}\s+-\s+\d{2}:\d{2}:\d{2}:\s+World\s+triggered\s+"Match_Start"`)

	// Пример: MatchStatus: Score: 13:6 on map "cs_office" RoundsPlayed: 19
	matchStatusRe := regexp.MustCompile(
		`^L\s+\d{2}/\d{2}/\d{4}\s+-\s+\d{2}:\d{2}:\d{2}:\s+MatchStatus:\s+Score:\s+(\d+):(\d+)\s+on\s+map\s+"[^"]+"\s+RoundsPlayed:\s+(-?\d+)`)

	// Пример: Game Over: competitive cs_office score 13:6 after 28 min
	gameOverRe := regexp.MustCompile(
		`^L\s+\d{2}/\d{2}/\d{4}\s+-\s+\d{2}:\d{2}:\d{2}:\s+Game Over:`)

	// Пример: Team "CT" triggered "SFUI_Notice_CTs_Win" (CT "1") (T "0")
	ctWinRe := regexp.MustCompile(
		`^L\s+\d{2}/\d{2}/\d{4}\s+-\s+\d{2}:\d{2}:\d{2}:\s+Team\s+"CT"\s+triggered\s+"SFUI_Notice_CTs_Win"`)

	// Пример: Team "TERRORIST" triggered "SFUI_Notice_Terrorists_Win" (CT "0") (T "1")
	terroristWinRe := regexp.MustCompile(
		`^L\s+\d{2}/\d{2}/\d{4}\s+-\s+\d{2}:\d{2}:\d{2}:\s+Team\s+"TERRORIST"\s+triggered\s+"SFUI_Notice_Terrorists_Win"`)

	return &LogRegexps{
		KillPattern:            killRe,
		FlashPattern:           flashRe,
		DefuseBeginPattern:     defuseBeginRe,
		DefuseSuccessPattern:   defuseSuccessRe,
		DefuseAbandonedPattern: defuseAbandonedRe,
		BombExplodedPattern:    bombExplodedRe,
		MatchStartPattern:      matchStartRe,
		MatchStatusPattern:     matchStatusRe,
		GameOverPattern:        gameOverRe,
		CTWinPattern:           ctWinRe,
		TerroristWinPattern:    terroristWinRe,
	}
}

// ExtractDateFromLogLine извлекает дату из строки лога в формате YYYY-MM-DD
// Пример входной строки: L 02/14/2024 - 12:34:56: ...
func ExtractDateFromLogLine(line string) string {
	// Ищем паттерн даты в начале строки: L MM/DD/YYYY
	if !strings.HasPrefix(line, "L ") {
		return ""
	}

	// Находим позицию первой цифры после "L "
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return ""
	}

	datePart := parts[1] // Должно быть MM/DD/YYYY
	dateComponents := strings.Split(datePart, "/")
	if len(dateComponents) != 3 {
		return ""
	}

	// Преобразуем MM/DD/YYYY -> YYYY-MM-DD
	month := dateComponents[0]
	day := dateComponents[1]
	year := dateComponents[2]

	// Добавляем ведущий ноль если нужно
	if len(month) == 1 {
		month = "0" + month
	}
	if len(day) == 1 {
		day = "0" + day
	}

	return year + "-" + month + "-" + day
}
