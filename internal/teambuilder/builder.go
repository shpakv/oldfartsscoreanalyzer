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
	if len(players)%2 != 0 {
		teamSize++
	}
	// Сортируем игроков по убыванию веса
	sort.Slice(config.Players, func(i, j int) bool {
		return players[i].Score > players[j].Score
	})

	var bestTeam1, bestTeam2 Team
	minDiff := math.Inf(1)

	// Функция проверки ограничений
	isConstraintSatisfied := func(team1, team2 Team) bool {
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

	var backtrack func(int, Team, Team, float64, float64)
	backtrack = func(index int, team1, team2 Team, weight1, weight2 float64) {
		// Базовый случай - достигнут конец списка игроков
		if index == len(players) {
			// Проверяем, что все игроки распределены
			if len(team1)+len(team2) == len(players) {
				// Разрешаем небольшую разницу в размерах команд
				if math.Abs(float64(len(team1)-len(team2))) <= 1 {
					// Проверяем ограничения
					if isConstraintSatisfied(team1, team2) {
						diff := math.Abs(weight1 - weight2)
						if diff < minDiff {
							minDiff = diff
							bestTeam1 = make(Team, len(team1))
							bestTeam2 = make(Team, len(team2))
							copy(bestTeam1, team1)
							copy(bestTeam2, team2)
						}
					}
				}
			}
			return
		}

		// Пробуем добавить игрока в первую команду
		if len(team1) < teamSize {
			backtrack(
				index+1,
				append(team1, players[index]),
				team2,
				weight1+players[index].Score,
				weight2,
			)
		}

		// Пробуем добавить игрока во вторую команду
		if len(team2) < teamSize {
			backtrack(
				index+1,
				team1,
				append(team2, players[index]),
				weight1,
				weight2+players[index].Score,
			)
		}
	}

	// Запускаем поиск
	backtrack(0, Team{}, Team{}, 0, 0)

	return bestTeam1, bestTeam2
}

// playerInTeam вспомогательная функция для проверки наличия игрока в команде
//
// Параметры:
//   - team: команда для проверки
//   - playerName: имя игрока
//
// Возвращает:
//   - true, если игрок есть в команде
//   - false, если игрок отсутствует.
//
// Используется для проверки ограничений при распределении.
func playerInTeam(team Team, playerName string) bool {
	for _, player := range team {
		if player.NickName == playerName {
			return true
		}
	}
	return false
}
