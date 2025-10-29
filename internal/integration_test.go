package internal

import (
	"oldfartscounter/internal/teambuilder"
	"oldfartscounter/internal/teamtable"
	"oldfartscounter/internal/telegram"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Integration tests that verify the full flow from team building to formatting

// Mock repository for integration tests
type mockIntegrationRepository struct{}

func (m *mockIntegrationRepository) FindByName(name string) *teambuilder.Player {
	scores := map[string]float64{
		"Alice":   1000,
		"Bob":     2000,
		"Charlie": 1500,
		"Dave":    1800,
		"Eve":     1600,
		"Frank":   1700,
		"Grace":   1400,
		"Henry":   1900,
	}
	if score, ok := scores[name]; ok {
		return &teambuilder.Player{NickName: name, Score: score}
	}
	return nil
}

func (m *mockIntegrationRepository) GetAll() []teambuilder.Player {
	return []teambuilder.Player{}
}

func (m *mockIntegrationRepository) GetTop(n int) []teambuilder.Player {
	return []teambuilder.Player{}
}

func TestIntegration_TwoTeams_FullFlow(t *testing.T) {
	// Setup
	repo := &mockIntegrationRepository{}
	builder := teambuilder.NewTeamBuilder(repo)
	formatter := telegram.NewTeamTableFormatter()

	config := &teambuilder.TeamConfiguration{
		Players: teambuilder.Team{
			{NickName: "Alice", Score: 1000},
			{NickName: "Bob", Score: 2000},
			{NickName: "Charlie", Score: 1500},
			{NickName: "Dave", Score: 1800},
		},
		Constraints: teambuilder.Constraints{},
		NumTeams:    2,
	}

	// Execute
	teams := builder.BuildMultiple(config)

	// Verify teams were created
	assert.Len(t, teams, 2, "Should create 2 teams")

	// Create table
	table := teamtable.NewTeamTableMultiple(teams, "")

	// Verify table structure
	assert.Len(t, table.Headers, 2, "Should have 2 headers")
	assert.Equal(t, "Team 1", table.Headers[0])
	assert.Equal(t, "Team 2", table.Headers[1])

	// Format the table
	formatted := formatter.Format(table)

	// Verify formatting
	assert.Contains(t, formatted, "Team 1", "Should contain Team 1 header")
	assert.Contains(t, formatted, "Team 2", "Should contain Team 2 header")
	assert.Contains(t, formatted, "TS:", "Should contain team scores")
	assert.Contains(t, formatted, "Diff:", "Should contain diff for 2 teams")
	assert.Contains(t, formatted, "начинает за", "Should contain side suggestions for 2 teams")
	assert.Contains(t, formatted, "цель 'Сорян, Братан'", "Should contain footer notes for 2 teams")

	// Verify all players are present
	for _, team := range teams {
		for _, player := range team {
			assert.Contains(t, formatted, player.NickName, "Player %s should be in formatted output", player.NickName)
		}
	}
}

func TestIntegration_FourTeams_FullFlow(t *testing.T) {
	// Setup
	repo := &mockIntegrationRepository{}
	builder := teambuilder.NewTeamBuilder(repo)
	formatter := telegram.NewTeamTableFormatter()

	config := &teambuilder.TeamConfiguration{
		Players: teambuilder.Team{
			{NickName: "Alice", Score: 1000},
			{NickName: "Bob", Score: 2000},
			{NickName: "Charlie", Score: 1500},
			{NickName: "Dave", Score: 1800},
			{NickName: "Eve", Score: 1600},
			{NickName: "Frank", Score: 1700},
			{NickName: "Grace", Score: 1400},
			{NickName: "Henry", Score: 1900},
		},
		Constraints: teambuilder.Constraints{},
		NumTeams:    4,
	}

	// Execute
	teams := builder.BuildMultiple(config)

	// Verify teams were created
	assert.Len(t, teams, 4, "Should create 4 teams")

	// Verify all players are distributed
	totalPlayers := 0
	for _, team := range teams {
		totalPlayers += len(team)
	}
	assert.Equal(t, len(config.Players), totalPlayers, "All players should be distributed")

	// Create table
	table := teamtable.NewTeamTableMultiple(teams, "SorryBroPlayer")

	// Verify table structure
	assert.Len(t, table.Headers, 4, "Should have 4 headers")
	for i := 1; i <= 4; i++ {
		assert.Contains(t, table.Headers, "Team "+string(rune('0'+i)), "Should have Team %d header", i)
	}

	// Format the table
	formatted := formatter.Format(table)

	// Verify formatting
	assert.Contains(t, formatted, "Team 1", "Should contain Team 1 header")
	assert.Contains(t, formatted, "Team 2", "Should contain Team 2 header")
	assert.Contains(t, formatted, "Team 3", "Should contain Team 3 header")
	assert.Contains(t, formatted, "Team 4", "Should contain Team 4 header")
	assert.Contains(t, formatted, "TS:", "Should contain team scores")

	// Verify 4-team specific formatting (what should NOT be present)
	assert.NotContains(t, formatted, "Diff:", "Should NOT contain diff for 4 teams")
	assert.NotContains(t, formatted, "начинает за", "Should NOT contain side suggestions for 4 teams")
	assert.NotContains(t, formatted, "цель 'Сорян, Братан'", "Should NOT contain footer notes for 4 teams")
	assert.NotContains(t, formatted, "(✱)", "Should NOT contain SorryBro marker for 4 teams")

	// Verify all players are present
	for _, team := range teams {
		for _, player := range team {
			assert.Contains(t, formatted, player.NickName, "Player %s should be in formatted output", player.NickName)
		}
	}
}

func TestIntegration_FourTeams_WithConstraints(t *testing.T) {
	// Setup
	repo := &mockIntegrationRepository{}
	builder := teambuilder.NewTeamBuilder(repo)

	config := &teambuilder.TeamConfiguration{
		Players: teambuilder.Team{
			{NickName: "Alice", Score: 1000},
			{NickName: "Bob", Score: 2000},
			{NickName: "Charlie", Score: 1500},
			{NickName: "Dave", Score: 1800},
			{NickName: "Eve", Score: 1600},
			{NickName: "Frank", Score: 1700},
			{NickName: "Grace", Score: 1400},
			{NickName: "Henry", Score: 1900},
		},
		Constraints: teambuilder.Constraints{
			{Type: teambuilder.ConstraintTogether, Player1: "Alice", Player2: "Bob"},
			{Type: teambuilder.ConstraintSeparate, Player1: "Charlie", Player2: "Dave"},
		},
		NumTeams: 4,
	}

	// Execute
	teams := builder.BuildMultiple(config)

	// Verify constraints are satisfied
	var aliceTeam, bobTeam, charlieTeam, daveTeam int = -1, -1, -1, -1

	for i, team := range teams {
		for _, player := range team {
			switch player.NickName {
			case "Alice":
				aliceTeam = i
			case "Bob":
				bobTeam = i
			case "Charlie":
				charlieTeam = i
			case "Dave":
				daveTeam = i
			}
		}
	}

	// Alice and Bob should be together
	assert.Equal(t, aliceTeam, bobTeam, "Alice and Bob should be in the same team")

	// Charlie and Dave should be separate
	assert.NotEqual(t, charlieTeam, daveTeam, "Charlie and Dave should be in different teams")
}

func TestIntegration_DefaultToTwoTeams(t *testing.T) {
	// Setup
	repo := &mockIntegrationRepository{}
	builder := teambuilder.NewTeamBuilder(repo)
	formatter := telegram.NewTeamTableFormatter()

	config := &teambuilder.TeamConfiguration{
		Players: teambuilder.Team{
			{NickName: "Alice", Score: 1000},
			{NickName: "Bob", Score: 2000},
			{NickName: "Charlie", Score: 1500},
			{NickName: "Dave", Score: 1800},
		},
		Constraints: teambuilder.Constraints{},
		NumTeams:    0, // Invalid - should default to 2
	}

	// Execute
	teams := builder.BuildMultiple(config)

	// Verify default behavior
	assert.Len(t, teams, 2, "Should default to 2 teams when NumTeams is invalid")

	// Create table
	table := teamtable.NewTeamTableMultiple(teams, "")

	// Format
	formatted := formatter.Format(table)

	// Should have 2-team formatting
	assert.Contains(t, formatted, "Diff:", "Should contain diff for 2 teams (default)")
}

func TestIntegration_SeparatorAlignment(t *testing.T) {
	// Setup
	repo := &mockIntegrationRepository{}
	builder := teambuilder.NewTeamBuilder(repo)
	formatter := telegram.NewTeamTableFormatter()

	config := &teambuilder.TeamConfiguration{
		Players: teambuilder.Team{
			{NickName: "A", Score: 1000},
			{NickName: "B", Score: 2000},
			{NickName: "C", Score: 1500},
			{NickName: "D", Score: 1800},
			{NickName: "E", Score: 1600},
			{NickName: "F", Score: 1700},
			{NickName: "G", Score: 1400},
			{NickName: "H", Score: 1900},
		},
		Constraints: teambuilder.Constraints{},
		NumTeams:    4,
	}

	// Execute
	teams := builder.BuildMultiple(config)
	table := teamtable.NewTeamTableMultiple(teams, "")
	formatted := formatter.Format(table)

	// Verify separator structure
	lines := strings.Split(formatted, "\n")
	var separatorLines []string

	for _, line := range lines {
		if strings.HasPrefix(line, "|") && strings.Contains(line, "---") {
			separatorLines = append(separatorLines, line)
		}
	}

	// All separators should be identical and properly formatted
	assert.Greater(t, len(separatorLines), 0, "Should have separator lines")

	if len(separatorLines) > 1 {
		firstSep := separatorLines[0]
		for i := 1; i < len(separatorLines); i++ {
			assert.Equal(t, firstSep, separatorLines[i], "All separators should be identical")
		}
	}

	// Separators should have proper column structure (not one long line)
	for _, sep := range separatorLines {
		pipeCount := strings.Count(sep, "|")
		assert.GreaterOrEqual(t, pipeCount, 5, "Separators should have proper column delimiters")
	}
}
