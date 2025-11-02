package stats

import (
	"fmt"
	"sort"
	"strings"

	"oldfartscounter/internal/logparser"
)

// Processor обрабатывает данные парсинга и создает статистику
type Processor struct{}

// New создает новый процессор статистики
func New() *Processor {
	return &Processor{}
}

// Process обрабатывает результаты парсинга и возвращает статистические данные
func (p *Processor) Process(parseResult *logparser.ParseResult, groupBy string) *StatsData {
	// Заполняем игроков
	players := make(map[string]Player)

	// Сортируем события убийств по дате, чтобы получить самые свежие ники
	sortedKillEvents := make([]logparser.KillEvent, len(parseResult.KillEvents))
	copy(sortedKillEvents, parseResult.KillEvents)
	sort.Slice(sortedKillEvents, func(i, j int) bool {
		return sortedKillEvents[i].Date < sortedKillEvents[j].Date
	})

	// Обрабатываем события убийств (всегда обновляем Title на последний)
	for _, event := range sortedKillEvents {
		kKey, kTitle := logparser.KeyAndTitle(groupBy, event.KillerName, event.KillerSID)
		vKey, vTitle := logparser.KeyAndTitle(groupBy, event.VictimName, event.VictimSID)

		// Всегда обновляем Title - так получим последний (актуальный) ник
		players[kKey] = Player{Key: kKey, Title: kTitle}
		players[vKey] = Player{Key: vKey, Title: vTitle}
	}

	// Сортируем события флешек по дате
	sortedFlashEvents := make([]logparser.FlashEvent, len(parseResult.FlashEvents))
	copy(sortedFlashEvents, parseResult.FlashEvents)
	sort.Slice(sortedFlashEvents, func(i, j int) bool {
		return sortedFlashEvents[i].Date < sortedFlashEvents[j].Date
	})

	// Обрабатываем события флешек (всегда обновляем Title на последний)
	for _, event := range sortedFlashEvents {
		fKey, fTitle := logparser.KeyAndTitle(groupBy, event.FlasherName, event.FlasherSID)
		vKey, vTitle := logparser.KeyAndTitle(groupBy, event.VictimName, event.VictimSID)

		// Всегда обновляем Title - так получим последний (актуальный) ник
		players[fKey] = Player{Key: fKey, Title: fTitle}
		players[vKey] = Player{Key: vKey, Title: vTitle}
	}

	// Сортируем события дефьюза по дате
	sortedDefuseEvents := make([]logparser.DefuseEvent, len(parseResult.DefuseEvents))
	copy(sortedDefuseEvents, parseResult.DefuseEvents)
	sort.Slice(sortedDefuseEvents, func(i, j int) bool {
		return sortedDefuseEvents[i].Date < sortedDefuseEvents[j].Date
	})

	// Обрабатываем события дефьюза (всегда обновляем Title на последний)
	for _, event := range sortedDefuseEvents {
		pKey, pTitle := logparser.KeyAndTitle(groupBy, event.PlayerName, event.PlayerSID)
		players[pKey] = Player{Key: pKey, Title: pTitle}
	}

	// Создаем упорядоченный список игроков
	playerList := make([]Player, 0, len(players))
	for _, player := range players {
		playerList = append(playerList, player)
	}
	sort.Slice(playerList, func(i, j int) bool {
		return strings.ToLower(playerList[i].Title) < strings.ToLower(playerList[j].Title)
	})

	// Создаем список оружия
	weapons := make([]string, 0, len(parseResult.WeaponSet))
	for weapon := range parseResult.WeaponSet {
		weapons = append(weapons, weapon)
	}
	sort.Slice(weapons, func(i, j int) bool {
		return strings.ToLower(weapons[i]) < strings.ToLower(weapons[j])
	})

	// Создаем индексы для быстрого поиска
	playerIndex := make(map[string]int)
	for i, player := range playerList {
		playerIndex[player.Key] = i
	}
	weaponIndex := make(map[string]int)
	for i, weapon := range weapons {
		weaponIndex[weapon] = i
	}

	// Формируем строку периода данных
	dateRange := ""
	if parseResult.StartDate != "" && parseResult.EndDate != "" {
		if parseResult.StartDate == parseResult.EndDate {
			dateRange = parseResult.StartDate
		} else {
			dateRange = parseResult.StartDate + " — " + parseResult.EndDate
		}
	}

	// Группируем события по датам
	dailyKills := make(map[string][]logparser.KillEvent)
	for _, e := range parseResult.KillEvents {
		if e.Date != "" {
			dailyKills[e.Date] = append(dailyKills[e.Date], e)
		}
	}

	dailyFlash := make(map[string][]logparser.FlashEvent)
	for _, e := range parseResult.FlashEvents {
		if e.Date != "" {
			dailyFlash[e.Date] = append(dailyFlash[e.Date], e)
		}
	}

	dailyDefuse := make(map[string][]logparser.DefuseEvent)
	for _, e := range parseResult.DefuseEvents {
		if e.Date != "" {
			dailyDefuse[e.Date] = append(dailyDefuse[e.Date], e)
		}
	}

	dailyRounds := make(map[string][]logparser.RoundStats)
	for _, r := range parseResult.RoundStats {
		if r.Date != "" {
			dailyRounds[r.Date] = append(dailyRounds[r.Date], r)
		}
	}

	// Строим рейтинги
	playerRatings := p.buildPlayerRatings(parseResult.RoundStats, parseResult.KillEvents, parseResult.FlashEvents, parseResult.DefuseEvents)

	// Вычисляем средний EPI (μ) из реальных данных
	var totalEPI float64
	var totalRounds int
	for _, rating := range playerRatings {
		totalEPI += rating.TotalEPI
		totalRounds += rating.RoundsPlayed
	}
	averageMu := 0.6 // Дефолтное значение
	if totalRounds > 0 {
		averageMu = totalEPI / float64(totalRounds)
	}

	return &StatsData{
		Players:            playerList,
		Weapons:            weapons,
		KillMatrix:         p.buildKillMatrix(parseResult.KillEvents, playerList, playerIndex, groupBy),
		WeaponData:         p.buildWeaponData(parseResult.KillEvents, playerList, weapons, playerIndex, weaponIndex, groupBy),
		FlashData:          p.buildFlashData(parseResult.FlashEvents, playerList, playerIndex, groupBy),
		DefuseData:         p.buildDefuseData(parseResult.DefuseEvents, playerList, playerIndex, groupBy),
		DateRange:          dateRange,
		MinRoundsForRating: 100.0,     // Константа K для байесовского рейтинга
		AverageMu:          averageMu, // Средний EPI всех игроков
		KillEvents:         parseResult.KillEvents,
		FlashEvents:        parseResult.FlashEvents,
		DefuseEvents:       parseResult.DefuseEvents,
		RoundStats:         parseResult.RoundStats,
		PlayerRatings:      playerRatings,
		DailyKills:         dailyKills,
		DailyFlash:         dailyFlash,
		DailyDefuse:        dailyDefuse,
		DailyRounds:        dailyRounds,
	}
}

// buildKillMatrix создает матрицу убийств
func (p *Processor) buildKillMatrix(events []logparser.KillEvent, players []Player, playerIndex map[string]int, groupBy string) KillMatrix {
	matrix := make([][]int, len(players))
	for i := range matrix {
		matrix[i] = make([]int, len(players))
	}

	maxKills := 0
	for _, event := range events {
		kKey, _ := logparser.KeyAndTitle(groupBy, event.KillerName, event.KillerSID)
		vKey, _ := logparser.KeyAndTitle(groupBy, event.VictimName, event.VictimSID)

		if kIdx, ok := playerIndex[kKey]; ok {
			if vIdx, ok := playerIndex[vKey]; ok {
				matrix[kIdx][vIdx]++
				if matrix[kIdx][vIdx] > maxKills {
					maxKills = matrix[kIdx][vIdx]
				}
			}
		}
	}

	if maxKills == 0 {
		maxKills = 1
	}

	return KillMatrix{
		Matrix: matrix,
		Max:    maxKills,
	}
}

// buildWeaponData создает данные по оружию
func (p *Processor) buildWeaponData(events []logparser.KillEvent, players []Player, weapons []string, playerIndex, weaponIndex map[string]int, groupBy string) WeaponData {
	killerWeaponMatrix := make([][]int, len(players))
	victimWeaponMatrix := make([][]int, len(players))
	for i := range killerWeaponMatrix {
		killerWeaponMatrix[i] = make([]int, len(weapons))
		victimWeaponMatrix[i] = make([]int, len(weapons))
	}

	killerMax, victimMax := 0, 0

	for _, event := range events {
		if event.Weapon == "" {
			continue
		}

		kKey, _ := logparser.KeyAndTitle(groupBy, event.KillerName, event.KillerSID)
		vKey, _ := logparser.KeyAndTitle(groupBy, event.VictimName, event.VictimSID)

		if kIdx, ok := playerIndex[kKey]; ok {
			if wIdx, ok := weaponIndex[event.Weapon]; ok {
				killerWeaponMatrix[kIdx][wIdx]++
				if killerWeaponMatrix[kIdx][wIdx] > killerMax {
					killerMax = killerWeaponMatrix[kIdx][wIdx]
				}
			}
		}

		if vIdx, ok := playerIndex[vKey]; ok {
			if wIdx, ok := weaponIndex[event.Weapon]; ok {
				victimWeaponMatrix[vIdx][wIdx]++
				if victimWeaponMatrix[vIdx][wIdx] > victimMax {
					victimMax = victimWeaponMatrix[vIdx][wIdx]
				}
			}
		}
	}

	if killerMax == 0 {
		killerMax = 1
	}
	if victimMax == 0 {
		victimMax = 1
	}

	// Создаем транспонированную матрицу для "Кто с чего убивает" (Weapons × Players)
	weaponKillsMatrix := transposeMatrix(killerWeaponMatrix)

	return WeaponData{
		KillerWeaponMatrix: killerWeaponMatrix,
		VictimWeaponMatrix: victimWeaponMatrix,
		WeaponKillsMatrix:  weaponKillsMatrix,
		KillerMax:          killerMax,
		VictimMax:          victimMax,
	}
}

// buildFlashData создает данные по флешкам
func (p *Processor) buildFlashData(events []logparser.FlashEvent, players []Player, playerIndex map[string]int, groupBy string) FlashData {
	countMatrix := make([][]int, len(players))
	secondsMatrix := make([][]float64, len(players))
	for i := range countMatrix {
		countMatrix[i] = make([]int, len(players))
		secondsMatrix[i] = make([]float64, len(players))
	}

	countMax := 0
	var secondsMax float64

	for _, event := range events {
		fKey, _ := logparser.KeyAndTitle(groupBy, event.FlasherName, event.FlasherSID)
		vKey, _ := logparser.KeyAndTitle(groupBy, event.VictimName, event.VictimSID)

		if fIdx, ok := playerIndex[fKey]; ok {
			if vIdx, ok := playerIndex[vKey]; ok {
				countMatrix[fIdx][vIdx]++
				secondsMatrix[fIdx][vIdx] += event.Duration

				if countMatrix[fIdx][vIdx] > countMax {
					countMax = countMatrix[fIdx][vIdx]
				}
				if secondsMatrix[fIdx][vIdx] > secondsMax {
					secondsMax = secondsMatrix[fIdx][vIdx]
				}
			}
		}
	}

	if countMax == 0 {
		countMax = 1
	}
	if secondsMax == 0 {
		secondsMax = 1
	}

	return FlashData{
		CountMatrix:   countMatrix,
		SecondsMatrix: secondsMatrix,
		CountMax:      countMax,
		SecondsMax:    secondsMax,
	}
}

// buildDefuseData создает данные по дефьюзу
func (p *Processor) buildDefuseData(events []logparser.DefuseEvent, players []Player, playerIndex map[string]int, groupBy string) DefuseData {
	attempts := make([]int, len(players))
	withKit := make([]int, len(players))
	withoutKit := make([]int, len(players))
	successWithKit := make([]int, len(players))
	successWithoutKit := make([]int, len(players))
	abandoned := make([]int, len(players))
	failed := make([]int, len(players))

	// Отслеживаем активные дефьюзы по игрокам
	activeDefuses := make(map[string]logparser.DefuseEvent)

	totalMax := 0

	for _, event := range events {
		switch event.EventType {
		case "begin":
			pKey, _ := logparser.KeyAndTitle(groupBy, event.PlayerName, event.PlayerSID)
			if pIdx, ok := playerIndex[pKey]; ok {
				attempts[pIdx]++
				if event.WithKit {
					withKit[pIdx]++
				} else {
					withoutKit[pIdx]++
				}

				// Сохраняем активный дефьюз
				activeDefuses[pKey] = event

				if attempts[pIdx] > totalMax {
					totalMax = attempts[pIdx]
				}
			}

		case "success":
			// Находим последний активный дефьюз
			for playerKey, defuseEvent := range activeDefuses {
				if pIdx, ok := playerIndex[playerKey]; ok {
					if defuseEvent.WithKit {
						successWithKit[pIdx]++
					} else {
						successWithoutKit[pIdx]++
					}
				}
				delete(activeDefuses, playerKey) // Удаляем после успешного дефьюза
				break                            // Обычно успешен только один дефьюз
			}

		case "abandoned":
			pKey, _ := logparser.KeyAndTitle(groupBy, event.PlayerName, event.PlayerSID)
			if pIdx, ok := playerIndex[pKey]; ok {
				abandoned[pIdx]++
			}
			delete(activeDefuses, pKey)

		case "failed":
			// Взрыв бомбы - все активные дефьюзы считаются неудачными
			for playerKey := range activeDefuses {
				if pIdx, ok := playerIndex[playerKey]; ok {
					failed[pIdx]++
				}
				delete(activeDefuses, playerKey)
			}
		}
	}

	if totalMax == 0 {
		totalMax = 1
	}

	return DefuseData{
		Attempts:          attempts,
		WithKit:           withKit,
		WithoutKit:        withoutKit,
		SuccessWithKit:    successWithKit,
		SuccessWithoutKit: successWithoutKit,
		Abandoned:         abandoned,
		Failed:            failed,
		TotalMax:          totalMax,
	}
}

// transposeMatrix транспонирует матрицу
func transposeMatrix(matrix [][]int) [][]int {
	if len(matrix) == 0 {
		return [][]int{}
	}
	rows, cols := len(matrix), len(matrix[0])
	transposed := make([][]int, cols)
	for j := 0; j < cols; j++ {
		transposed[j] = make([]int, rows)
		for i := 0; i < rows; i++ {
			transposed[j][i] = matrix[i][j]
		}
	}
	return transposed
}

// buildPlayerRatings строит агрегированные рейтинги игроков
func (p *Processor) buildPlayerRatings(roundStats []logparser.RoundStats, killEvents []logparser.KillEvent, flashEvents []logparser.FlashEvent, defuseEvents []logparser.DefuseEvent) []PlayerRating {
	// Константы для байесовского рейтинга
	const K = 100.0 // Минимальное количество раундов для "достоверности"

	// Создаем маппинг AccountID -> имя игрока
	playerNames := make(map[int64]string)

	// Собираем имена из событий убийств
	for _, event := range killEvents {
		if len(event.KillerSID) > 6 {
			var accountID int64
			if n, _ := fmt.Sscanf(event.KillerSID[5:len(event.KillerSID)-1], "%d", &accountID); n > 0 && accountID > 0 {
				playerNames[accountID] = event.KillerName
			}
		}
		if len(event.VictimSID) > 6 {
			var accountID int64
			if n, _ := fmt.Sscanf(event.VictimSID[5:len(event.VictimSID)-1], "%d", &accountID); n > 0 && accountID > 0 {
				playerNames[accountID] = event.VictimName
			}
		}
	}

	// Собираем имена из событий флешек
	for _, event := range flashEvents {
		if len(event.FlasherSID) > 6 {
			var accountID int64
			if n, _ := fmt.Sscanf(event.FlasherSID[5:len(event.FlasherSID)-1], "%d", &accountID); n > 0 && accountID > 0 && playerNames[accountID] == "" {
				playerNames[accountID] = event.FlasherName
			}
		}
		if len(event.VictimSID) > 6 {
			var accountID int64
			if n, _ := fmt.Sscanf(event.VictimSID[5:len(event.VictimSID)-1], "%d", &accountID); n > 0 && accountID > 0 && playerNames[accountID] == "" {
				playerNames[accountID] = event.VictimName
			}
		}
	}

	// Собираем имена из событий дефьюза
	for _, event := range defuseEvents {
		if len(event.PlayerSID) > 6 {
			var accountID int64
			if n, _ := fmt.Sscanf(event.PlayerSID[5:len(event.PlayerSID)-1], "%d", &accountID); n > 0 && accountID > 0 && playerNames[accountID] == "" {
				playerNames[accountID] = event.PlayerName
			}
		}
	}

	// Агрегируем данные по игрокам
	playerData := make(map[int64]*PlayerRating)

	for _, round := range roundStats {
		for _, playerStats := range round.Players {
			if playerStats.AccountID == 0 {
				continue
			}

			if playerData[playerStats.AccountID] == nil {
				playerData[playerStats.AccountID] = &PlayerRating{
					AccountID: playerStats.AccountID,
					Name:      playerNames[playerStats.AccountID],
				}
			}

			rating := playerData[playerStats.AccountID]
			rating.RoundsPlayed++
			rating.TotalEPI += playerStats.Rating
			rating.TotalDamage += playerStats.Damage
			rating.TotalKills += playerStats.Kills
			rating.TotalDeaths += playerStats.Deaths
			rating.TotalAssists += playerStats.Assists

			// Проверяем победу по полю Winner
			// Winner: 2=T, 3=CT, 0=неизвестно/ничья
			if (playerStats.Team == 3 && round.Winner == 3) || (playerStats.Team == 2 && round.Winner == 2) {
				rating.WinRounds++
			}

			// Обновляем последнюю дату игры
			if round.Date != "" && (rating.LastPlayed == "" || round.Date > rating.LastPlayed) {
				rating.LastPlayed = round.Date
			}
		}
	}

	// Рассчитываем среднее EPI по всем игрокам для использования как μ (мю)
	var totalEPI float64
	var totalRounds int
	for _, rating := range playerData {
		totalEPI += rating.TotalEPI
		totalRounds += rating.RoundsPlayed
	}

	mu := 0.6 // Дефолтное значение если нет данных
	if totalRounds > 0 {
		mu = totalEPI / float64(totalRounds)
	}

	// Рассчитываем финальные рейтинги
	ratings := make([]PlayerRating, 0, len(playerData))
	for _, rating := range playerData {
		// Простое среднее
		if rating.RoundsPlayed > 0 {
			rating.AverageEPI = rating.TotalEPI / float64(rating.RoundsPlayed)
		}

		// Байесовский рейтинг с регуляризацией
		// BayesianEPI = (TotalEPI + K * μ) / (RoundsPlayed + K)
		rating.BayesianEPI = (rating.TotalEPI + K*mu) / (float64(rating.RoundsPlayed) + K)

		// Если нет имени, используем AccountID
		if rating.Name == "" {
			rating.Name = fmt.Sprintf("Player_%d", rating.AccountID)
		}

		ratings = append(ratings, *rating)
	}

	// Сортируем по байесовскому рейтингу (по убыванию)
	sort.Slice(ratings, func(i, j int) bool {
		return ratings[i].BayesianEPI > ratings[j].BayesianEPI
	})

	return ratings
}
