package main

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

func priorityColor(pValue string) string {
	switch pValue {
	case "high":
		return Red
	case "medium":
		return Yellow
	case "low":
		return Green
	default:
		return Reset
	}
}
