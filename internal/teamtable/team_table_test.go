package teamtable

import (
	"oldfartscounter/internal/teambuilder"
	"testing"

	"github.com/stretchr/testify/assert"
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

	table := NewTeamTable(team1, team2, "SorryBro")

	assert.Equal(t, expectedHeaders, table.Headers, "Headers mismatch")
	assert.Equal(t, expectedRows, table.Rows, "Rows mismatch")
	assert.Equal(t, expectedTeamScore, table.TeamScore, "TeamScore mismatch")
	assert.Equal(t, expectedScoreDifference, table.ScoreDifference, "ScoreDifference mismatch")
}

func TestNewTeamTableMultiple_FourTeams(t *testing.T) {
	team1 := teambuilder.Team{
		{NickName: "Player1", Score: 100},
		{NickName: "Player2", Score: 200},
	}
	team2 := teambuilder.Team{
		{NickName: "Player3", Score: 150},
		{NickName: "Player4", Score: 50},
	}
	team3 := teambuilder.Team{
		{NickName: "Player5", Score: 120},
		{NickName: "Player6", Score: 80},
	}
	team4 := teambuilder.Team{
		{NickName: "Player7", Score: 90},
		{NickName: "Player8", Score: 110},
	}

	teams := []teambuilder.Team{team1, team2, team3, team4}

	expectedHeaders := []string{"Team 1", "Team 2", "Team 3", "Team 4"}
	expectedRows := [][]string{
		{"Player1", "Player3", "Player5", "Player7"},
		{"Player2", "Player4", "Player6", "Player8"},
	}
	expectedTeamScore := []string{"300.00", "200.00", "200.00", "200.00"}
	expectedScoreDifference := "100.00" // max 300 - min 200

	table := NewTeamTableMultiple(teams, "")

	assert.Equal(t, expectedHeaders, table.Headers, "Headers mismatch")
	assert.Equal(t, expectedRows, table.Rows, "Rows mismatch")
	assert.Equal(t, expectedTeamScore, table.TeamScore, "TeamScore mismatch")
	assert.Equal(t, expectedScoreDifference, table.ScoreDifference, "ScoreDifference mismatch")
}

func TestNewTeamTableMultiple_TwoTeams(t *testing.T) {
	team1 := teambuilder.Team{
		{NickName: "Player1", Score: 100},
		{NickName: "Player2", Score: 200},
	}
	team2 := teambuilder.Team{
		{NickName: "Player3", Score: 150},
		{NickName: "Player4", Score: 50},
	}

	teams := []teambuilder.Team{team1, team2}

	expectedHeaders := []string{"Team 1", "Team 2"}
	expectedRows := [][]string{
		{"Player1", "Player3"},
		{"Player2", "Player4"},
	}
	expectedTeamScore := []string{"300.00", "200.00"}
	expectedScoreDifference := "100.00"

	table := NewTeamTableMultiple(teams, "")

	assert.Equal(t, expectedHeaders, table.Headers, "Headers mismatch")
	assert.Equal(t, expectedRows, table.Rows, "Rows mismatch")
	assert.Equal(t, expectedTeamScore, table.TeamScore, "TeamScore mismatch")
	assert.Equal(t, expectedScoreDifference, table.ScoreDifference, "ScoreDifference mismatch")
}

func TestNewTeamTableMultiple_UnevenTeams(t *testing.T) {
	team1 := teambuilder.Team{
		{NickName: "Player1", Score: 100},
		{NickName: "Player2", Score: 200},
		{NickName: "Player3", Score: 150},
	}
	team2 := teambuilder.Team{
		{NickName: "Player4", Score: 50},
	}
	team3 := teambuilder.Team{
		{NickName: "Player5", Score: 120},
		{NickName: "Player6", Score: 80},
	}
	team4 := teambuilder.Team{
		{NickName: "Player7", Score: 90},
		{NickName: "Player8", Score: 110},
		{NickName: "Player9", Score: 100},
	}

	teams := []teambuilder.Team{team1, team2, team3, team4}

	expectedHeaders := []string{"Team 1", "Team 2", "Team 3", "Team 4"}
	expectedRows := [][]string{
		{"Player1", "Player4", "Player5", "Player7"},
		{"Player2", "", "Player6", "Player8"},
		{"Player3", "", "", "Player9"},
	}
	expectedTeamScore := []string{"450.00", "50.00", "200.00", "300.00"}

	table := NewTeamTableMultiple(teams, "")

	assert.Equal(t, expectedHeaders, table.Headers, "Headers mismatch")
	assert.Equal(t, expectedRows, table.Rows, "Rows mismatch")
	assert.Equal(t, expectedTeamScore, table.TeamScore, "TeamScore mismatch")
	// Max 450 - Min 50 = 400
	assert.Equal(t, "400.00", table.ScoreDifference, "ScoreDifference mismatch")
}

func TestNewTeamTableMultiple_EmptyTeams(t *testing.T) {
	teams := []teambuilder.Team{}

	table := NewTeamTableMultiple(teams, "")

	assert.Empty(t, table.Headers, "Headers should be empty")
	assert.Empty(t, table.Rows, "Rows should be empty")
	assert.Empty(t, table.TeamScore, "TeamScore should be empty")
	assert.Equal(t, "0.00", table.ScoreDifference, "ScoreDifference should be 0.00")
}
