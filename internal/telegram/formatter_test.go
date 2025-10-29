package telegram

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"oldfartscounter/internal/teamtable"
)

const ExpectedOutput = "" +
	"```\n" +
	"| Team 1      | Team 2          | \n" +
	"|-------------|-----------------|\n" +
	"| Player1     | Player3         | \n" +
	"| Player2     | Player4         | \n" +
	"|             | (✱) Mr. Tits... | \n" +
	"|-------------|-----------------|\n" +
	"| TS: 8000.00 | TS: 7000.00     | \n" +
	"| Diff: 1000.00 (13.3%)         |\n" +
	"|-------------|-----------------|\n" +
	"| Team 1 начинает за T          |\n" +
	"| Team 2 начинает за CT         |\n" +
	"|-------------|-----------------|\n" +
	"| (✱) - цель 'Сорян, Братан'    |\n" +
	"| Приз - Phoenix Case $8        |\n" +
	"|-------------|-----------------|\n" +
	"```\n"

func TestFormatter_Format_TwoTeams(t *testing.T) {
	table := &teamtable.TeamTable{
		Headers: []string{"Team 1", "Team 2"},
		Rows: [][]string{
			{"Player1", "Player3"},
			{"Player2", "Player4"},
			{"", "Mr. Titspervert"},
		},
		TeamScore:       []string{"8000.00", "7000.00"},
		ScoreDifference: "1000.00",
		SorryBro:        "Mr. Titspervert",
	}

	f := &TeamTableFormatter{}
	formatted := f.Format(table)
	fmt.Println(formatted)

	assert.Equal(t, ExpectedOutput, formatted, "Formatted table mismatch")
}

func TestFormatter_Format_FourTeams(t *testing.T) {
	table := &teamtable.TeamTable{
		Headers: []string{"Team 1", "Team 2", "Team 3", "Team 4"},
		Rows: [][]string{
			{"Player1", "Player3", "Player5", "Player7"},
			{"Player2", "Player4", "Player6", "Player8"},
		},
		TeamScore:       []string{"300.00", "200.00", "200.00", "200.00"},
		ScoreDifference: "100.00",
		SorryBro:        "Player1", // Should NOT be shown for 4 teams
	}

	f := &TeamTableFormatter{}
	formatted := f.Format(table)
	fmt.Println(formatted)

	// Verify that SorryBro marker is NOT present
	assert.NotContains(t, formatted, "(✱)", "SorryBro marker should not be shown for 4 teams")

	// Verify that Diff line is NOT present
	assert.NotContains(t, formatted, "Diff:", "Diff line should not be shown for 4 teams")

	// Verify that footer notes are NOT present
	assert.NotContains(t, formatted, "цель 'Сорян, Братан'", "Footer notes should not be shown for 4 teams")
	assert.NotContains(t, formatted, "Phoenix Case", "Prize note should not be shown for 4 teams")

	// Verify that side suggestions are NOT present
	assert.NotContains(t, formatted, "начинает за CT", "Side suggestions should not be shown for 4 teams")
	assert.NotContains(t, formatted, "начинает за T", "Side suggestions should not be shown for 4 teams")

	// Verify that all team headers are present
	assert.Contains(t, formatted, "Team 1", "Team 1 header should be present")
	assert.Contains(t, formatted, "Team 2", "Team 2 header should be present")
	assert.Contains(t, formatted, "Team 3", "Team 3 header should be present")
	assert.Contains(t, formatted, "Team 4", "Team 4 header should be present")

	// Verify that all team scores are present
	assert.Contains(t, formatted, "TS: 300.00", "Team 1 score should be present")
	assert.Contains(t, formatted, "TS: 200.00", "Team 2/3/4 scores should be present")

	// Verify that all players are present
	for i := 1; i <= 8; i++ {
		assert.Contains(t, formatted, fmt.Sprintf("Player%d", i), "Player%d should be present", i)
	}
}

func TestFormatter_Format_FourTeams_UnevenSizes(t *testing.T) {
	table := &teamtable.TeamTable{
		Headers: []string{"Team 1", "Team 2", "Team 3", "Team 4"},
		Rows: [][]string{
			{"Player1", "Player4", "Player5", "Player7"},
			{"Player2", "", "Player6", "Player8"},
			{"Player3", "", "", "Player9"},
		},
		TeamScore:       []string{"450.00", "50.00", "200.00", "300.00"},
		ScoreDifference: "400.00",
		SorryBro:        "",
	}

	f := &TeamTableFormatter{}
	formatted := f.Format(table)

	// Verify structure is correct with empty cells
	assert.Contains(t, formatted, "Player1", "Player1 should be present")
	assert.Contains(t, formatted, "Player4", "Player4 should be present")
	assert.Contains(t, formatted, "Player9", "Player9 should be present")

	// Verify no SorryBro elements for 4 teams
	assert.NotContains(t, formatted, "(✱)", "SorryBro marker should not be shown for 4 teams")
	assert.NotContains(t, formatted, "Diff:", "Diff line should not be shown for 4 teams")
}

func TestFormatter_Format_TwoTeams_NoSorryBro(t *testing.T) {
	table := &teamtable.TeamTable{
		Headers: []string{"Team 1", "Team 2"},
		Rows: [][]string{
			{"Player1", "Player3"},
			{"Player2", "Player4"},
		},
		TeamScore:       []string{"5000.00", "5000.00"},
		ScoreDifference: "0.00",
		SorryBro:        "", // No SorryBro player
	}

	f := &TeamTableFormatter{}
	formatted := f.Format(table)

	// Should still have all 2-team elements except SorryBro marker
	assert.Contains(t, formatted, "Diff:", "Diff line should be present for 2 teams")
	assert.Contains(t, formatted, "начинает за", "Side suggestions should be present for 2 teams")
	assert.Contains(t, formatted, "цель 'Сорян, Братан'", "Footer notes should be present for 2 teams")

	// But no actual SorryBro marker since SorryBro is empty
	// The skull should only appear next to a player name
	lines := strings.Split(formatted, "\n")
	skullCount := 0
	for _, line := range lines {
		if strings.Contains(line, "(✱)") && !strings.Contains(line, "цель") {
			skullCount++
		}
	}
	assert.Equal(t, 0, skullCount, "No skull marker should appear when SorryBro is empty")
}

func TestFormatter_SeparatorAlignment_FourTeams(t *testing.T) {
	table := &teamtable.TeamTable{
		Headers: []string{"Team 1", "Team 2", "Team 3", "Team 4"},
		Rows: [][]string{
			{"P1", "P3", "P5", "P7"},
			{"P2", "P4", "P6", "P8"},
		},
		TeamScore:       []string{"100.00", "200.00", "150.00", "175.00"},
		ScoreDifference: "100.00",
		SorryBro:        "",
	}

	f := &TeamTableFormatter{}
	formatted := f.Format(table)

	// Parse the formatted output to verify separator alignment
	lines := strings.Split(formatted, "\n")

	var separatorLines []string

	for _, line := range lines {
		if strings.HasPrefix(line, "|") && strings.Contains(line, "---") {
			separatorLines = append(separatorLines, line)
		}
	}

	// Verify we have separator lines
	assert.Greater(t, len(separatorLines), 0, "Should have at least one separator line")

	// Separators should have | between columns, not just one long line
	for i, sepLine := range separatorLines {
		pipeCount := strings.Count(sepLine, "|")
		// Should have: | at start + | between each column (3 for 4 teams) + | at end = 5 total
		assert.GreaterOrEqual(t, pipeCount, 5,
			"Separator line %d should have proper column separators", i)
	}

	// All separators should have the same structure
	if len(separatorLines) > 1 {
		firstSep := separatorLines[0]
		for i := 1; i < len(separatorLines); i++ {
			assert.Equal(t, firstSep, separatorLines[i],
				"All separator lines should be identical")
		}
	}
}
