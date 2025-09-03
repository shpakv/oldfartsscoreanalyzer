/*
Package teambuilder — ваш личный инструктор по созданию честных матчей в CS2.
Если в ваших играх слишком часто звучит фраза «Ну это было нечестно», то вы по адресу.

Основные возможности:
  - **Справедливое распределение игроков**: Мы стараемся сделать ваши команды честными, но если вы всё равно проигрываете — это не наша вина. 🤷‍♂️
  - **Поддержка ограничений**: Хотите, чтобы два друга не играли вместе (или наоборот)? Легко.
  - **Баланс по рейтингу**: Мы учитываем скиллы игроков, чтобы дать всем равные шансы на победу.
  - **Несколько алгоритмов**: Змейка, пары, жадный метод — всё, чтобы ваш матч прошёл на высшем уровне.

Пример использования:

	builder := &TeamBuilder{}
	config := &TeamConfiguration{
		Players: []TeamPlayer{
			{"Player1", 1000},
			{"Player2", 2000},
		},
		Constraints: []Constraint{
			{Type: ConstraintTogether, Player1: "Player1", Player2: "Player2"},
		},
	}
	team1, team2 := builder.Build(config)
*/
package teambuilder

import (
	"log"
	"math"
	"sort"
	"time"
)

// TeamBuilder — это тот самый алгоритмический гений, который берёт ваш список игроков
// и создаёт команды, настолько честные, насколько это возможно в CS2.
type TeamBuilder struct{}

// Build выполняет распределение игроков по командам с учетом заданных ограничений.
// Использует комбинацию нескольких алгоритмов распределения для достижения оптимального результата.
//
// Параметры:
//   - config *TeamConfiguration: конфигурация, содержащая список игроков и ограничений
//
// Возвращает:
//   - (Team, Team): две сбалансированные команды, удовлетворяющие заданным ограничениям
//
// Алгоритм работы:
//  1. Обработка ограничений типа ConstraintTogether
//  2. Сортировка игроков по убыванию рейтинга
//  3. Распределение связанных игроков
//  4. Распределение оставшихся игроков
//  5. Оптимизация полученного результата
//
// Сложность: O(2^n), где n - количество игроков
func (b *TeamBuilder) Build(config *TeamConfiguration) (Team, Team) {
	players := config.Players
	constraints := config.Constraints
	economicConfig := config.EconomicConfig

	// Проверка на пустой список игроков
	if len(players) == 0 {
		return Team{}, Team{}
	}

	// Если экономические настройки не заданы, используем дефолтные для равных команд
	if !economicConfig.Enabled && economicConfig.BasePercentage == 0 {
		economicConfig = EconomicConfig{
			Enabled:        false, // Отключено для равных команд по умолчанию
			BasePercentage: 15.0,
			MaxPercentage:  5.0,
			MinPercentage:  1.0,
		}
	}

	// Сначала обработаем ограничения типа ConstraintTogether
	linkedPlayers := make(map[string]string)
	for _, constraint := range constraints {
		if constraint.Type == ConstraintTogether {
			linkedPlayers[constraint.Player1] = constraint.Player2
			linkedPlayers[constraint.Player2] = constraint.Player1
		}
	}

	// Логирование начала процесса
	startTime := time.Now()
	log.Printf("🚀 TeamBuilder: Начинаем распределение %d игроков с %d ограничениями",
		len(players), len(constraints))

	if len(linkedPlayers) > 0 {
		log.Printf("🔗 Найдены связанные игроки: %d пар", len(linkedPlayers)/2)
	}

	if economicConfig.Enabled {
		log.Printf("💰 Экономические преимущества включены: Base=%.1f%%, Max=%.1f%%, Min=%.1f%%",
			economicConfig.BasePercentage, economicConfig.MaxPercentage, economicConfig.MinPercentage)
	}

	// Сортируем игроков по убыванию веса
	sort.Slice(players, func(i, j int) bool {
		return players[i].Score > players[j].Score
	})

	// Параллельно выполняем все методы распределения и выбираем лучший результат
	type algorithmResult struct {
		team1        Team
		team2        Team
		diff         float64
		valid        bool
		methodName   string
		duration     time.Duration
		team1Score   float64
		team2Score   float64
		economicDiff float64
	}

	// Канал для получения результатов от горутин
	results := make(chan algorithmResult, 4)

	// Запускаем все алгоритмы параллельно с детальным логированием
	go func() {
		start := time.Now()
		team1, team2 := distributeWithLinkedPlayers(players, linkedPlayers)
		duration := time.Since(start)
		valid := isConstraintSatisfied(team1, team2, constraints)
		diff := calculateTeamBalanceDifference(team1, team2, economicConfig)
		team1Score := team1.Score()
		team2Score := team2.Score()

		// Для неравных команд показываем экономическую разницу
		economicDiff := diff
		if len(team1) != len(team2) && economicConfig.Enabled {
			economicDiff = math.Abs(team1Score - team2Score)
		}

		results <- algorithmResult{
			team1, team2, diff, valid, "LinkedPlayers", duration,
			team1Score, team2Score, economicDiff,
		}
	}()

	go func() {
		start := time.Now()
		team1, team2 := distributeSnake(players)
		duration := time.Since(start)
		valid := isConstraintSatisfied(team1, team2, constraints)
		diff := calculateTeamBalanceDifference(team1, team2, economicConfig)
		team1Score := team1.Score()
		team2Score := team2.Score()

		economicDiff := diff
		if len(team1) != len(team2) && economicConfig.Enabled {
			economicDiff = math.Abs(team1Score - team2Score)
		}

		results <- algorithmResult{
			team1, team2, diff, valid, "Snake", duration,
			team1Score, team2Score, economicDiff,
		}
	}()

	go func() {
		start := time.Now()
		team1, team2 := distributePairs(players)
		duration := time.Since(start)
		valid := isConstraintSatisfied(team1, team2, constraints)
		diff := calculateTeamBalanceDifference(team1, team2, economicConfig)
		team1Score := team1.Score()
		team2Score := team2.Score()

		economicDiff := diff
		if len(team1) != len(team2) && economicConfig.Enabled {
			economicDiff = math.Abs(team1Score - team2Score)
		}

		results <- algorithmResult{
			team1, team2, diff, valid, "Pairs", duration,
			team1Score, team2Score, economicDiff,
		}
	}()

	go func() {
		start := time.Now()
		team1, team2 := distributeGreedyWithConfig(players, economicConfig)
		duration := time.Since(start)
		valid := isConstraintSatisfied(team1, team2, constraints)
		diff := calculateTeamBalanceDifference(team1, team2, economicConfig)
		team1Score := team1.Score()
		team2Score := team2.Score()

		economicDiff := diff
		if len(team1) != len(team2) && economicConfig.Enabled {
			economicDiff = math.Abs(team1Score - team2Score)
		}

		results <- algorithmResult{
			team1, team2, diff, valid, "GreedyEcon", duration,
			team1Score, team2Score, economicDiff,
		}
	}()

	// Собираем результаты от всех горутин с подробным логированием
	var bestTeam1, bestTeam2 Team
	bestDiff := math.Inf(1)
	var bestMethod string
	var algorithmResults []algorithmResult

	log.Printf("⚡ Алгоритмы запущены параллельно, ждем результаты...")

	for i := 0; i < 4; i++ {
		result := <-results
		algorithmResults = append(algorithmResults, result)

		// Логирование результата каждого алгоритма
		status := "❌ INVALID"
		if result.valid {
			status = "✅ VALID"
		}

		// Показываем экономическое влияние для неравных команд
		if len(result.team1) != len(result.team2) && economicConfig.Enabled {
			log.Printf("📊 %s: %s | Время: %v | Баланс: %.1f vs %.1f | До эконом.: %.1f | После: %.1f",
				result.methodName, status, result.duration,
				result.team1Score, result.team2Score, result.economicDiff, result.diff)
		} else {
			log.Printf("📊 %s: %s | Время: %v | Баланс: %.1f vs %.1f (разница: %.1f)",
				result.methodName, status, result.duration,
				result.team1Score, result.team2Score, result.diff)
		}

		if result.valid && result.diff < bestDiff {
			bestDiff = result.diff
			bestMethod = result.methodName
			bestTeam1 = make(Team, len(result.team1))
			bestTeam2 = make(Team, len(result.team2))
			copy(bestTeam1, result.team1)
			copy(bestTeam2, result.team2)
		}
	}

	// Сводная статистика по алгоритмам
	validAlgorithms := 0
	totalDuration := time.Duration(0)
	for _, result := range algorithmResults {
		if result.valid {
			validAlgorithms++
		}
		totalDuration += result.duration
	}

	log.Printf("📈 Сводка: %d/%d алгоритмов дали валидные результаты за %v",
		validAlgorithms, len(algorithmResults), totalDuration)

	// Если нашли хотя бы одно валидное решение, оптимизируем его
	if bestDiff != math.Inf(1) {
		log.Printf("🎯 Лучший алгоритм: %s с разницей %.1f", bestMethod, bestDiff)
		log.Printf("⚙️  Запускаем оптимизацию...")

		optimizationStart := time.Now()
		optimizedTeam1, optimizedTeam2 := optimizeTeamsWithConfig(bestTeam1, bestTeam2, constraints, economicConfig)
		optimizationTime := time.Since(optimizationStart)

		finalDiff := calculateTeamBalanceDifference(optimizedTeam1, optimizedTeam2, economicConfig)
		improvement := bestDiff - finalDiff

		if improvement > 0.1 {
			log.Printf("✨ Оптимизация улучшила баланс на %.1f за %v (%.1f -> %.1f)",
				improvement, optimizationTime, bestDiff, finalDiff)
		} else {
			log.Printf("🔧 Оптимизация завершена за %v (баланс: %.1f)",
				optimizationTime, finalDiff)
		}

		totalTime := time.Since(startTime)
		log.Printf("🏁 Команды сформированы за %v! Итоговый баланс: %.1f",
			totalTime, finalDiff)

		return optimizedTeam1, optimizedTeam2
	}

	// Если не нашли валидного решения, используем fallback
	log.Printf("⚠️  Ни один алгоритм не удовлетворил ограничения, используем fallback")
	fallbackStart := time.Now()
	fallbackTeam1, fallbackTeam2 := distributeGreedyWithConfig(players, economicConfig)
	fallbackTime := time.Since(fallbackStart)

	optimizedTeam1, optimizedTeam2 := optimizeTeamsWithConfig(fallbackTeam1, fallbackTeam2, constraints, economicConfig)
	totalTime := time.Since(startTime)
	finalDiff := calculateTeamBalanceDifference(optimizedTeam1, optimizedTeam2, economicConfig)

	log.Printf("🔄 Fallback выполнен за %v, итого %v | Финальный баланс: %.1f",
		fallbackTime, totalTime, finalDiff)

	return optimizedTeam1, optimizedTeam2
}

// calculateTeamBalanceDifference вычисляет разницу в балансе команд
func calculateTeamBalanceDifference(team1, team2 Team, config EconomicConfig) float64 {
	if len(team1) == len(team2) {
		// Для равных команд используем обычный баланс
		return math.Abs(team1.Score() - team2.Score())
	}
	// Для неравных команд используем эффективный баланс с экономикой
	score1, score2 := GetEffectiveTeamScoreWithConfig(team1, team2, config)
	return math.Abs(score1 - score2)
}

// getTeamScore вычисляет суммарный рейтинг команды.
//
// Параметры:
//   - team Team: команда для подсчета рейтинга
//
// Возвращает:
//   - float64: суммарный рейтинг всех игроков команды
func getTeamScore(team Team) float64 {
	score := 0.0
	for _, player := range team {
		score += player.Score
	}
	return score
}

// isConstraintSatisfied проверяет соответствие распределения игроков заданным ограничениям.
//
// Параметры:
//   - team1, team2 Team: проверяемые команды
//   - constraints []Constraint: список ограничений
//
// Возвращает:
//   - bool: true если все ограничения соблюдены, false в противном случае
func isConstraintSatisfied(team1, team2 Team, constraints Constraints) bool {
	for _, constraint := range constraints {
		player1InTeam1 := playerInTeam(team1, constraint.Player1)
		player1InTeam2 := playerInTeam(team2, constraint.Player1)
		player2InTeam1 := playerInTeam(team1, constraint.Player2)
		player2InTeam2 := playerInTeam(team2, constraint.Player2)

		// Проверяем, что оба игрока существуют
		if (!player1InTeam1 && !player1InTeam2) || (!player2InTeam1 && !player2InTeam2) {
			// Пропускаем ограничение, если один из игроков не существует
			continue
		}

		switch constraint.Type {
		case ConstraintTogether:
			if (player1InTeam1 && player2InTeam2) || (player1InTeam2 && player2InTeam1) {
				return false
			}
		case ConstraintSeparate:
			if (player1InTeam1 && player2InTeam1) || (player1InTeam2 && player2InTeam2) {
				return false
			}
		}
	}
	return true
}

// distributeSnake реализует распределение игроков методом "змейки".
// Распределяет игроков поочередно между командами, учитывая текущий баланс.
//
// Параметры:
//   - players []TeamPlayer: отсортированный список игроков
//
// Возвращает:
//   - (Team, Team): распределенные команды
func distributeSnake(players []TeamPlayer) (Team, Team) {
	teamSize := len(players) / 2
	team1 := make(Team, 0, teamSize)
	team2 := make(Team, 0, teamSize)

	for i := 0; i < len(players); i += 2 {
		if i+1 < len(players) {
			if getTeamScore(team1) <= getTeamScore(team2) {
				team1 = append(team1, players[i])
				team2 = append(team2, players[i+1])
			} else {
				team2 = append(team2, players[i])
				team1 = append(team1, players[i+1])
			}
		} else {
			if getTeamScore(team1) <= getTeamScore(team2) {
				team1 = append(team1, players[i])
			} else {
				team2 = append(team2, players[i])
			}
		}
	}
	return team1, team2
}

// Вспомогательная функция для распределения с учетом связанных игроков
func distributeWithLinkedPlayers(players Team, linkedPlayers map[string]string) (Team, Team) {
	teamSize := len(players) / 2
	if len(players)%2 != 0 {
		teamSize++
	}

	team1 := make(Team, 0, teamSize)
	team2 := make(Team, 0, teamSize)
	used := make(map[string]bool)

	// Создаем карту для быстрого поиска игроков по имени
	playerMap := make(map[string]TeamPlayer)
	for _, p := range players {
		playerMap[p.NickName] = p
	}

	// Находим группы связанных игроков (компоненты связности)
	groups := findConnectedGroups(players, linkedPlayers)

	// Распределяем группы связанных игроков
	for _, group := range groups {
		// Собираем игроков из группы
		groupPlayers := make([]TeamPlayer, 0, len(group))
		for _, name := range group {
			if player, ok := playerMap[name]; ok {
				groupPlayers = append(groupPlayers, player)
				used[name] = true
			}
		}

		// Если группа пустая, пропускаем
		if len(groupPlayers) == 0 {
			continue
		}

		// Распределяем всю группу в одну команду
		if team1.Score() <= team2.Score() && len(team1)+len(groupPlayers) <= teamSize {
			team1 = append(team1, groupPlayers...)
		} else {
			team2 = append(team2, groupPlayers...)
		}
	}

	// Распределяем оставшихся игроков
	for i := 0; i < len(players); i++ {
		if used[players[i].NickName] {
			continue
		}

		if team1.Score() <= team2.Score() && len(team1) < teamSize {
			team1 = append(team1, players[i])
		} else {
			team2 = append(team2, players[i])
		}
		used[players[i].NickName] = true
	}

	return team1, team2
}

// findConnectedGroups находит группы связанных игроков (компоненты связности в графе)
func findConnectedGroups(players Team, linkedPlayers map[string]string) [][]string {
	// Создаем граф связей между игроками
	graph := make(map[string][]string)
	for _, player := range players {
		graph[player.NickName] = []string{}
	}

	// Заполняем граф связями
	for player, linked := range linkedPlayers {
		// Проверяем, что оба игрока существуют
		if _, ok := graph[player]; ok {
			if _, ok := graph[linked]; ok {
				graph[player] = append(graph[player], linked)
				graph[linked] = append(graph[linked], player)
			}
		}
	}

	// Находим компоненты связности с помощью поиска в глубину
	visited := make(map[string]bool)
	var groups [][]string

	for player := range graph {
		if !visited[player] {
			group := []string{}
			dfs(player, graph, visited, &group)
			if len(group) > 0 {
				groups = append(groups, group)
			}
		}
	}

	return groups
}

// dfs выполняет поиск в глубину для нахождения всех связанных игроков
func dfs(player string, graph map[string][]string, visited map[string]bool, group *[]string) {
	visited[player] = true
	*group = append(*group, player)

	for _, neighbor := range graph[player] {
		if !visited[neighbor] {
			dfs(neighbor, graph, visited, group)
		}
	}
}

// distributePairs реализует распределение игроков парами.
// Формирует пары из сильнейшего и слабейшего игроков.
//
// Параметры:
//   - players []TeamPlayer: отсортированный список игроков
//
// Возвращает:
//   - (Team, Team): распределенные команды
func distributePairs(players []TeamPlayer) (Team, Team) {
	teamSize := len(players) / 2
	team1 := make(Team, 0, teamSize)
	team2 := make(Team, 0, teamSize)

	for i := 0; i < len(players)/2; i++ {
		if getTeamScore(team1) <= getTeamScore(team2) {
			team1 = append(team1, players[i])
			team2 = append(team2, players[len(players)-1-i])
		} else {
			team2 = append(team2, players[i])
			team1 = append(team1, players[len(players)-1-i])
		}
	}

	// Если осталось нечетное количество игроков
	if len(players)%2 != 0 {
		if getTeamScore(team1) <= getTeamScore(team2) {
			team1 = append(team1, players[len(players)/2])
		} else {
			team2 = append(team2, players[len(players)/2])
		}
	}

	return team1, team2
}

// distributeGreedy реализует жадный алгоритм распределения игроков.
// Добавляет каждого следующего игрока в команду с меньшим суммарным рейтингом.
//
// Параметры:
//   - players []TeamPlayer: список игроков
//
// Возвращает:
//   - (Team, Team): распределенные команды
func distributeGreedy(players []TeamPlayer) (Team, Team) {
	teamSize := len(players) / 2
	if len(players)%2 != 0 {
		teamSize++
	}
	team1 := make(Team, 0, teamSize)
	team2 := make(Team, 0, teamSize)

	for _, player := range players {
		// Calculate projected team scores including economic advantage
		projTeam1 := append(Team{}, team1...)
		projTeam1 = append(projTeam1, player)

		projTeam2 := append(Team{}, team2...)
		projTeam2 = append(projTeam2, player)

		// Calculate effective scores for both scenarios
		score1Team1, score1Team2 := getEffectiveTeamScore(projTeam1, team2)
		score2Team1, score2Team2 := getEffectiveTeamScore(team1, projTeam2)

		// Choose the option that minimizes score difference
		diff1 := math.Abs(score1Team1 - score1Team2)
		diff2 := math.Abs(score2Team1 - score2Team2)

		if (diff1 <= diff2 && len(team1) < teamSize) || len(team2) >= teamSize {
			team1 = append(team1, player)
		} else {
			team2 = append(team2, player)
		}
	}
	return team1, team2
}

// optimizeTeams выполняет оптимизацию начального распределения игроков
// путем попарного обмена игроков между командами.
//
// Параметры:
//   - team1, team2 Team: исходные команды
//   - constraints []Constraint: список ограничений
//
// Возвращает:
//   - (Team, Team): оптимизированные команды
//
// Особенности:
//   - Выполняет до 3 попыток оптимизации
//   - Прекращает оптимизацию, если улучшение не достигнуто
//   - Сохраняет все ограничения при обмене игроками
func optimizeTeams(team1, team2 Team, constraints Constraints) (Team, Team) {
	bestTeam1 := make(Team, len(team1))
	bestTeam2 := make(Team, len(team2))
	copy(bestTeam1, team1)
	copy(bestTeam2, team2)

	// Calculate effective scores considering economic advantage
	team1Score, team2Score := getEffectiveTeamScore(bestTeam1, bestTeam2)
	bestDiff := math.Abs(team1Score - team2Score)

	for attempt := 0; attempt < 3; attempt++ {
		improved := false

		for i := 0; i < len(bestTeam1); i++ {
			for j := 0; j < len(bestTeam2); j++ {
				// Create copies of current best teams
				newTeam1 := make(Team, len(bestTeam1))
				newTeam2 := make(Team, len(bestTeam2))
				copy(newTeam1, bestTeam1)
				copy(newTeam2, bestTeam2)

				// Swap players
				newTeam1[i], newTeam2[j] = newTeam2[j], newTeam1[i]

				if !isConstraintSatisfied(newTeam1, newTeam2, constraints) {
					continue
				}

				// Calculate effective scores with economic advantage
				newTeam1Score, newTeam2Score := getEffectiveTeamScore(newTeam1, newTeam2)
				newDiff := math.Abs(newTeam1Score - newTeam2Score)

				if newDiff < bestDiff {
					bestDiff = newDiff
					copy(bestTeam1, newTeam1)
					copy(bestTeam2, newTeam2)
					improved = true
				}
			}
		}

		if !improved {
			break
		}
	}

	return bestTeam1, bestTeam2
}

// GetEffectiveTeamScoreWithConfig рассчитывает эффективные счета команд с учетом экономических преимуществ
func GetEffectiveTeamScoreWithConfig(team1, team2 Team, config EconomicConfig) (float64, float64) {
	team1Score := team1.Score()
	team2Score := team2.Score()

	// Если экономические преимущества отключены, возвращаем базовые счета
	if !config.Enabled {
		return team1Score, team2Score
	}

	// Рассчитываем адаптивное экономическое преимущество
	if len(team1) < len(team2) {
		// Team1 имеет меньше игроков, добавляем экономический бонус
		percentPerPlayer := calculateAdaptivePercentage(len(team1), config)
		economicBonus := team1Score * percentPerPlayer * float64(len(team2)-len(team1))
		team1Score += economicBonus
	} else if len(team2) < len(team1) {
		// Team2 имеет меньше игроков, добавляем экономический бонус
		percentPerPlayer := calculateAdaptivePercentage(len(team2), config)
		economicBonus := team2Score * percentPerPlayer * float64(len(team1)-len(team2))
		team2Score += economicBonus
	}

	return team1Score, team2Score
}

// calculateAdaptivePercentage вычисляет адаптивный процент экономического преимущества
// Формула: BasePercentage / smallerTeamSize с ограничениями min/max
func calculateAdaptivePercentage(smallerTeamSize int, config EconomicConfig) float64 {
	if smallerTeamSize == 0 {
		return 0.0
	}

	// Основная формула: чем больше команда, тем меньше влияние экономики
	percentage := config.BasePercentage / float64(smallerTeamSize)

	// Применяем ограничения
	if percentage > config.MaxPercentage {
		percentage = config.MaxPercentage
	}
	if percentage < config.MinPercentage {
		percentage = config.MinPercentage
	}

	// Возвращаем как долю (например, 0.05 для 5%)
	return percentage / 100.0
}

// getEffectiveTeamScore обертка для обратной совместимости с дефолтными настройками
func getEffectiveTeamScore(team1, team2 Team) (float64, float64) {
	defaultConfig := EconomicConfig{
		Enabled:        true,
		BasePercentage: 15.0, // 15% базовый процент
		MaxPercentage:  5.0,  // максимум 5% за игрока
		MinPercentage:  1.0,  // минимум 1% за игрока
	}
	return GetEffectiveTeamScoreWithConfig(team1, team2, defaultConfig)
}

// distributeGreedyWithConfig реализует жадный алгоритм с учетом экономической конфигурации
func distributeGreedyWithConfig(players []TeamPlayer, config EconomicConfig) (Team, Team) {
	teamSize := len(players) / 2
	if len(players)%2 != 0 {
		teamSize++
	}
	team1 := make(Team, 0, teamSize)
	team2 := make(Team, 0, teamSize)

	for _, player := range players {
		// Calculate projected team scores including economic advantage
		projTeam1 := append(Team{}, team1...)
		projTeam1 = append(projTeam1, player)

		projTeam2 := append(Team{}, team2...)
		projTeam2 = append(projTeam2, player)

		// Calculate effective scores for both scenarios using config
		score1Team1, score1Team2 := GetEffectiveTeamScoreWithConfig(projTeam1, team2, config)
		score2Team1, score2Team2 := GetEffectiveTeamScoreWithConfig(team1, projTeam2, config)

		// Choose the option that minimizes score difference
		diff1 := math.Abs(score1Team1 - score1Team2)
		diff2 := math.Abs(score2Team1 - score2Team2)

		if (diff1 <= diff2 && len(team1) < teamSize) || len(team2) >= teamSize {
			team1 = append(team1, player)
		} else {
			team2 = append(team2, player)
		}
	}
	return team1, team2
}

// optimizeTeamsWithConfig выполняет оптимизацию с учетом экономической конфигурации
func optimizeTeamsWithConfig(team1, team2 Team, constraints Constraints, config EconomicConfig) (Team, Team) {
	bestTeam1 := make(Team, len(team1))
	bestTeam2 := make(Team, len(team2))
	copy(bestTeam1, team1)
	copy(bestTeam2, team2)

	// Calculate effective scores considering economic advantage
	team1Score, team2Score := GetEffectiveTeamScoreWithConfig(bestTeam1, bestTeam2, config)
	bestDiff := math.Abs(team1Score - team2Score)

	for attempt := 0; attempt < 3; attempt++ {
		improved := false

		for i := 0; i < len(bestTeam1); i++ {
			for j := 0; j < len(bestTeam2); j++ {
				// Create copies of current best teams
				newTeam1 := make(Team, len(bestTeam1))
				newTeam2 := make(Team, len(bestTeam2))
				copy(newTeam1, bestTeam1)
				copy(newTeam2, bestTeam2)

				// Swap players
				newTeam1[i], newTeam2[j] = newTeam2[j], newTeam1[i]

				if !isConstraintSatisfied(newTeam1, newTeam2, constraints) {
					continue
				}

				// Calculate effective scores with economic advantage
				newTeam1Score, newTeam2Score := GetEffectiveTeamScoreWithConfig(newTeam1, newTeam2, config)
				newDiff := math.Abs(newTeam1Score - newTeam2Score)

				if newDiff < bestDiff {
					bestDiff = newDiff
					copy(bestTeam1, newTeam1)
					copy(bestTeam2, newTeam2)
					improved = true
				}
			}
		}

		if !improved {
			break
		}
	}

	return bestTeam1, bestTeam2
}

// playerInTeam проверяет наличие игрока в команде.
//
// Параметры:
//   - team Team: команда для проверки
//   - playerName string: имя искомого игрока
//
// Возвращает:
//   - bool: true если игрок найден в команде, false в противном случае
func playerInTeam(team Team, playerName string) bool {
	for _, player := range team {
		if player.NickName == playerName {
			return true
		}
	}
	return false
}
