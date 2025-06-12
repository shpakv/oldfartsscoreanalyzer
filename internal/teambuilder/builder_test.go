package teambuilder

import (
	"math"
	"reflect"
	"sort"
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
		{
			name: "Odd number of players distribution",
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
			},
			constraints: []Constraint{},
			wantCheck: func(team1, team2 Team) bool {
				// Проверка правильного распределения при нечетном количестве
				totalPlayers := len(team1) + len(team2)
				if totalPlayers != 9 {
					return false
				}
				// Разница в размере команд должна быть не больше 1
				if math.Abs(float64(len(team1)-len(team2))) > 1 {
					return false
				}
				// Проверка баланса
				weightDiff := math.Abs(getTeamScore(team1) - getTeamScore(team2))
				return weightDiff < 150 // Увеличенный порог для нечетного случая
			},
		},
		{
			name: "Minimal team size - three players",
			players: []TeamPlayer{
				{"Player A", 1000},
				{"Player B", 2000},
				{"Player C", 1500},
			},
			constraints: []Constraint{},
			wantCheck: func(team1, team2 Team) bool {
				// Проверка минимального размера команд
				if len(team1)+len(team2) != 3 {
					return false
				}
				return math.Abs(float64(len(team1)-len(team2))) <= 1
			},
		},
		{
			name: "Complex constraints with odd number",
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
			},
			constraints: []Constraint{
				{Type: ConstraintTogether, Player1: "Player A", Player2: "Player B"},
				{Type: ConstraintSeparate, Player1: "Player C", Player2: "Player D"},
			},
			wantCheck: func(team1, team2 Team) bool {
				// Проверка всех условий
				// 1. Нечетное количество игроков
				if math.Abs(float64(len(team1)-len(team2))) > 1 {
					return false
				}
				// 2. Проверка ограничения "вместе"
				playerAInTeam1 := playerInTeam(team1, "Player A")
				playerBInTeam1 := playerInTeam(team1, "Player B")
				playerAInTeam2 := playerInTeam(team2, "Player A")
				playerBInTeam2 := playerInTeam(team2, "Player B")
				if !((playerAInTeam1 && playerBInTeam1) || (playerAInTeam2 && playerBInTeam2)) {
					return false
				}
				// 3. Проверка ограничения "раздельно"
				playerCInTeam1 := playerInTeam(team1, "Player C")
				playerDInTeam1 := playerInTeam(team1, "Player D")
				if playerCInTeam1 && playerDInTeam1 {
					return false
				}
				playerCInTeam2 := playerInTeam(team2, "Player C")
				playerDInTeam2 := playerInTeam(team2, "Player D")
				if playerCInTeam2 && playerDInTeam2 {
					return false
				}
				return true
			},
		},
		{
			name: "Edge case - all players with same score",
			players: []TeamPlayer{
				{"Player A", 1000},
				{"Player B", 1000},
				{"Player C", 1000},
				{"Player D", 1000},
				{"Player E", 1000},
			},
			constraints: []Constraint{},
			wantCheck: func(team1, team2 Team) bool {
				if math.Abs(float64(len(team1)-len(team2))) > 1 {
					return false
				}
				return math.Abs(getTeamScore(team1)-getTeamScore(team2)) == 1000
			},
		},
		{
			name: "Edge case - extreme score differences",
			players: []TeamPlayer{
				{"Player A", 100},
				{"Player B", 200},
				{"Player C", 5000},
				{"Player D", 6000},
				{"Player E", 150},
			},
			constraints: []Constraint{},
			wantCheck: func(team1, team2 Team) bool {
				return math.Abs(getTeamScore(team1)-getTeamScore(team2)) < 1000
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

// Вспомогательные функции для тестов
func sortTeamByName(team Team) Team {
	sortedTeam := make(Team, len(team))
	copy(sortedTeam, team)
	sort.Slice(sortedTeam, func(i, j int) bool {
		return sortedTeam[i].NickName < sortedTeam[j].NickName
	})
	return sortedTeam
}

func calculateTeamScoreDifference(team1, team2 Team) float64 {
	return math.Abs(team1.Score() - team2.Score())
}

// Тест на базовое распределение без ограничений
func TestBasicDistribution(t *testing.T) {
	builder := &TeamBuilder{}
	config := &TeamConfiguration{
		Players: Team{
			{NickName: "Player1", Score: 1000},
			{NickName: "Player2", Score: 2000},
			{NickName: "Player3", Score: 1500},
			{NickName: "Player4", Score: 1800},
		},
		Constraints: Constraints{},
	}

	team1, team2 := builder.Build(config)

	// Проверяем, что все игроки распределены
	if len(team1)+len(team2) != len(config.Players) {
		t.Errorf("Not all players were distributed. Expected %d, got %d",
			len(config.Players), len(team1)+len(team2))
	}

	// Проверяем, что команды примерно равны по силе
	diff := calculateTeamScoreDifference(team1, team2)
	if diff > 500 {
		t.Errorf("Teams are not balanced. Score difference: %.2f", diff)
	}
}

// Тест на ограничение ConstraintTogether
func TestConstraintTogether(t *testing.T) {
	builder := &TeamBuilder{}
	config := &TeamConfiguration{
		Players: Team{
			{NickName: "Player1", Score: 1000},
			{NickName: "Player2", Score: 2000},
			{NickName: "Player3", Score: 1500},
			{NickName: "Player4", Score: 1800},
		},
		Constraints: Constraints{
			{Type: ConstraintTogether, Player1: "Player1", Player2: "Player2"},
		},
	}

	team1, team2 := builder.Build(config)

	// Проверяем, что игроки Player1 и Player2 находятся в одной команде
	player1InTeam1 := playerInTeam(team1, "Player1")
	player2InTeam1 := playerInTeam(team1, "Player2")
	player1InTeam2 := playerInTeam(team2, "Player1")
	player2InTeam2 := playerInTeam(team2, "Player2")

	if !((player1InTeam1 && player2InTeam1) || (player1InTeam2 && player2InTeam2)) {
		t.Errorf("ConstraintTogether not satisfied. Players should be in the same team")
		t.Logf("Team1: %v", team1)
		t.Logf("Team2: %v", team2)
	}
}

// Тест на ограничение ConstraintSeparate
func TestConstraintSeparate(t *testing.T) {
	builder := &TeamBuilder{}
	config := &TeamConfiguration{
		Players: Team{
			{NickName: "Player1", Score: 1000},
			{NickName: "Player2", Score: 2000},
			{NickName: "Player3", Score: 1500},
			{NickName: "Player4", Score: 1800},
		},
		Constraints: Constraints{
			{Type: ConstraintSeparate, Player1: "Player1", Player2: "Player2"},
		},
	}

	team1, team2 := builder.Build(config)

	// Проверяем, что игроки Player1 и Player2 находятся в разных командах
	player1InTeam1 := playerInTeam(team1, "Player1")
	player2InTeam1 := playerInTeam(team1, "Player2")
	player1InTeam2 := playerInTeam(team2, "Player1")
	player2InTeam2 := playerInTeam(team2, "Player2")

	if (player1InTeam1 && player2InTeam1) || (player1InTeam2 && player2InTeam2) {
		t.Errorf("ConstraintSeparate not satisfied. Players should be in different teams")
		t.Logf("Team1: %v", team1)
		t.Logf("Team2: %v", team2)
	}
}

// Тест на множественные ограничения
func TestMultipleConstraints(t *testing.T) {
	builder := &TeamBuilder{}
	config := &TeamConfiguration{
		Players: Team{
			{NickName: "Player1", Score: 1000},
			{NickName: "Player2", Score: 2000},
			{NickName: "Player3", Score: 1500},
			{NickName: "Player4", Score: 1800},
			{NickName: "Player5", Score: 1600},
			{NickName: "Player6", Score: 1700},
		},
		Constraints: Constraints{
			{Type: ConstraintTogether, Player1: "Player1", Player2: "Player2"},
			{Type: ConstraintSeparate, Player1: "Player3", Player2: "Player4"},
			{Type: ConstraintTogether, Player1: "Player5", Player2: "Player6"},
		},
	}

	team1, _ := builder.Build(config)

	// Проверяем все ограничения
	player1InTeam1 := playerInTeam(team1, "Player1")
	player2InTeam1 := playerInTeam(team1, "Player2")
	player3InTeam1 := playerInTeam(team1, "Player3")
	player4InTeam1 := playerInTeam(team1, "Player4")
	player5InTeam1 := playerInTeam(team1, "Player5")
	player6InTeam1 := playerInTeam(team1, "Player6")

	// Проверка ConstraintTogether для Player1 и Player2
	if !((player1InTeam1 && player2InTeam1) || (!player1InTeam1 && !player2InTeam1)) {
		t.Errorf("ConstraintTogether not satisfied for Player1 and Player2")
	}

	// Проверка ConstraintSeparate для Player3 и Player4
	if (player3InTeam1 && player4InTeam1) || (!player3InTeam1 && !player4InTeam1) {
		t.Errorf("ConstraintSeparate not satisfied for Player3 and Player4")
	}

	// Проверка ConstraintTogether для Player5 и Player6
	if !((player5InTeam1 && player6InTeam1) || (!player5InTeam1 && !player6InTeam1)) {
		t.Errorf("ConstraintTogether not satisfied for Player5 and Player6")
	}
}

// Тест на конфликтующие ограничения
func TestConflictingConstraints(t *testing.T) {
	builder := &TeamBuilder{}
	config := &TeamConfiguration{
		Players: Team{
			{NickName: "Player1", Score: 1000},
			{NickName: "Player2", Score: 2000},
			{NickName: "Player3", Score: 1500},
			{NickName: "Player4", Score: 1800},
		},
		Constraints: Constraints{
			{Type: ConstraintTogether, Player1: "Player1", Player2: "Player2"},
			{Type: ConstraintTogether, Player1: "Player2", Player2: "Player3"},
			{Type: ConstraintSeparate, Player1: "Player1", Player2: "Player3"},
		},
	}

	team1, team2 := builder.Build(config)

	// В этом случае невозможно удовлетворить все ограничения
	// Проверяем, что функция вернула какое-то решение
	if len(team1)+len(team2) != len(config.Players) {
		t.Errorf("Not all players were distributed with conflicting constraints")
	}

	// Логируем полученные команды для анализа
	t.Logf("Team1: %v", team1)
	t.Logf("Team2: %v", team2)
}

// Тест на нечетное количество игроков
func TestOddNumberOfPlayers(t *testing.T) {
	builder := &TeamBuilder{}
	config := &TeamConfiguration{
		Players: Team{
			{NickName: "Player1", Score: 1000},
			{NickName: "Player2", Score: 2000},
			{NickName: "Player3", Score: 1500},
			{NickName: "Player4", Score: 1800},
			{NickName: "Player5", Score: 1600},
		},
		Constraints: Constraints{},
	}

	team1, team2 := builder.Build(config)

	// Проверяем, что все игроки распределены
	if len(team1)+len(team2) != len(config.Players) {
		t.Errorf("Not all players were distributed. Expected %d, got %d",
			len(config.Players), len(team1)+len(team2))
	}

	// Проверяем, что размеры команд отличаются максимум на 1
	if math.Abs(float64(len(team1)-len(team2))) > 1 {
		t.Errorf("Teams are not balanced by size. Team1: %d, Team2: %d",
			len(team1), len(team2))
	}
}

// Тест на отсутствующих в ограничениях игроков
func TestMissingPlayersInConstraints(t *testing.T) {
	builder := &TeamBuilder{}
	config := &TeamConfiguration{
		Players: Team{
			{NickName: "Player1", Score: 1000},
			{NickName: "Player2", Score: 2000},
		},
		Constraints: Constraints{
			{Type: ConstraintTogether, Player1: "Player1", Player2: "NonExistentPlayer"},
		},
	}

	team1, team2 := builder.Build(config)

	// Проверяем, что функция не паникует и возвращает какое-то решение
	if len(team1)+len(team2) != len(config.Players) {
		t.Errorf("Not all players were distributed with missing players in constraints")
	}
}

// Тест на оптимизацию команд
func TestTeamOptimization(t *testing.T) {
	builder := &TeamBuilder{}
	config := &TeamConfiguration{
		Players: Team{
			{NickName: "Player1", Score: 1000},
			{NickName: "Player2", Score: 2000},
			{NickName: "Player3", Score: 1500},
			{NickName: "Player4", Score: 1800},
			{NickName: "Player5", Score: 1600},
			{NickName: "Player6", Score: 1700},
		},
		Constraints: Constraints{},
	}

	team1, team2 := builder.Build(config)

	// Проверяем, что разница в силе команд минимальна
	diff := calculateTeamScoreDifference(team1, team2)

	// Создаем команды вручную и проверяем, что алгоритм дает лучший или такой же результат
	manualTeam1 := Team{
		{NickName: "Player2", Score: 2000},
		{NickName: "Player5", Score: 1600},
		{NickName: "Player1", Score: 1000},
	}
	manualTeam2 := Team{
		{NickName: "Player4", Score: 1800},
		{NickName: "Player6", Score: 1700},
		{NickName: "Player3", Score: 1500},
	}

	manualDiff := calculateTeamScoreDifference(manualTeam1, manualTeam2)

	if diff > manualDiff+100 { // Допускаем небольшую погрешность
		t.Errorf("Team optimization is not optimal. Got diff: %.2f, expected around: %.2f",
			diff, manualDiff)
		t.Logf("Team1: %v", team1)
		t.Logf("Team2: %v", team2)
	}
}

// Тест на пустой список игроков
func TestEmptyPlayersList(t *testing.T) {
	builder := &TeamBuilder{}
	config := &TeamConfiguration{
		Players:     Team{},
		Constraints: Constraints{},
	}

	team1, team2 := builder.Build(config)

	if len(team1) != 0 || len(team2) != 0 {
		t.Errorf("Expected empty teams for empty player list")
	}
}

// Тест на идемпотентность - повторный запуск с теми же данными
func TestIdempotence(t *testing.T) {
	builder := &TeamBuilder{}
	config := &TeamConfiguration{
		Players: Team{
			{NickName: "Player1", Score: 1000},
			{NickName: "Player2", Score: 2000},
			{NickName: "Player3", Score: 1500},
			{NickName: "Player4", Score: 1800},
		},
		Constraints: Constraints{
			{Type: ConstraintTogether, Player1: "Player1", Player2: "Player2"},
		},
	}

	team1First, team2First := builder.Build(config)
	team1Second, team2Second := builder.Build(config)

	// Сортируем команды для сравнения
	sortedTeam1First := sortTeamByName(team1First)
	sortedTeam2First := sortTeamByName(team2First)
	sortedTeam1Second := sortTeamByName(team1Second)
	sortedTeam2Second := sortTeamByName(team2Second)

	// Проверяем, что результаты одинаковые или хотя бы эквивалентные по балансу
	if !reflect.DeepEqual(sortedTeam1First, sortedTeam1Second) ||
		!reflect.DeepEqual(sortedTeam2First, sortedTeam2Second) {
		diffFirst := calculateTeamScoreDifference(team1First, team2First)
		diffSecond := calculateTeamScoreDifference(team1Second, team2Second)

		if math.Abs(diffFirst-diffSecond) > 100 {
			t.Errorf("Algorithm is not idempotent. First run diff: %.2f, Second run diff: %.2f",
				diffFirst, diffSecond)
		}
	}
}

// Тест на использование метода Team.Score() вместо getTeamScore
func TestTeamScoreMethod(t *testing.T) {
	team := Team{
		{NickName: "Player1", Score: 1000},
		{NickName: "Player2", Score: 2000},
		{NickName: "Player3", Score: 1500},
	}

	// Проверяем, что метод Team.Score() возвращает то же значение, что и getTeamScore
	methodScore := team.Score()
	functionScore := getTeamScore(team)

	if methodScore != functionScore {
		t.Errorf("Team.Score() method returns different value than getTeamScore function. "+
			"Method: %.2f, Function: %.2f", methodScore, functionScore)
	}
}

// Тест на проверку обработки ограничений с несуществующими игроками
func TestConstraintsWithNonexistentPlayers(t *testing.T) {
	builder := &TeamBuilder{}
	config := &TeamConfiguration{
		Players: Team{
			{NickName: "Player1", Score: 1000},
			{NickName: "Player2", Score: 2000},
			{NickName: "Player3", Score: 1500},
		},
		Constraints: Constraints{
			{Type: ConstraintTogether, Player1: "Player1", Player2: "NonExistentPlayer"},
			{Type: ConstraintSeparate, Player1: "NonExistentPlayer1", Player2: "NonExistentPlayer2"},
		},
	}

	// Функция не должна паниковать
	team1, team2 := builder.Build(config)

	// Проверяем, что все существующие игроки распределены
	if len(team1)+len(team2) != len(config.Players) {
		t.Errorf("Not all players were distributed. Expected %d, got %d",
			len(config.Players), len(team1)+len(team2))
	}
}

// Тест на обработку циклических зависимостей в ограничениях
func TestCyclicConstraints(t *testing.T) {
	builder := &TeamBuilder{}
	config := &TeamConfiguration{
		Players: Team{
			{NickName: "Player1", Score: 1000},
			{NickName: "Player2", Score: 2000},
			{NickName: "Player3", Score: 1500},
			{NickName: "Player4", Score: 1800},
		},
		Constraints: Constraints{
			{Type: ConstraintTogether, Player1: "Player1", Player2: "Player2"},
			{Type: ConstraintTogether, Player1: "Player2", Player2: "Player3"},
			{Type: ConstraintTogether, Player1: "Player3", Player2: "Player1"},
		},
	}

	team1, team2 := builder.Build(config)

	// Проверяем, что все игроки распределены
	if len(team1)+len(team2) != len(config.Players) {
		t.Errorf("Not all players were distributed. Expected %d, got %d",
			len(config.Players), len(team1)+len(team2))
	}

	// Проверяем, что все три игрока находятся в одной команде
	player1InTeam1 := playerInTeam(team1, "Player1")
	player2InTeam1 := playerInTeam(team1, "Player2")
	player3InTeam1 := playerInTeam(team1, "Player3")

	if !((player1InTeam1 && player2InTeam1 && player3InTeam1) ||
		(!player1InTeam1 && !player2InTeam1 && !player3InTeam1)) {
		t.Errorf("Cyclic ConstraintTogether not satisfied. Players 1, 2, and 3 should be in the same team")
		t.Logf("Team1: %v", team1)
		t.Logf("Team2: %v", team2)
	}
}
