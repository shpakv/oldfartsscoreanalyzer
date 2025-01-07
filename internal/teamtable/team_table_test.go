package teamtable

import (
	"github.com/stretchr/testify/assert"
	"oldfartscounter/internal/teambuilder"
	"testing"
)

func TestNewTeamTable(t *testing.T) {
	team1 := teambuilder.Team{
		{NickName: "Player1", Score: 100},
		{NickName: "Player2", Score: 200},
	}
	team2 := teambuilder.Team{
		{NickName: "Player3", Score: 150},
		{NickName: "Player4", Score: 50},
	}

	expectedHeaders := []string{"Team 1", "Team 2"}
	expectedRows := [][]string{
		{"Player1", "Player3"},
		{"Player2", "Player4"},
	}
	expectedTeamScore := []string{"300.00", "200.00"}
	expectedScoreDifference := "100.00"

	table := NewTeamTable(team1, team2)

	assert.Equal(t, expectedHeaders, table.Headers, "Headers mismatch")
	assert.Equal(t, expectedRows, table.Rows, "Rows mismatch")
	assert.Equal(t, expectedTeamScore, table.TeamScore, "TeamScore mismatch")
	assert.Equal(t, expectedScoreDifference, table.ScoreDifference, "ScoreDifference mismatch")
}
