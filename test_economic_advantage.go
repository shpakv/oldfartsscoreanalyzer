package main

import (
	"fmt"
	"math"
	"oldfartscounter/internal/teambuilder"
)

func main() {
	// Тестовые игроки для демонстрации адаптивной экономики
	players := []teambuilder.TeamPlayer{
		{"Player1", 2000},
		{"Player2", 1800},
		{"Player3", 1600},
		{"Player4", 1400},
		{"Player5", 1200},
		{"Player6", 1000},
		{"Player7", 800},
		{"Player8", 600},
		{"Player9", 400},
	}

	// Разные сценарии экономических настроек
	scenarios := []struct {
		name           string
		config         teambuilder.EconomicConfig
		expectedResult string
	}{
		{
			name: "Отключена экономика (равные команды)",
			config: teambuilder.EconomicConfig{
				Enabled: false,
			},
			expectedResult: "Стандартный баланс по рейтингу",
		},
		{
			name: "Консервативная (для 9v10)",
			config: teambuilder.EconomicConfig{
				Enabled:        true,
				BasePercentage: 15.0,
				MaxPercentage:  2.0,
				MinPercentage:  1.0,
			},
			expectedResult: "15%/9 = 1.67% → 2% (ограничение min)",
		},
		{
			name: "Сбалансированная (для 4v5)",
			config: teambuilder.EconomicConfig{
				Enabled:        true,
				BasePercentage: 20.0,
				MaxPercentage:  6.0,
				MinPercentage:  2.0,
			},
			expectedResult: "20%/4 = 5.0% за недостающего игрока",
		},
		{
			name: "Агрессивная (для малых команд)",
			config: teambuilder.EconomicConfig{
				Enabled:        true,
				BasePercentage: 30.0,
				MaxPercentage:  8.0,
				MinPercentage:  3.0,
			},
			expectedResult: "30%/4 = 7.5% за недостающего игрока",
		},
	}

	builder := &teambuilder.TeamBuilder{}

	fmt.Println("🎮 Демонстрация адаптивной экономической системы CS2")
	fmt.Println("============================================================")

	for i, scenario := range scenarios {
		fmt.Printf("\n%d. %s\n", i+1, scenario.name)
		fmt.Printf("   %s\n", scenario.expectedResult)

		config := &teambuilder.TeamConfiguration{
			Players:        players,
			Constraints:    []teambuilder.Constraint{},
			EconomicConfig: scenario.config,
		}

		team1, team2 := builder.Build(config)

		fmt.Printf("   Команда 1: %d игроков, рейтинг %.1f\n", len(team1), team1.Score())
		fmt.Printf("   Команда 2: %d игроков, рейтинг %.1f\n", len(team2), team2.Score())

		// Показываем эффективные счета с учетом экономики
		if scenario.config.Enabled && len(team1) != len(team2) {
			// Импортируем функцию из teambuilder пакета
			effectiveScore1, effectiveScore2 := teambuilder.GetEffectiveTeamScoreWithConfig(team1, team2, scenario.config)
			fmt.Printf("   Эффективные счета: %.1f vs %.1f (разница: %.1f)\n",
				effectiveScore1, effectiveScore2,
				math.Abs(effectiveScore1-effectiveScore2))
		}
	}

	fmt.Println("\n🔍 Формула адаптивной экономики:")
	fmt.Println("   percentage = BasePercentage / SmallerTeamSize")
	fmt.Println("   Ограничения: min ≤ percentage ≤ max")
	fmt.Println("   Бонус = TeamScore × percentage × PlayerDifference")
}
