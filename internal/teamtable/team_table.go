package teamtable

import (
	"fmt"
	"oldfartscounter/internal/teambuilder"
)

type TeamTable struct {
	Headers         []string
	Rows            [][]string
	TeamScore       []string
	ScoreDifference string
	SorryBro        string
}

func NewTeamTable(team1, team2 teambuilder.Team, sorryBro string) *TeamTable {
	scoreDifference := team1.Score() - team2.Score()
	if scoreDifference < 0 {
		scoreDifference = -scoreDifference
	}

	maxRows := len(team1)
	if len(team2) > maxRows {
		maxRows = len(team2)
	}

	rows := make([][]string, maxRows)
	for i := 0; i < maxRows; i++ {
		var row []string
		if i < len(team1) {
			row = append(row, team1[i].NickName)
		} else {
			row = append(row, "")
		}
		if i < len(team2) {
			row = append(row, team2[i].NickName)
		} else {
			row = append(row, "")
		}
		rows[i] = row
	}

	return &TeamTable{
		Headers:         []string{"Team 1", "Team 2"},
		Rows:            rows,
		TeamScore:       []string{fmt.Sprintf("%.2f", team1.Score()), fmt.Sprintf("%.2f", team2.Score())},
		ScoreDifference: fmt.Sprintf("%.2f", scoreDifference),
		SorryBro:        sorryBro,
	}
}

// NewTeamTableMultiple создает таблицу для переменного количества команд
func NewTeamTableMultiple(teams []teambuilder.Team, sorryBro string) *TeamTable {
	if len(teams) == 0 {
		return &TeamTable{
			Headers:         []string{},
			Rows:            [][]string{},
			TeamScore:       []string{},
			ScoreDifference: "0.00",
			SorryBro:        sorryBro,
		}
	}

	// Создаем заголовки
	headers := make([]string, len(teams))
	for i := range teams {
		headers[i] = fmt.Sprintf("Team %d", i+1)
	}

	// Находим максимальное количество строк
	maxRows := 0
	for _, team := range teams {
		if len(team) > maxRows {
			maxRows = len(team)
		}
	}

	// Создаем строки
	rows := make([][]string, maxRows)
	for i := 0; i < maxRows; i++ {
		row := make([]string, len(teams))
		for j, team := range teams {
			if i < len(team) {
				row[j] = team[i].NickName
			} else {
				row[j] = ""
			}
		}
		rows[i] = row
	}

	// Вычисляем счета команд
	teamScores := make([]string, len(teams))
	scores := make([]float64, len(teams))
	for i, team := range teams {
		scores[i] = team.Score()
		teamScores[i] = fmt.Sprintf("%.2f", scores[i])
	}

	// Вычисляем разницу между максимальной и минимальной командой
	minScore := scores[0]
	maxScore := scores[0]
	for _, score := range scores {
		if score < minScore {
			minScore = score
		}
		if score > maxScore {
			maxScore = score
		}
	}
	scoreDifference := maxScore - minScore

	return &TeamTable{
		Headers:         headers,
		Rows:            rows,
		TeamScore:       teamScores,
		ScoreDifference: fmt.Sprintf("%.2f", scoreDifference),
		SorryBro:        sorryBro,
	}
}
