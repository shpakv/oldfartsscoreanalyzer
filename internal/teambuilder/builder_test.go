package teambuilder

import (
	"math"
	"testing"
)

// TestTeamBuilder_Calculate тестирует функцию распределения игроков по командам.
// Покрывает следующие сценарии:
//  1. Базовое распределение без ограничений
//  2. Распределение с требованием совместить определенных игроков
//  3. Распределение с требованием разделить определенных игроков
func TestTeamBuilder_Calculate(t *testing.T) {
	c := &TeamBuilder{}

	// testCases описывает различные сценарии тестирования
	testCases := []struct {
		name        string                       // Название тестового сценария
		players     Team                         // Список игроков
		constraints []Constraint                 // Ограничения на распределение
		wantCheck   func(team1, team2 Team) bool // Функция проверки результата
	}{
		{
			// Базовый сценарий: равномерное распределение без дополнительных условий
			name: "Balanced team distribution without constraints",
			players: []TeamPlayer{
				{"Player A", 928.57},
				{"Player B", 1863.03},
				{"Player C", 1492.15},
				{"Player D", 1512.89},
				{"Player E", 1290.22},
				{"Player F", 1552.77},
				{"Player G", 2520.86},
				{"Player H", 2383.79},
				{"Player I", 1936.04},
				{"Player J", 1905.34},
			},
			constraints: []Constraint{},
			wantCheck: func(team1, team2 Team) bool {
				// Проверка размера команд
				if len(team1) != 5 || len(team2) != 5 {
					return false
				}

				// Вычисление суммарного веса команд
				team1Weight := 0.0
				team2Weight := 0.0
				for _, player := range team1 {
					team1Weight += player.Score
				}
				for _, player := range team2 {
					team2Weight += player.Score
				}

				// Проверка разницы весов
				weightDiff := math.Abs(team1Weight - team2Weight)
				return weightDiff < 100 // Допустимая разница в весе
			},
		},
		{
			// Сценарий с требованием совместить двух игроков
			name: "Forcing specific players to be in the same team",
			players: []TeamPlayer{
				{"Player A", 928.57},
				{"Player B", 1863.03},
				{"Player C", 1492.15},
				{"Player D", 1512.89},
				{"Player E", 1290.22},
				{"Player F", 1552.77},
				{"Player G", 2520.86},
				{"Player H", 2383.79},
				{"Player I", 1936.04},
				{"Player J", 1905.34},
			},
			constraints: []Constraint{
				{Type: ConstraintTogether, Player1: "Player I", Player2: "Player G"},
			},
			wantCheck: func(team1, team2 Team) bool {
				// Проверка, что указанные игроки вместе
				playerIInTeam1 := playerInTeam(team1, "Player I")
				playerGInTeam1 := playerInTeam(team1, "Player G")

				playerIInTeam2 := playerInTeam(team2, "Player I")
				playerGInTeam2 := playerInTeam(team2, "Player G")

				return (playerIInTeam1 && playerGInTeam1) ||
					(playerIInTeam2 && playerGInTeam2)
			},
		},
		{
			// Сценарий с требованием разделить двух игроков
			name: "Forcing specific players to be in different teams",
			players: []TeamPlayer{
				{"Player A", 928.57},
				{"Player B", 1863.03},
				{"Player C", 1492.15},
				{"Player D", 1512.89},
				{"Player E", 1290.22},
				{"Player F", 1552.77},
				{"Player G", 2520.86},
				{"Player H", 2383.79},
				{"Player I", 1936.04},
				{"Player J", 1905.34},
			},
			constraints: []Constraint{
				{Type: ConstraintSeparate, Player1: "Player G", Player2: "Player H"},
			},
			wantCheck: func(team1, team2 Team) bool {
				// Проверка, что указанные игроки в разных командах
				playerGInTeam1 := playerInTeam(team1, "Player G")
				playerHInTeam1 := playerInTeam(team1, "Player H")

				playerGInTeam2 := playerInTeam(team2, "Player G")
				playerHInTeam2 := playerInTeam(team2, "Player H")

				return (playerGInTeam1 && !playerHInTeam1) ||
					(playerGInTeam2 && !playerHInTeam2)
			},
		},
	}

	// Выполнение тестов для каждого сценария
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Создание конфигурации
			config := &TeamConfiguration{
				Players:     tc.players,
				Constraints: tc.constraints,
			}

			// Распределение команд
			team1, team2 := c.Build(config)

			// Вывод информации о командах
			printTeamDetails(t, team1, team2)

			// Проверка соответствия требованиям
			if !tc.wantCheck(team1, team2) {
				t.Errorf("Test case %s failed", tc.name)

				// Дополнительная отладочная информация
				t.Logf("Team 1:")
				for _, player := range team1 {
					t.Logf("%s (%.2f)", player.NickName, player.Score)
				}
				t.Logf("Team 2:")
				for _, player := range team2 {
					t.Logf("%s (%.2f)", player.NickName, player.Score)
				}
			}
		})
	}
}

// printTeamDetails выводит подробную информацию о сформированных командах
func printTeamDetails(t *testing.T, team1, team2 Team) {
	t.Log("=== Team 1 ===")
	totalWeight1 := 0.0
	for _, player := range team1 {
		t.Logf("%s (%.2f)", player.NickName, player.Score)
		totalWeight1 += player.Score
	}
	t.Logf("= Total Score %.2f =\n", totalWeight1)

	t.Log("=== Team 2 ===")
	totalWeight2 := 0.0
	for _, player := range team2 {
		t.Logf("%s (%.2f)", player.NickName, player.Score)
		totalWeight2 += player.Score
	}
	t.Logf("= Total Score: %.2f =", totalWeight2)
	t.Logf("**** Score Difference: %.2f ****", math.Abs(totalWeight1-totalWeight2))
}

// BenchmarkTeamBuilder_Calculate тестирует производительность распределения команд.
// Замеряет время выполнения для различных входных данных
func BenchmarkTeamBuilder_Calculate(b *testing.B) {
	c := &TeamBuilder{}

	// Тестовый набор игроков
	players := []TeamPlayer{
		{"Player A", 928.57},
		{"Player B", 1863.03},
		{"Player C", 1492.15},
		{"Player D", 1512.89},
		{"Player E", 1290.22},
		{"Player F", 1552.77},
		{"Player G", 2520.86},
		{"Player H", 2383.79},
		{"Player I", 1936.04},
		{"Player J", 1905.34},
	}

	// Тестовые ограничения
	constraints := []Constraint{
		{Type: ConstraintTogether, Player1: "Player I", Player2: "Player J"},
		{Type: ConstraintSeparate, Player1: "Player G", Player2: "Player H"},
	}

	// Запуск бенчмарка
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Создание конфигурации
		config := &TeamConfiguration{
			Players:     players,
			Constraints: constraints,
		}
		c.Build(config)
	}
}
