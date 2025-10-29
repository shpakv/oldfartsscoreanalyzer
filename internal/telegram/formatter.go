package telegram

import (
	"fmt"
	"math"
	"oldfartscounter/internal/teamtable"
	"strings"

	"github.com/mattn/go-runewidth"
)

// Set a maximum length for strings (visible columns, not bytes)
const maxStringLength = 15

const skull = "(✱)" // маркер цели

type TeamTableFormatter struct{}

func NewTeamTableFormatter() *TeamTableFormatter { return &TeamTableFormatter{} }

func (f *TeamTableFormatter) Format(table *teamtable.TeamTable) string {
	// Find the longest visible width in each column
	colWidths := make([]int, len(table.Headers))
	for i, header := range table.Headers {
		colWidths[i] = runewidth.StringWidth(header)
	}
	for _, row := range table.Rows {
		for j, cell := range row {
			var displayCell string
			// Account for skull width when calculating column widths (only for 2 teams)
			if len(table.Headers) == 2 && cell == table.SorryBro && cell != "" {
				withSkull := skull + " " + cell
				skullWidth := runewidth.StringWidth(skull)
				if runewidth.StringWidth(withSkull) > maxStringLength {
					// Reserve space for skull + space (1 char) + "..." (3 chars)
					truncatedName := runewidth.Truncate(cell, maxStringLength-skullWidth-4, "")
					displayCell = skull + " " + truncatedName + "..."
				} else {
					displayCell = withSkull
				}
			} else {
				displayCell = truncateVisible(cell, maxStringLength)
			}
			if w := runewidth.StringWidth(displayCell); w > colWidths[j] {
				colWidths[j] = w
			}
		}
	}
	for i, score := range table.TeamScore {
		tsString := "TS: " + score
		if w := runewidth.StringWidth(tsString); w > colWidths[i] {
			colWidths[i] = w
		}
	}

	// Calculate total table width for single-row entries like Diff
	totalWidth := 1 // initial "|"
	for _, width := range colWidths {
		totalWidth += width + 3 // width + " | "
	}
	totalWidth -= 1 // remove extra space after last column

	var sb strings.Builder
	sb.WriteString("```\n")

	// Headers
	sb.WriteString("| ")
	for i, header := range table.Headers {
		sb.WriteString(padRight(header, colWidths[i]))
		sb.WriteString(" | ")
	}
	sb.WriteString("\n")

	// Separator
	sb.WriteString("|")
	for i, width := range colWidths {
		sb.WriteString(strings.Repeat("-", width+2))
		if i < len(colWidths)-1 {
			sb.WriteString("|")
		}
	}
	sb.WriteString("|\n")

	// Rows
	for _, row := range table.Rows {
		sb.WriteString("| ")
		for j, cell := range row {
			// Add skull emoji if this is the SorryBro player (only for 2 teams)
			displayCell := cell
			if len(table.Headers) == 2 && cell == table.SorryBro && cell != "" {
				// Check if we need to truncate
				withSkull := skull + " " + cell
				skullWidth := runewidth.StringWidth(skull)
				if runewidth.StringWidth(withSkull) > maxStringLength {
					// Truncate the name, add skull before the name
					truncatedName := runewidth.Truncate(cell, maxStringLength-skullWidth-4, "")
					displayCell = skull + " " + truncatedName + "..."
				} else {
					displayCell = withSkull
				}
			} else {
				displayCell = truncateVisible(cell, maxStringLength)
			}
			sb.WriteString(padRight(displayCell, colWidths[j]))
			sb.WriteString(" | ")
		}
		sb.WriteString("\n")
	}

	// Totals
	sb.WriteString("|")
	for i, width := range colWidths {
		sb.WriteString(strings.Repeat("-", width+2))
		if i < len(colWidths)-1 {
			sb.WriteString("|")
		}
	}
	sb.WriteString("|\n")

	// Format team scores row
	sb.WriteString("| ")
	for i := 0; i < len(table.TeamScore); i++ {
		tsString := "TS: " + table.TeamScore[i]
		sb.WriteString(padRight(tsString, colWidths[i]))
		sb.WriteString(" | ")
	}
	sb.WriteString("\n")

	// Additional info (only for 2 teams)
	if len(table.TeamScore) == 2 {
		// Percentage diff
		percentDiff := calculatePercentageDifferenceMultiple(table.TeamScore)
		diffText := fmt.Sprintf("Diff: %s (%s%%)", table.ScoreDifference, percentDiff)
		sb.WriteString(fmt.Sprintf("| %s |\n", padRight(diffText, totalWidth-3)))

		// Side suggestion
		sb.WriteString("|")
		for i, width := range colWidths {
			sb.WriteString(strings.Repeat("-", width+2))
			if i < len(colWidths)-1 {
				sb.WriteString("|")
			}
		}
		sb.WriteString("|\n")
		t1Side := "Team 1 начинает за CT"
		t2Side := "Team 2 начинает за T"
		if table.TeamScore[0] > table.TeamScore[1] {
			t1Side = "Team 1 начинает за T"
			t2Side = "Team 2 начинает за CT"
		}
		sb.WriteString(fmt.Sprintf("| %s |\n", padRight(t1Side, totalWidth-3)))
		sb.WriteString(fmt.Sprintf("| %s |\n", padRight(t2Side, totalWidth-3)))
	}

	// Footer notes (only for 2 teams)
	if len(table.TeamScore) == 2 {
		sb.WriteString("|")
		for i, width := range colWidths {
			sb.WriteString(strings.Repeat("-", width+2))
			if i < len(colWidths)-1 {
				sb.WriteString("|")
			}
		}
		sb.WriteString("|\n")
		skullText := `(✱) - цель 'Сорян, Братан'`
		giftText := `Приз - Phoenix Case $8`
		sb.WriteString(fmt.Sprintf("| %s |\n", padRight(skullText, totalWidth-3)))
		sb.WriteString(fmt.Sprintf("| %s |\n", padRight(giftText, totalWidth-3)))

		sb.WriteString("|")
		for i, width := range colWidths {
			sb.WriteString(strings.Repeat("-", width+2))
			if i < len(colWidths)-1 {
				sb.WriteString("|")
			}
		}
		sb.WriteString("|\n")
	} else {
		// For 4+ teams, just add final separator
		sb.WriteString("|")
		for i, width := range colWidths {
			sb.WriteString(strings.Repeat("-", width+2))
			if i < len(colWidths)-1 {
				sb.WriteString("|")
			}
		}
		sb.WriteString("|\n")
	}

	sb.WriteString("```\n")
	return sb.String()
}

func calculatePercentageDifference(score1, score2 string) string {
	var s1, s2 float64
	fmt.Sscanf(score1, "%f", &s1)
	fmt.Sscanf(score2, "%f", &s2)
	if s1 == 0 && s2 == 0 {
		return "0.0"
	}
	avg := (s1 + s2) / 2
	diff := math.Abs(s1 - s2)
	return fmt.Sprintf("%.1f", (diff/avg)*100)
}

func calculatePercentageDifferenceMultiple(scores []string) string {
	if len(scores) == 0 {
		return "0.0"
	}

	// Parse all scores
	floatScores := make([]float64, len(scores))
	sum := 0.0
	for i, scoreStr := range scores {
		fmt.Sscanf(scoreStr, "%f", &floatScores[i])
		sum += floatScores[i]
	}

	if sum == 0 {
		return "0.0"
	}

	// Find min and max
	minScore := floatScores[0]
	maxScore := floatScores[0]
	for _, score := range floatScores {
		if score < minScore {
			minScore = score
		}
		if score > maxScore {
			maxScore = score
		}
	}

	// Calculate average and percentage difference
	avg := sum / float64(len(scores))
	diff := math.Abs(maxScore - minScore)
	return fmt.Sprintf("%.1f", (diff/avg)*100)
}

func truncateVisible(s string, maxLen int) string {
	if runewidth.StringWidth(s) <= maxLen {
		return s
	}
	return runewidth.Truncate(s, maxLen, "...")
}

func padRight(s string, width int) string {
	return runewidth.FillRight(s, width)
}
