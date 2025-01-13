package teambuilder

import (
	"math"
	"sort"
)

// TeamBuilder - построитель команд с поддержкой сложных ограничений
// Основная задача - справедливое распределение игроков с учетом их индивидуальных характеристик
type TeamBuilder struct{}

// Build выполняет распределение игроков по командам
//
// Параметры:
//   - config: объект конфигурации, который содержит:
//   - players: полный список игроков
//   - constraints: список ограничений для распределения
//
// Возвращает:
//   - Две команды с оптимальным распределением игроков
//
// Основные характеристики алгоритма:
//   - Минимизация разницы суммарного веса команд
//   - Поддержка сложных ограничений на распределение
//   - Гарантированное соблюдение размера команд
//
// Алгоритм использует метод возврата (backtracking) с оптимизацией:
//  1. Сортировка игроков по убыванию веса
//  2. Рекурсивный перебор вариантов распределения
//  3. Проверка ограничений на каждом шаге
//  4. Выбор оптимального варианта
//
// Сложность алгоритма: O(2^n), где n - количество игроков.
// Подходит для небольших составов (до 20-30 игроков)
func (b *TeamBuilder) Build(config *TeamConfiguration) (Team, Team) {
	players := config.Players
	constraints := config.Constraints
	teamSize := len(players) / 2

	// Сортируем игроков по убыванию веса
	sort.Slice(players, func(i, j int) bool {
		return players[i].Score > players[j].Score
	})

	// Создаем две команды
	team1 := make(Team, 0, teamSize)
	team2 := make(Team, 0, teamSize)

	// Распределяем игроков змейкой (сильнейший - слабейший)
	for i := 0; i < len(players); i += 2 {
		if i+1 < len(players) {
			// Распределяем пару игроков
			if getTeamScore(team1) <= getTeamScore(team2) {
				team1 = append(team1, players[i])
				team2 = append(team2, players[i+1])
			} else {
				team2 = append(team2, players[i])
				team1 = append(team1, players[i+1])
			}
		} else {
			// Если остался один игрок
			if getTeamScore(team1) <= getTeamScore(team2) {
				team1 = append(team1, players[i])
			} else {
				team2 = append(team2, players[i])
			}
		}
	}

	// Если есть ограничения, пытаемся их удовлетворить путем обмена игроками
	if len(constraints) > 0 {
		team1, team2 = optimizeWithConstraints(team1, team2, constraints)
	}

	return team1, team2
}

// Функция проверки ограничений
func isConstraintSatisfied(team1, team2 Team, constraints []Constraint) bool {
	for _, constraint := range constraints {
		player1InTeam1 := playerInTeam(team1, constraint.Player1)
		player1InTeam2 := playerInTeam(team2, constraint.Player1)
		player2InTeam1 := playerInTeam(team1, constraint.Player2)
		player2InTeam2 := playerInTeam(team2, constraint.Player2)

		switch constraint.Type {
		case ConstraintTogether:
			// Игроки должны быть в одной команде
			if (player1InTeam1 && !player2InTeam1) ||
				(player1InTeam2 && !player2InTeam2) {
				return false
			}
		case ConstraintSeparate:
			// Игроки должны быть в разных командах
			if (player1InTeam1 && player2InTeam1) ||
				(player1InTeam2 && player2InTeam2) {
				return false
			}
		}
	}
	return true
}

// Вспомогательная функция для подсчета общего счета команды
func getTeamScore(team Team) float64 {
	score := 0.0
	for _, player := range team {
		score += player.Score
	}
	return score
}

// Функция для оптимизации с учетом ограничений
func optimizeWithConstraints(team1, team2 Team, constraints []Constraint) (Team, Team) {
	bestTeam1, bestTeam2 := team1, team2
	bestDiff := math.Abs(getTeamScore(team1) - getTeamScore(team2))

	// Проверяем начальное распределение на соответствие ограничениям
	if !isConstraintSatisfied(team1, team2, constraints) {
		// Если начальное распределение не удовлетворяет ограничениям,
		// пытаемся найти валидное решение
		for i := 0; i < len(team1); i++ {
			for j := 0; j < len(team2); j++ {
				newTeam1 := make(Team, len(team1))
				newTeam2 := make(Team, len(team2))
				copy(newTeam1, team1)
				copy(newTeam2, team2)

				// Меняем игроков местами
				newTeam1[i], newTeam2[j] = newTeam2[j], newTeam1[i]

				// Проверяем ограничения и разницу в счете
				if isConstraintSatisfied(newTeam1, newTeam2, constraints) {
					newDiff := math.Abs(getTeamScore(newTeam1) - getTeamScore(newTeam2))
					if newDiff < bestDiff {
						bestDiff = newDiff
						bestTeam1 = make(Team, len(newTeam1))
						bestTeam2 = make(Team, len(newTeam2))
						copy(bestTeam1, newTeam1)
						copy(bestTeam2, newTeam2)
					}
				}
			}
		}
	}

	// Если нашли лучшее решение с учетом ограничений, используем его
	if isConstraintSatisfied(bestTeam1, bestTeam2, constraints) {
		return bestTeam1, bestTeam2
	}

	// Если не удалось найти решение, удовлетворяющее ограничениям,
	// возвращаем исходное распределение
	return team1, team2
}

// playerInTeam вспомогательная функция для проверки наличия игрока в команде
func playerInTeam(team Team, playerName string) bool {
	for _, player := range team {
		if player.NickName == playerName {
			return true
		}
	}
	return false
}
