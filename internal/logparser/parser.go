package logparser

import (
	"bufio"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Parser отвечает за парсинг log файлов
type Parser struct {
	regexps *LogRegexps
}

// New создает новый парсер
func New() *Parser {
	return &Parser{
		regexps: NewLogRegexps(),
	}
}

// ParseDirectory парсит все файлы в директории
func (p *Parser) ParseDirectory(dir, ext string) (*ParseResult, error) {
	result := &ParseResult{
		Players:      make(map[string]Player),
		KillEvents:   []KillEvent{},
		FlashEvents:  []FlashEvent{},
		DefuseEvents: []DefuseEvent{},
		WeaponSet:    make(map[string]struct{}),
		RoundStats:   []RoundStats{},
	}

	var fileDates []time.Time
	dateRegex := regexp.MustCompile(`(\d{4})_(\d{2})_(\d{2})_\d{6}`)

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if ext != "" && !strings.HasSuffix(strings.ToLower(d.Name()), strings.ToLower(ext)) {
			return nil
		}

		// Извлекаем дату из названия файла
		if matches := dateRegex.FindStringSubmatch(d.Name()); matches != nil {
			if date, err := time.Parse("2006-01-02", matches[1]+"-"+matches[2]+"-"+matches[3]); err == nil {
				fileDates = append(fileDates, date)
			}
		}

		return p.parseFile(path, result)
	})

	if err != nil {
		return result, err
	}

	// Определяем диапазон дат
	if len(fileDates) > 0 {
		sort.Slice(fileDates, func(i, j int) bool {
			return fileDates[i].Before(fileDates[j])
		})
		result.StartDate = fileDates[0].Format("02-01-2006")
		result.EndDate = fileDates[len(fileDates)-1].Format("02-01-2006")
	}

	return result, err
}

// MatchBounds представляет границы одного матча
type MatchBounds struct {
	StartLine int
	EndLine   int
}

// parseFile парсит один файл и добавляет результаты в ParseResult
func (p *Parser) parseFile(path string, result *ParseResult) error {
	// Сначала читаем весь файл чтобы найти границы матчей
	lines, err := p.readLines(path)
	if err != nil {
		return err
	}

	// Находим границы всех матчей в файле
	matches := p.findMatchBounds(lines)

	// Парсим события только внутри матчей
	for _, match := range matches {
		roundCountBefore := len(result.RoundStats)
		p.parseMatchLines(lines, match.StartLine, match.EndLine, result)
		// Пересчитываем рейтинги для раундов этого матча после того как Winner проставлен
		for i := roundCountBefore; i < len(result.RoundStats); i++ {
			calculateRoundRatings(&result.RoundStats[i])
		}
	}

	return nil
}

// readLines читает все строки из файла
func (p *Parser) readLines(path string) ([]string, error) {
	f, err := os.Open(path) // #nosec G304 - path is controlled by user input for log parsing
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// findMatchBounds находит границы всех матчей в файле
func (p *Parser) findMatchBounds(lines []string) []MatchBounds {
	var bounds []MatchBounds
	var currentMatchStart = -1

	for i, line := range lines {
		// Ищем начало матча
		if p.regexps.MatchStartPattern.MatchString(line) {
			currentMatchStart = i
			continue
		}

		// Ищем конец матча по событию "Game Over:"
		if p.regexps.GameOverPattern.MatchString(line) && currentMatchStart != -1 {
			bounds = append(bounds, MatchBounds{
				StartLine: currentMatchStart,
				EndLine:   i,
			})
			currentMatchStart = -1
		}
	}

	return bounds
}

// parseMatchLines парсит строки матча и добавляет события в result
func (p *Parser) parseMatchLines(lines []string, startLine, endLine int, result *ParseResult) {
	for i := startLine; i <= endLine && i < len(lines); i++ {
		line := lines[i]
		date := ExtractDateFromLogLine(line)

		// Проверяем, не начало ли это JSON_BEGIN блока
		if strings.Contains(line, "JSON_BEGIN{") {
			roundStats, consumed := p.parseJSONBlockFromLines(lines, i, date)
			if roundStats != nil {
				result.RoundStats = append(result.RoundStats, *roundStats)
			}
			i += consumed // Пропускаем обработанные строки
			continue
		}

		// Проверяем события победы команды
		if p.regexps.CTWinPattern.MatchString(line) {
			// CT выиграли - проставляем победителя последнему раунду
			if len(result.RoundStats) > 0 {
				result.RoundStats[len(result.RoundStats)-1].Winner = 3
			}
			continue
		}

		if p.regexps.TerroristWinPattern.MatchString(line) {
			// T выиграли - проставляем победителя последнему раунду
			if len(result.RoundStats) > 0 {
				result.RoundStats[len(result.RoundStats)-1].Winner = 2
			}
			continue
		}

		// Попытка парсинга убийства
		if matches := p.regexps.KillPattern.FindStringSubmatch(line); matches != nil {
			event := KillEvent{
				KillerName: matches[1],
				KillerSID:  matches[2],
				VictimName: matches[3],
				VictimSID:  matches[4],
				Weapon:     strings.TrimSpace(matches[5]),
				Date:       date,
			}
			result.KillEvents = append(result.KillEvents, event)

			if event.Weapon != "" {
				result.WeaponSet[event.Weapon] = struct{}{}
			}
			continue
		}

		// Попытка парсинга флешки
		if matches := p.regexps.FlashPattern.FindStringSubmatch(line); matches != nil {
			duration, _ := strconv.ParseFloat(matches[3], 64)
			event := FlashEvent{
				VictimName:  matches[1],
				VictimSID:   matches[2],
				FlasherName: matches[4],
				FlasherSID:  matches[5],
				Duration:    duration,
				Date:        date,
			}
			result.FlashEvents = append(result.FlashEvents, event)
			continue
		}

		// Попытка парсинга начала дефьюза
		if matches := p.regexps.DefuseBeginPattern.FindStringSubmatch(line); matches != nil {
			withKit := matches[3] == "With"
			event := DefuseEvent{
				PlayerName: matches[1],
				PlayerSID:  matches[2],
				WithKit:    withKit,
				EventType:  "begin",
				Date:       date,
			}
			result.DefuseEvents = append(result.DefuseEvents, event)
			continue
		}

		// Успешный дефьюз бомбы
		if matches := p.regexps.DefuseSuccessPattern.FindStringSubmatch(line); matches != nil {
			// Создаем событие успешного дефьюза (без указания конкретного игрока)
			event := DefuseEvent{
				PlayerName: "", // Будет определено при обработке
				PlayerSID:  "",
				WithKit:    false, // Будет определено при обработке
				EventType:  "success",
				Date:       date,
			}
			result.DefuseEvents = append(result.DefuseEvents, event)
			continue
		}

		// Брошенный дефьюз
		if matches := p.regexps.DefuseAbandonedPattern.FindStringSubmatch(line); matches != nil {
			event := DefuseEvent{
				PlayerName: matches[1],
				PlayerSID:  matches[2],
				WithKit:    false, // Будет определено при обработке
				EventType:  "abandoned",
				Date:       date,
			}
			result.DefuseEvents = append(result.DefuseEvents, event)
			continue
		}

		// Взрыв бомбы
		if matches := p.regexps.BombExplodedPattern.FindStringSubmatch(line); matches != nil {
			// Создаем событие взрыва бомбы
			event := DefuseEvent{
				PlayerName: "", // Будет определено при обработке
				PlayerSID:  "",
				WithKit:    false, // Будет определено при обработке
				EventType:  "failed",
				Date:       date,
			}
			result.DefuseEvents = append(result.DefuseEvents, event)
			continue
		}
	}
}

// parseJSONBlockFromLines парсит блок JSON_BEGIN...JSON_END из массива строк
// Возвращает RoundStats и количество обработанных строк
func (p *Parser) parseJSONBlockFromLines(lines []string, startIdx int, date string) (*RoundStats, int) {
	stats := &RoundStats{
		Date:    date,
		Players: []PlayerStats{},
	}

	firstLine := lines[startIdx]

	// Извлекаем время из первой строки
	if strings.HasPrefix(firstLine, "L ") {
		parts := strings.Fields(firstLine)
		if len(parts) >= 3 {
			stats.Time = strings.TrimSuffix(parts[2], ":")
		}
	}

	// Читаем строки до JSON_END
	consumed := 0
	var blockLines []string
	for i := startIdx + 1; i < len(lines); i++ {
		consumed++
		line := lines[i]
		if strings.Contains(line, "JSON_END") {
			break
		}
		blockLines = append(blockLines, line)
	}

	// Парсим метаданные
	for _, line := range blockLines {
		// Убираем префикс лога "L 09/05/2025 - 18:11:25: "
		if colonIdx := strings.Index(line, ": "); colonIdx != -1 {
			line = line[colonIdx+2:]
		}
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, `"round_number"`):
			if val := extractQuotedValue(line); val != "" {
				stats.RoundNumber, _ = strconv.Atoi(val)
			}
		case strings.HasPrefix(line, `"score_t"`):
			if val := extractQuotedValue(line); val != "" {
				stats.ScoreT, _ = strconv.Atoi(val)
			}
		case strings.HasPrefix(line, `"score_ct"`):
			if val := extractQuotedValue(line); val != "" {
				stats.ScoreCT, _ = strconv.Atoi(val)
			}
		case strings.HasPrefix(line, `"map"`):
			stats.Map = extractQuotedValue(line)
		case strings.HasPrefix(line, `"server"`):
			stats.Server = extractQuotedValue(line)
		case strings.HasPrefix(line, `"player_`):
			// Парсим статистику игрока
			playerStats := parsePlayerStats(line)
			if playerStats != nil {
				stats.Players = append(stats.Players, *playerStats)
			}
		}
	}

	// Рейтинги будут рассчитаны позже, после проставления Winner
	return stats, consumed
}

// calculateRoundRatings рассчитывает EPI рейтинг для всех игроков в раунде
func calculateRoundRatings(round *RoundStats) {
	if len(round.Players) == 0 {
		return
	}

	// Подсчитываем количество игроков в каждой команде
	ctCount := 0
	tCount := 0
	for _, p := range round.Players {
		switch p.Team {
		case 3:
			ctCount++
		case 2:
			tCount++
		}
	}

	// Рассчитываем рейтинг для каждого игрока
	for i := range round.Players {
		p := &round.Players[i]

		// Определяем параметры для формулы
		var oppCount, teamCount int
		var win float64

		switch p.Team {
		case 3: // CT
			oppCount = tCount
			teamCount = ctCount
			// Проверяем победу по полю Winner (будет проставлено позже)
			if round.Winner == 3 {
				win = 1.0
			}
		case 2: // T
			oppCount = ctCount
			teamCount = tCount
			// Проверяем победу по полю Winner (будет проставлено позже)
			if round.Winner == 2 {
				win = 1.0
			}
		}

		// Защита от деления на ноль
		if oppCount == 0 {
			oppCount = 5
		}
		if teamCount == 0 {
			teamCount = 5
		}

		// Формула EPI:
		// EPIraw = (Dmg/100) * (5/OppCount)^0.7 * (OppCount/TeamCount)^0.5 + 0.15*Kills + 0.08*Assists
		//          ----------------------------------------------------------------------------
		//                              1 + 0.35 * Deaths
		//          * (1 + 0.10 * Win + MultiKillBonus + ClutchBonus)

		dmg := float64(p.Damage)
		kills := float64(p.Kills)
		deaths := float64(p.Deaths)
		assists := float64(p.Assists)

		// Импакт от урона
		damageImpact := (dmg / 100.0) * math.Pow(5.0/float64(oppCount), 0.7) * math.Pow(float64(oppCount)/float64(teamCount), 0.5)

		// Импакт от убийств и ассистов
		fragImpact := 0.15*kills + 0.08*assists

		// Числитель
		numerator := damageImpact + fragImpact

		// Знаменатель
		denominator := 1.0 + 0.35*deaths

		// Базовый рейтинг
		baseRating := numerator / denominator

		// Бонус за многокиллы (процент убитых противников)
		killRatio := kills / float64(oppCount)
		var multiKillBonus float64
		switch {
		case killRatio >= 1.0: // Убил всех противников (ACE)
			multiKillBonus = 0.3
		case killRatio >= 0.8: // Убил 80%+ противников (например, 4 из 5)
			multiKillBonus = 0.2
		case killRatio >= 0.6: // Убил 60%+ противников (например, 3 из 5)
			multiKillBonus = 0.1
		}

		// Дополнительный бонус за клатч (в меньшинстве + победа + минимум 2 килла)
		clutchBonus := 0.0
		if teamCount < oppCount && win == 1.0 && kills >= 2 {
			// Бонус зависит от разницы в численности: каждый игрок в минусе дает +5%
			outnumberedDiff := float64(oppCount - teamCount)
			clutchBonus = outnumberedDiff * 0.05
		}

		// Финальный рейтинг
		p.Rating = baseRating * (1.0 + 0.10*win + multiKillBonus + clutchBonus)
	}
}

// ParseResult содержит результаты парсинга логов
type ParseResult struct {
	Players      map[string]Player
	KillEvents   []KillEvent
	FlashEvents  []FlashEvent
	DefuseEvents []DefuseEvent
	WeaponSet    map[string]struct{}
	RoundStats   []RoundStats // Статистика раундов из JSON_BEGIN блоков
	StartDate    string       // Начальная дата в формате DD-MM-YYYY
	EndDate      string       // Конечная дата в формате DD-MM-YYYY
}

// KeyAndTitle возвращает ключ и заголовок для группировки игроков
func KeyAndTitle(groupBy, name, sid string) (key, title string) {
	switch strings.ToLower(groupBy) {
	case "steamid":
		return sid, name // ключуем по SID, подписываем ником
	default:
		return name, name
	}
}

// extractQuotedValue извлекает значение из строки вида "key" : "value"
func extractQuotedValue(line string) string {
	parts := strings.Split(line, ":")
	if len(parts) < 2 {
		return ""
	}
	value := strings.TrimSpace(parts[1])
	value = strings.Trim(value, `",`)
	return value
}

// parsePlayerStats парсит строку со статистикой игрока
// Формат: "player_0" : "            26840160,      2,  16000,      0,      0,      0,      0,   0.00,   0.00,      0,      0,      0,      0,      0,      0,      0,      0,      0,      0,      0,      0,      0,      0,      0,      0,      0"
func parsePlayerStats(line string) *PlayerStats {
	// Находим значение в кавычках после :
	parts := strings.SplitN(line, ":", 2)
	if len(parts) < 2 {
		return nil
	}

	// Извлекаем CSV данные
	value := strings.TrimSpace(parts[1])
	value = strings.Trim(value, `"`)

	// Разбиваем по запятым
	fields := strings.Split(value, ",")
	if len(fields) < 26 {
		return nil
	}

	// Парсим каждое поле
	ps := &PlayerStats{}
	ps.AccountID, _ = strconv.ParseInt(strings.TrimSpace(fields[0]), 10, 64)
	ps.Team, _ = strconv.Atoi(strings.TrimSpace(fields[1]))
	ps.Money, _ = strconv.Atoi(strings.TrimSpace(fields[2]))
	ps.Kills, _ = strconv.Atoi(strings.TrimSpace(fields[3]))
	ps.Deaths, _ = strconv.Atoi(strings.TrimSpace(fields[4]))
	ps.Assists, _ = strconv.Atoi(strings.TrimSpace(fields[5]))

	// Damage может быть как int так и float (в новых версиях CS2)
	if damageFloat, err := strconv.ParseFloat(strings.TrimSpace(fields[6]), 64); err == nil {
		ps.Damage = int(damageFloat)
	}

	ps.HSP, _ = strconv.ParseFloat(strings.TrimSpace(fields[7]), 64)
	ps.KDR, _ = strconv.ParseFloat(strings.TrimSpace(fields[8]), 64)
	ps.ADR, _ = strconv.ParseFloat(strings.TrimSpace(fields[9]), 64)
	ps.MVP, _ = strconv.Atoi(strings.TrimSpace(fields[10]))
	ps.EF, _ = strconv.Atoi(strings.TrimSpace(fields[11]))
	ps.UD, _ = strconv.Atoi(strings.TrimSpace(fields[12]))
	ps.ThreeK, _ = strconv.Atoi(strings.TrimSpace(fields[13]))
	ps.FourK, _ = strconv.Atoi(strings.TrimSpace(fields[14]))
	ps.FiveK, _ = strconv.Atoi(strings.TrimSpace(fields[15]))
	ps.ClutchK, _ = strconv.Atoi(strings.TrimSpace(fields[16]))
	ps.FirstK, _ = strconv.Atoi(strings.TrimSpace(fields[17]))
	ps.PistolK, _ = strconv.Atoi(strings.TrimSpace(fields[18]))
	ps.SniperK, _ = strconv.Atoi(strings.TrimSpace(fields[19]))
	ps.BlindK, _ = strconv.Atoi(strings.TrimSpace(fields[20]))
	ps.BombK, _ = strconv.Atoi(strings.TrimSpace(fields[21]))

	// FireDmg и UniqueK тоже могут быть float в новых версиях
	if fireDmgFloat, err := strconv.ParseFloat(strings.TrimSpace(fields[22]), 64); err == nil {
		ps.FireDmg = int(fireDmgFloat)
	}
	if uniqueKFloat, err := strconv.ParseFloat(strings.TrimSpace(fields[23]), 64); err == nil {
		ps.UniqueK = int(uniqueKFloat)
	}

	ps.Dinks, _ = strconv.Atoi(strings.TrimSpace(fields[24]))
	ps.ChickenK, _ = strconv.Atoi(strings.TrimSpace(fields[25]))

	return ps
}
