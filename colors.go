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

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[90m"
	Bold   = "\033[1m"
)

func colorize(color, text string) string {
	return color + text + Reset
}

func priorityColor(p Priority) string {
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
