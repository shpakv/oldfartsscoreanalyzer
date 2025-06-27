package telegram

import (
	"fmt"
	"math"
	"oldfartscounter/internal/teamtable"
	"strings"
)

// Set a maximum length for strings (e.g., 15 characters)
const maxStringLength = 15

type TeamTableFormatter struct {
}

func NewTeamTableFormatter() *TeamTableFormatter {
	return &TeamTableFormatter{}
}

func (f *TeamTableFormatter) Format(table *teamtable.TeamTable) string {
	// Find the longest string in each column
	colWidths := make([]int, len(table.Headers))
	for i, header := range table.Headers {
		colWidths[i] = len(header)
	}
	for _, row := range table.Rows {
		for j, cell := range row {
			truncated := truncateWithEllipsis(cell, maxStringLength)
			if len(truncated) > colWidths[j] {
				colWidths[j] = len(truncated)
			}
		}
	}
	for i, score := range table.TeamScore {
		tsString := "TS: " + score
		if len(tsString) > colWidths[i] {
			colWidths[i] = len(tsString)
		}
	}

	// Calculate total table width for single-row entries like Diff
	totalWidth := 1 // Account for the initial "|"
	for _, width := range colWidths {
		totalWidth += width + 3 // Column width + padding and borders
	}
	totalWidth -= 1 // Remove the extra space after the last column

	// Create a formatted table
	var sb strings.Builder

	// Open telegram formatting
	sb.WriteString("```\n")

	// Format headers
	sb.WriteString("| ")
	for i, header := range table.Headers {
		sb.WriteString(fmt.Sprintf("%-*s | ", colWidths[i], header))
	}
	sb.WriteString("\n")

	// Add separator
	sb.WriteString("|-")
	for _, width := range colWidths {
		sb.WriteString(strings.Repeat("-", width+2))
	}
	sb.WriteString("|\n")

	// Format rows
	for _, row := range table.Rows {
		sb.WriteString("| ")
		for j, cell := range row {
			truncated := truncateWithEllipsis(cell, maxStringLength)
			sb.WriteString(fmt.Sprintf("%-*s | ", colWidths[j], truncated))
		}
		sb.WriteString("\n")
	}

	// Add total scores
	sb.WriteString("|-")
	for _, width := range colWidths {
		sb.WriteString(strings.Repeat("-", width+2))
	}
	sb.WriteString("|\n")
	sb.WriteString(fmt.Sprintf("| %-*s | %-*s |\n", colWidths[0], "TS: "+table.TeamScore[0], colWidths[1], "TS: "+table.TeamScore[1]))

	// Calculate percentage difference
	percentDiff := calculatePercentageDifference(table.TeamScore[0], table.TeamScore[1])

	// Add final score difference as a single-row entry with percentage
	diffText := fmt.Sprintf("Diff: %s (%s%%)", table.ScoreDifference, percentDiff)
	sb.WriteString(fmt.Sprintf("| %-*s |\n", totalWidth-3, diffText))

	sb.WriteString("|-")
	for _, width := range colWidths {
		sb.WriteString(strings.Repeat("-", width+2))
	}
	sb.WriteString("|\n")
	t1Side := "Team 1 начинает за CT"
	t2Side := "Team 2 начинает за T"
	if table.TeamScore[0] > table.TeamScore[1] {
		t1Side = "Team 1 начинает за T"
		t2Side = "Team 2 начинает за CT"
	}
	sb.WriteString(fmt.Sprintf("| %-*s |\n", totalWidth-3, t1Side))
	sb.WriteString(fmt.Sprintf("| %-*s |\n", totalWidth-3, t2Side))
	sb.WriteString("|-")
	for _, width := range colWidths {
		sb.WriteString(strings.Repeat("-", width+2))
	}
	sb.WriteString("|\n")

	// Close telegram formatting
	sb.WriteString("```\n")

	return sb.String()
}

// calculatePercentageDifference calculates the percentage difference between two score strings
func calculatePercentageDifference(score1, score2 string) string {
	// Parse scores to float
	var s1, s2 float64
	fmt.Sscanf(score1, "%f", &s1)
	fmt.Sscanf(score2, "%f", &s2)

	// Avoid division by zero
	if s1 == 0 && s2 == 0 {
		return "0.0"
	}

	// Calculate average score
	avgScore := (s1 + s2) / 2

	// Calculate absolute difference
	absDiff := math.Abs(s1 - s2)

	// Calculate percentage difference relative to average
	percentDiff := (absDiff / avgScore) * 100

	// Format to one decimal place
	return fmt.Sprintf("%.1f", percentDiff)
}

func truncateWithEllipsis(s string, maxLen int) string {
	if len(s) > maxLen {
		if maxLen > 3 {
			return s[:maxLen-3] + "..."
		}
		return s[:maxLen] // Если maxLen <= 3, возвращаем обрезанное без ...
	}
	return s
}
