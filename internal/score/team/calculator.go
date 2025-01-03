package team

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type Player struct {
	ID     string
	Rating float64
}

type Constraint struct {
	MustTogether []string // Игроки, которые должны быть в одной команде
	MustSeparate []string // Игроки, которые должны быть в разных командах
}

type Builder struct {
	players     []Player
	constraints Constraint
}

type TeamAllocation struct {
	Team1       []Player
	Team2       []Player
	Team1Rating float64
	Team2Rating float64
}

func NewTeamBuilder(players []Player, constraints Constraint) *Builder {
	return &Builder{
		players:     players,
		constraints: constraints,
	}
}

func (tb *Builder) BuildTeams() TeamAllocation {
	// Сортируем игроков по рейтингу в порядке убывания
	sort.Slice(tb.players, func(i, j int) bool {
		return tb.players[i].Rating > tb.players[j].Rating
	})

	// Применяем жадный алгоритм сbacktracking
	return tb.distributeTeams()
}

func (tb *Builder) distributeTeams() TeamAllocation {
	bestAllocation := TeamAllocation{}
	minRatingDifference := math.Inf(1)

	// Генерируем все возможные комбинации с учетом ограничений
	combinations := tb.generateTeamCombinations()

	for _, combination := range combinations {
		if !tb.validateConstraints(combination) {
			continue
		}

		team1Rating := calculateTeamRating(combination.Team1)
		team2Rating := calculateTeamRating(combination.Team2)
		ratingDifference := math.Abs(team1Rating - team2Rating)

		if ratingDifference < minRatingDifference {
			minRatingDifference = ratingDifference
			bestAllocation = TeamAllocation{
				Team1:       combination.Team1,
				Team2:       combination.Team2,
				Team1Rating: team1Rating,
				Team2Rating: team2Rating,
			}
		}
	}

	return bestAllocation
}

func (tb *Builder) generateTeamCombinations() []TeamAllocation {
	n := len(tb.players)
	combinations := []TeamAllocation{}

	// Генерируем все возможные подмножества
	for i := 0; i < (1 << n); i++ {
		team1 := []Player{}
		team2 := []Player{}

		for j := 0; j < n; j++ {
			if (i & (1 << j)) > 0 {
				team1 = append(team1, tb.players[j])
			} else {
				team2 = append(team2, tb.players[j])
			}
		}

		// Убеждаемся, что команды не пустые
		if len(team1) > 0 && len(team2) > 0 {
			combinations = append(combinations, TeamAllocation{
				Team1: team1,
				Team2: team2,
			})
		}
	}

	return combinations
}

func (tb *Builder) validateConstraints(allocation TeamAllocation) bool {
	// Проверка обязательного нахождения вместе
	for _, mustTogetherGroup := range tb.splitPlayerGroups(tb.constraints.MustTogether) {
		if !tb.checkGroupTogether(mustTogetherGroup, allocation) {
			return false
		}
	}

	// Проверка обязательного разделения
	for _, mustSeparateGroup := range tb.splitPlayerGroups(tb.constraints.MustSeparate) {
		if !tb.checkGroupSeparate(mustSeparateGroup, allocation) {
			return false
		}
	}

	return true
}

func (tb *Builder) splitPlayerGroups(groups []string) [][]string {
	result := [][]string{}
	for _, group := range groups {
		result = append(result, strings.Split(group, ","))
	}
	return result
}

func (tb *Builder) checkGroupTogether(group []string, allocation TeamAllocation) bool {
	team1ContainsGroup := tb.checkGroupInTeam(group, allocation.Team1)
	team2ContainsGroup := tb.checkGroupInTeam(group, allocation.Team2)
	return team1ContainsGroup || team2ContainsGroup
}

func (tb *Builder) checkGroupSeparate(group []string, allocation TeamAllocation) bool {
	team1ContainsGroup := tb.checkGroupInTeam(group, allocation.Team1)
	team2ContainsGroup := tb.checkGroupInTeam(group, allocation.Team2)
	return !team1ContainsGroup || !team2ContainsGroup
}

func (tb *Builder) checkGroupInTeam(group []string, team []Player) bool {
	for _, playerID := range group {
		found := false
		for _, teamPlayer := range team {
			if teamPlayer.ID == playerID {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func calculateTeamRating(team []Player) float64 {
	totalRating := 0.0
	for _, player := range team {
		totalRating += player.Rating
	}
	return totalRating
}

// Пример использования
func ExampleTeamBuilder() {
	players := []Player{
		{ID: "1", Rating: 10.0},
		{ID: "2", Rating: 9.0},
		{ID: "3", Rating: 8.0},
		{ID: "4", Rating: 7.0},
	}

	constraints := Constraint{
		MustTogether: []string{"2,3"}, // Игроки 2 и 3 вместе
		MustSeparate: []string{"1,4"}, // Игроки 1 и 4 в разных командах
	}

	teamBuilder := NewTeamBuilder(players, constraints)
	allocation := teamBuilder.BuildTeams()

	fmt.Printf("Team 1: %v (Rating: %.2f)\n", allocation.Team1, allocation.Team1Rating)
}
