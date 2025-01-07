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
	"|             | IncrediblyLo... | \n" +
	"|-------------------------------|\n" +
	"| TS: 8000.00 | TS: 7000.00     |\n" +
	"| Diff: 1000.00                 |\n" +
	"```\n"

func TestFormatter_Format(t *testing.T) {
	table := &teamtable.TeamTable{
		Headers: []string{"Team 1", "Team 2"},
		Rows: [][]string{
			{"Player1", "Player3"},
			{"Player2", "Player4"},
			{"", "IncrediblyLongNickNameOfPlayer5"},
		},
		TeamScore:       []string{"8000.00", "7000.00"},
		ScoreDifference: "1000.00",
	}

	f := &Formatter{}
	formatted := f.Format(table)
	fmt.Println(formatted)

	assert.Equal(t, ExpectedOutput, formatted, "Formatted table mismatch")
}
