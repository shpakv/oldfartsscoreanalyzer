package telegram

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"oldfartscounter/internal/teamtable"
	"testing"
)

const ExpectedOutput = "" +
	"```\n" +
	"| Team 1      | Team 2          | \n" +
	"|-------------------------------|\n" +
	"| Player1     | Player3         | \n" +
	"| Player2     | Player4         | \n" +
	"|             | (✱) Mr. Tits... | \n" +
	"|-------------------------------|\n" +
	"| TS: 8000.00 | TS: 7000.00     |\n" +
	"| Diff: 1000.00 (13.3%)         |\n" +
	"|-------------------------------|\n" +
	"| Team 1 начинает за T          |\n" +
	"| Team 2 начинает за CT         |\n" +
	"|-------------------------------|\n" +
	"| (✱) - цель 'Сорян, Братан'    |\n" +
	"| Приз - Phoenix Case $8        |\n" +
	"|-------------------------------|\n" +
	"```\n"

func TestFormatter_Format(t *testing.T) {
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
