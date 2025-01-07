package telegram

import (
	"fmt"
	"oldfartscounter/internal/teamtable"
	"strings"
)

// Set a maximum length for strings (e.g., 15 characters)
const maxStringLength = 15

type Formatter struct {
}

func NewFormatter() *Formatter {
	return &Formatter{}
}

func (f *Formatter) Format(table *teamtable.TeamTable) string {
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

	// Add final score difference as a single-row entry
	sb.WriteString(fmt.Sprintf("| %-*s |\n", totalWidth-3, "Diff: "+table.ScoreDifference))

	// Close telegram formatting
	sb.WriteString("```\n")

	return sb.String()
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
