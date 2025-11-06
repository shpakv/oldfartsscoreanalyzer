package teambuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock repository for testing
type mockRepository struct{}

func (m *mockRepository) FindByName(name string) *Player {
	scores := map[string]float64{
		"Player1": 1000,
		"Player2": 2000,
		"Player3": 1500,
		"Player4": 1800,
		"Player5": 1600,
		"Player6": 1700,
		"Player7": 1400,
		"Player8": 1900,
	}
	if score, ok := scores[name]; ok {
		return &Player{NickName: name, Score: score}
	}
	return nil
}

func (m *mockRepository) GetAll() []Player {
	return []Player{}
}

func (m *mockRepository) GetTop(n int) []Player {
	return []Player{}
}

func (m *mockRepository) GetAverageMu() float64 {
	return 3.0 // Mock value for testing
}

func TestBuildMultiple_TwoTeams(t *testing.T) {
	repo := &mockRepository{}
	builder := NewTeamBuilder(repo)

	config := &TeamConfiguration{
		Players: Team{
			{NickName: "Player1", Score: 1000},
			{NickName: "Player2", Score: 2000},
			{NickName: "Player3", Score: 1500},
			{NickName: "Player4", Score: 1800},
		},
		Constraints: Constraints{},
		NumTeams:    2,
	}

	teams := builder.Build(config)

	assert.Len(t, teams, 2, "Should return 2 teams")

	// Check that all players are distributed
	totalPlayers := len(teams[0]) + len(teams[1])
	assert.Equal(t, len(config.Players), totalPlayers, "All players should be distributed")

	// Check that teams are relatively balanced
	diff := teams[0].Score() - teams[1].Score()
	if diff < 0 {
		diff = -diff
	}
	assert.Less(t, diff, 1000.0, "Teams should be relatively balanced")
}

func TestBuildMultiple_FourTeams(t *testing.T) {
	repo := &mockRepository{}
	builder := NewTeamBuilder(repo)

	config := &TeamConfiguration{
		Players: Team{
			{NickName: "Player1", Score: 1000},
			{NickName: "Player2", Score: 2000},
			{NickName: "Player3", Score: 1500},
			{NickName: "Player4", Score: 1800},
			{NickName: "Player5", Score: 1600},
			{NickName: "Player6", Score: 1700},
			{NickName: "Player7", Score: 1400},
			{NickName: "Player8", Score: 1900},
		},
		Constraints: Constraints{},
		NumTeams:    4,
	}

	teams := builder.Build(config)

	assert.Len(t, teams, 4, "Should return 4 teams")

	// Check that all players are distributed
	totalPlayers := 0
	for _, team := range teams {
		totalPlayers += len(team)
	}
	assert.Equal(t, len(config.Players), totalPlayers, "All players should be distributed")

	// Check that each team has players
	for i, team := range teams {
		assert.Greater(t, len(team), 0, "Team %d should have at least one player", i+1)
	}

	// Check that teams are relatively balanced
	minScore := teams[0].Score()
	maxScore := teams[0].Score()
	for _, team := range teams {
		score := team.Score()
		if score < minScore {
			minScore = score
		}
		if score > maxScore {
			maxScore = score
		}
	}
	diff := maxScore - minScore
	assert.Less(t, diff, 1500.0, "Teams should be relatively balanced")
}

func TestBuildMultiple_DefaultToTwoTeams(t *testing.T) {
	repo := &mockRepository{}
	builder := NewTeamBuilder(repo)

	config := &TeamConfiguration{
		Players: Team{
			{NickName: "Player1", Score: 1000},
			{NickName: "Player2", Score: 2000},
			{NickName: "Player3", Score: 1500},
			{NickName: "Player4", Score: 1800},
		},
		Constraints: Constraints{},
		NumTeams:    0, // Invalid, should default to 2
	}

	teams := builder.Build(config)

	assert.Len(t, teams, 2, "Should default to 2 teams when NumTeams is invalid")
}

func TestBuildMultiple_WithConstraintsTogether(t *testing.T) {
	repo := &mockRepository{}
	builder := NewTeamBuilder(repo)

	config := &TeamConfiguration{
		Players: Team{
			{NickName: "Player1", Score: 1000},
			{NickName: "Player2", Score: 2000},
			{NickName: "Player3", Score: 1500},
			{NickName: "Player4", Score: 1800},
			{NickName: "Player5", Score: 1600},
			{NickName: "Player6", Score: 1700},
			{NickName: "Player7", Score: 1400},
			{NickName: "Player8", Score: 1900},
		},
		Constraints: Constraints{
			{Type: ConstraintTogether, Player1: "Player1", Player2: "Player2"},
		},
		NumTeams: 4,
	}

	teams := builder.Build(config)

	// Find which team has Player1
	player1Team := -1
	player2Team := -1

	for i, team := range teams {
		for _, player := range team {
			if player.NickName == "Player1" {
				player1Team = i
			}
			if player.NickName == "Player2" {
				player2Team = i
			}
		}
	}

	assert.Equal(t, player1Team, player2Team, "Player1 and Player2 should be in the same team")
}

func TestBuildMultiple_WithConstraintsSeparate(t *testing.T) {
	repo := &mockRepository{}
	builder := NewTeamBuilder(repo)

	config := &TeamConfiguration{
		Players: Team{
			{NickName: "Player1", Score: 1000},
			{NickName: "Player2", Score: 2000},
			{NickName: "Player3", Score: 1500},
			{NickName: "Player4", Score: 1800},
			{NickName: "Player5", Score: 1600},
			{NickName: "Player6", Score: 1700},
			{NickName: "Player7", Score: 1400},
			{NickName: "Player8", Score: 1900},
		},
		Constraints: Constraints{
			{Type: ConstraintSeparate, Player1: "Player1", Player2: "Player2"},
		},
		NumTeams: 4,
	}

	teams := builder.Build(config)

	// Find which team has Player1 and Player2
	player1Team := -1
	player2Team := -1

	for i, team := range teams {
		for _, player := range team {
			if player.NickName == "Player1" {
				player1Team = i
			}
			if player.NickName == "Player2" {
				player2Team = i
			}
		}
	}

	assert.NotEqual(t, player1Team, player2Team, "Player1 and Player2 should be in different teams")
}

func TestBuildMultiple_EmptyPlayers(t *testing.T) {
	repo := &mockRepository{}
	builder := NewTeamBuilder(repo)

	config := &TeamConfiguration{
		Players:     Team{},
		Constraints: Constraints{},
		NumTeams:    4,
	}

	teams := builder.Build(config)

	assert.Len(t, teams, 4, "Should return 4 teams")
	for _, team := range teams {
		assert.Empty(t, team, "Each team should be empty")
	}
}

func TestBuildMultiple_LessThanFourPlayers(t *testing.T) {
	repo := &mockRepository{}
	builder := NewTeamBuilder(repo)

	config := &TeamConfiguration{
		Players: Team{
			{NickName: "Player1", Score: 1000},
			{NickName: "Player2", Score: 2000},
		},
		Constraints: Constraints{},
		NumTeams:    4,
	}

	teams := builder.Build(config)

	assert.Len(t, teams, 4, "Should return 4 teams")

	// Count non-empty teams
	nonEmptyTeams := 0
	for _, team := range teams {
		if len(team) > 0 {
			nonEmptyTeams++
		}
	}
	assert.Equal(t, 2, nonEmptyTeams, "Should have 2 non-empty teams")
}

func TestDistributeFourTeamsSnake(t *testing.T) {
	players := Team{
		{NickName: "Player1", Score: 1000},
		{NickName: "Player2", Score: 2000},
		{NickName: "Player3", Score: 1500},
		{NickName: "Player4", Score: 1800},
		{NickName: "Player5", Score: 1600},
		{NickName: "Player6", Score: 1700},
		{NickName: "Player7", Score: 1400},
		{NickName: "Player8", Score: 1900},
	}

	teams := distributeFourTeamsSnake(players)

	assert.Len(t, teams, 4, "Should return 4 teams")

	// Check that all players are distributed
	totalPlayers := 0
	for _, team := range teams {
		totalPlayers += len(team)
	}
	assert.Equal(t, len(players), totalPlayers, "All players should be distributed")
}

func TestDistributeFourTeamsGreedy(t *testing.T) {
	players := Team{
		{NickName: "Player1", Score: 1000},
		{NickName: "Player2", Score: 2000},
		{NickName: "Player3", Score: 1500},
		{NickName: "Player4", Score: 1800},
		{NickName: "Player5", Score: 1600},
		{NickName: "Player6", Score: 1700},
		{NickName: "Player7", Score: 1400},
		{NickName: "Player8", Score: 1900},
	}

	teams := distributeFourTeamsGreedy(players)

	assert.Len(t, teams, 4, "Should return 4 teams")

	// Check that all players are distributed
	totalPlayers := 0
	for _, team := range teams {
		totalPlayers += len(team)
	}
	assert.Equal(t, len(players), totalPlayers, "All players should be distributed")

	// Greedy algorithm should create relatively balanced teams
	scores := make([]float64, 4)
	for i, team := range teams {
		scores[i] = team.Score()
	}

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

	diff := maxScore - minScore
	assert.Less(t, diff, 1500.0, "Teams should be relatively balanced with greedy algorithm")
}

func TestCalculateTeamsDifference(t *testing.T) {
	teams := []Team{
		{{NickName: "P1", Score: 100}, {NickName: "P2", Score: 200}},
		{{NickName: "P3", Score: 150}, {NickName: "P4", Score: 50}},
		{{NickName: "P5", Score: 120}, {NickName: "P6", Score: 80}},
		{{NickName: "P7", Score: 90}, {NickName: "P8", Score: 110}},
	}

	diff := calculateTeamsDifference(teams)

	// Team 1: 300, Team 2: 200, Team 3: 200, Team 4: 200
	// Max: 300, Min: 200, Diff: 100
	assert.Equal(t, 100.0, diff, "Difference should be 100")
}

func TestIsConstraintSatisfiedMultiple(t *testing.T) {
	teams := []Team{
		{{NickName: "Player1", Score: 100}, {NickName: "Player2", Score: 200}},
		{{NickName: "Player3", Score: 150}, {NickName: "Player4", Score: 50}},
		{{NickName: "Player5", Score: 120}, {NickName: "Player6", Score: 80}},
		{{NickName: "Player7", Score: 90}, {NickName: "Player8", Score: 110}},
	}

	// Test ConstraintTogether - Player1 and Player2 are in team 0
	constraintsTogether := Constraints{
		{Type: ConstraintTogether, Player1: "Player1", Player2: "Player2"},
	}
	assert.True(t, isConstraintSatisfiedMultiple(teams, constraintsTogether), "Together constraint should be satisfied")

	// Test ConstraintTogether - Player1 and Player3 are in different teams
	constraintsTogetherFail := Constraints{
		{Type: ConstraintTogether, Player1: "Player1", Player2: "Player3"},
	}
	assert.False(t, isConstraintSatisfiedMultiple(teams, constraintsTogetherFail), "Together constraint should not be satisfied")

	// Test ConstraintSeparate - Player1 and Player3 are in different teams
	constraintsSeparate := Constraints{
		{Type: ConstraintSeparate, Player1: "Player1", Player2: "Player3"},
	}
	assert.True(t, isConstraintSatisfiedMultiple(teams, constraintsSeparate), "Separate constraint should be satisfied")

	// Test ConstraintSeparate - Player1 and Player2 are in same team
	constraintsSeparateFail := Constraints{
		{Type: ConstraintSeparate, Player1: "Player1", Player2: "Player2"},
	}
	assert.False(t, isConstraintSatisfiedMultiple(teams, constraintsSeparateFail), "Separate constraint should not be satisfied")
}

func TestCopyTeams(t *testing.T) {
	original := []Team{
		{{NickName: "P1", Score: 100}},
		{{NickName: "P2", Score: 200}},
	}

	copied := copyTeams(original)

	// Check that it's a deep copy
	assert.Equal(t, original, copied, "Copied teams should be equal")

	// Modify the copy
	copied[0][0].Score = 999

	// Original should remain unchanged
	assert.NotEqual(t, original[0][0].Score, copied[0][0].Score, "Original should not be affected by changes to copy")
}
