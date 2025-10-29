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
