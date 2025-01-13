package teambuilder

import (
	"math"
	"sort"
)

/*
Package teambuilder предоставляет функционал для автоматического распределения игроков по командам
с учетом их индивидуальных характеристик и различных ограничений.

Основные возможности:
  - Справедливое распределение игроков по командам
  - Поддержка ограничений на совместное/раздельное размещение игроков
  - Оптимизация баланса команд по суммарному рейтингу
  - Несколько алгоритмов распределения (змейка, пары, жадный)

Примеры использования:
   builder := &TeamBuilder{}
   config := &TeamConfiguration{
       Players: []TeamPlayer{
           {"Player1", 1000},
           {"Player2", 2000},
           // ...
       },
       Constraints: []Constraint{
           {Type: ConstraintTogether, Player1: "Player1", Player2: "Player2"},
       },
   }
   team1, team2 := builder.Build(config)
*/

// TeamBuilder реализует построитель команд с поддержкой сложных ограничений.
// Основная задача - справедливое распределение игроков с учетом их индивидуальных характеристик.
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

	// Сначала обработаем ограничения типа ConstraintTogether
	linkedPlayers := make(map[string]string)
	for _, constraint := range constraints {
		if constraint.Type == ConstraintTogether {
			linkedPlayers[constraint.Player1] = constraint.Player2
			linkedPlayers[constraint.Player2] = constraint.Player1
		}
	}

	// Сортируем игроков по убыванию веса
	sort.Slice(players, func(i, j int) bool {
		return players[i].Score > players[j].Score
	})

	// Пробуем все три метода распределения и выбираем лучший результат
	var bestTeam1, bestTeam2 Team
	bestDiff := math.Inf(1)

	// Метод 1: Начальное распределение с учетом связанных игроков
	team1, team2 := distributeWithLinkedPlayers(players, linkedPlayers)
	if isConstraintSatisfied(team1, team2, constraints) {
		diff := math.Abs(getTeamScore(team1) - getTeamScore(team2))
		if diff < bestDiff {
			bestDiff = diff
			bestTeam1, bestTeam2 = team1, team2
		}
	}

	// Метод 2: Распределение змейкой
	team1, team2 = distributeSnake(players)
	if isConstraintSatisfied(team1, team2, constraints) {
		diff := math.Abs(getTeamScore(team1) - getTeamScore(team2))
		if diff < bestDiff {
			bestDiff = diff
			bestTeam1, bestTeam2 = team1, team2
		}
	}

	// Метод 3: Распределение парами
	team1, team2 = distributePairs(players)
	if isConstraintSatisfied(team1, team2, constraints) {
		diff := math.Abs(getTeamScore(team1) - getTeamScore(team2))
		if diff < bestDiff {
			bestDiff = diff
			bestTeam1, bestTeam2 = team1, team2
		}
	}

	// Метод 4: Жадное распределение
	team1, team2 = distributeGreedy(players)
	if isConstraintSatisfied(team1, team2, constraints) {
		diff := math.Abs(getTeamScore(team1) - getTeamScore(team2))
		if diff < bestDiff {
			bestDiff = diff
			bestTeam1, bestTeam2 = team1, team2
		}
	}

	// Если нашли хотя бы одно валидное решение, оптимизируем его
	if bestDiff != math.Inf(1) {
		bestTeam1, bestTeam2 = optimizeTeams(bestTeam1, bestTeam2, constraints)
		return bestTeam1, bestTeam2
	}

	// Если не нашли валидного решения, возвращаем результат жадного алгоритма
	// и пытаемся его оптимизировать
	team1, team2 = distributeGreedy(players)
	return optimizeTeams(team1, team2, constraints)
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
func isConstraintSatisfied(team1, team2 Team, constraints []Constraint) bool {
	for _, constraint := range constraints {
		player1InTeam1 := playerInTeam(team1, constraint.Player1)
		player1InTeam2 := playerInTeam(team2, constraint.Player1)
		player2InTeam1 := playerInTeam(team1, constraint.Player2)
		player2InTeam2 := playerInTeam(team2, constraint.Player2)

		switch constraint.Type {
		case ConstraintTogether:
			if (player1InTeam1 && !player2InTeam1) ||
				(player1InTeam2 && !player2InTeam2) {
				return false
			}
		case ConstraintSeparate:
			if (player1InTeam1 && player2InTeam1) ||
				(player1InTeam2 && player2InTeam2) {
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
func distributeWithLinkedPlayers(players []TeamPlayer, linkedPlayers map[string]string) (Team, Team) {
	teamSize := len(players) / 2
	if len(players)%2 != 0 {
		teamSize++
	}

	team1 := make(Team, 0, teamSize)
	team2 := make(Team, 0, teamSize)
	used := make(map[string]bool)

	// Сначала распределяем связанных игроков
	for i := 0; i < len(players); i++ {
		if used[players[i].NickName] {
			continue
		}

		if linked, ok := linkedPlayers[players[i].NickName]; ok {
			var linkedPlayer TeamPlayer
			for _, p := range players {
				if p.NickName == linked {
					linkedPlayer = p
					break
				}
			}

			if getTeamScore(team1) <= getTeamScore(team2) {
				team1 = append(team1, players[i], linkedPlayer)
			} else {
				team2 = append(team2, players[i], linkedPlayer)
			}

			used[players[i].NickName] = true
			used[linked] = true
		}
	}

	// Распределяем оставшихся игроков
	for i := 0; i < len(players); i++ {
		if used[players[i].NickName] {
			continue
		}

		if getTeamScore(team1) <= getTeamScore(team2) && len(team1) < teamSize {
			team1 = append(team1, players[i])
		} else {
			team2 = append(team2, players[i])
		}
		used[players[i].NickName] = true
	}

	return team1, team2
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
		if getTeamScore(team1) <= getTeamScore(team2) && len(team1) < teamSize {
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
func optimizeTeams(team1, team2 Team, constraints []Constraint) (Team, Team) {
	bestTeam1, bestTeam2 := team1, team2
	bestDiff := math.Abs(getTeamScore(team1) - getTeamScore(team2))

	for attempt := 0; attempt < 3; attempt++ {
		improved := false

		for i := 0; i < len(team1); i++ {
			for j := 0; j < len(team2); j++ {
				newTeam1 := make(Team, len(team1))
				newTeam2 := make(Team, len(team2))
				copy(newTeam1, bestTeam1)
				copy(newTeam2, bestTeam2)

				newTeam1[i], newTeam2[j] = newTeam2[j], newTeam1[i]

				if !isConstraintSatisfied(newTeam1, newTeam2, constraints) {
					continue
				}

				newDiff := math.Abs(getTeamScore(newTeam1) - getTeamScore(newTeam2))
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
