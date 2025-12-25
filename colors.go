package main

import (
	"database/sql"
	"time"
)

func formatDueDate(dueDate sql.NullTime) string {
	if !dueDate.Valid {
		return ""
	}

	due := dueDate.Time
	today := time.Now().Truncate(24 * time.Hour)
	dueDay := due.Truncate(24 * time.Hour)

	daysUntil := int(dueDay.Sub(today).Hours() / 24)

	dateStr := due.Format("2006-01-02")

	if daysUntil < 0 {
		return colorize(Red, dateStr+" (OVERDUE)")
	} else if daysUntil == 0 {
		return colorize(Red, dateStr+" (TODAY)")
	} else if daysUntil == 1 {
		return colorize(Yellow, dateStr+" (tomorrow)")
	} else if daysUntil <= 3 {
		return colorize(Yellow, dateStr)
	}
	return colorize(Green, dateStr)
}

// Color represents an ANSI color code
type Color string

// ANSI color codes
const (
	Reset  Color = "\033[0m"
	Red    Color = "\033[31m"
	Green  Color = "\033[32m"
	Yellow Color = "\033[33m"
	Blue   Color = "\033[34m"
	Purple Color = "\033[35m"
	Cyan   Color = "\033[36m"
	Gray   Color = "\033[90m"
	Bold   Color = "\033[1m"
)

func colorize(color Color, text string) string {
	return string(color) + text + string(Reset)
}

func priorityColor(p Priority) Color {
	switch p {
	case PriorityHigh:
		return Red
	case PriorityMedium:
		return Yellow
	case PriorityLow:
		return Green
	default:
		return Reset
	}
}
