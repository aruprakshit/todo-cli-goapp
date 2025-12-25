package main

import (
	"database/sql"
	"strings"
	"testing"
	"time"
)

func TestColorize(t *testing.T) {
	tests := []struct {
		name     string
		color    Color
		text     string
		expected string
	}{
		{
			name:     "red color",
			color:    Red,
			text:     "error",
			expected: string(Red) + "error" + string(Reset),
		},
		{
			name:     "green color",
			color:    Green,
			text:     "success",
			expected: string(Green) + "success" + string(Reset),
		},
		{
			name:     "empty text",
			color:    Blue,
			text:     "",
			expected: string(Blue) + "" + string(Reset),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := colorize(tt.color, tt.text)
			if result != tt.expected {
				t.Errorf("colorize(%q, %q) = %q, want %q", tt.color, tt.text, result, tt.expected)
			}
		})
	}
}

func TestPriorityColor(t *testing.T) {
	tests := []struct {
		name     string
		priority Priority
		expected Color
	}{
		{
			name:     "high priority returns red",
			priority: PriorityHigh,
			expected: Red,
		},
		{
			name:     "medium priority returns yellow",
			priority: PriorityMedium,
			expected: Yellow,
		},
		{
			name:     "low priority returns green",
			priority: PriorityLow,
			expected: Green,
		},
		{
			name:     "invalid priority returns reset",
			priority: Priority("invalid"),
			expected: Reset,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := priorityColor(tt.priority)
			if result != tt.expected {
				t.Errorf("priorityColor(%q) = %q, want %q", tt.priority, result, tt.expected)
			}
		})
	}
}

func TestFormatDueDate(t *testing.T) {
	today := time.Now().Truncate(24 * time.Hour)

	tests := []struct {
		name     string
		dueDate  sql.NullTime
		contains string // substring to check for
		color    Color  // expected color wrapper
	}{
		{
			name:     "null date returns empty string",
			dueDate:  sql.NullTime{Valid: false},
			contains: "",
			color:    "",
		},
		{
			name:     "overdue date",
			dueDate:  sql.NullTime{Time: today.AddDate(0, 0, -1), Valid: true},
			contains: "(OVERDUE)",
			color:    Red,
		},
		{
			name:     "due today",
			dueDate:  sql.NullTime{Time: today, Valid: true},
			contains: "(TODAY)",
			color:    Red,
		},
		{
			name:     "due tomorrow",
			dueDate:  sql.NullTime{Time: today.AddDate(0, 0, 1), Valid: true},
			contains: "(tomorrow)",
			color:    Yellow,
		},
		{
			name:     "due in 3 days",
			dueDate:  sql.NullTime{Time: today.AddDate(0, 0, 3), Valid: true},
			contains: today.AddDate(0, 0, 3).Format("2006-01-02"),
			color:    Yellow,
		},
		{
			name:     "due in 5 days",
			dueDate:  sql.NullTime{Time: today.AddDate(0, 0, 5), Valid: true},
			contains: today.AddDate(0, 0, 5).Format("2006-01-02"),
			color:    Green,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDueDate(tt.dueDate)

			// Check for null case
			if !tt.dueDate.Valid {
				if result != "" {
					t.Errorf("formatDueDate() = %q, want empty string", result)
				}
				return
			}

			// Check that result contains expected substring
			if tt.contains != "" && !strings.Contains(result, tt.contains) {
				t.Errorf("formatDueDate() = %q, want it to contain %q", result, tt.contains)
			}

			// Check that result starts with expected color
			if tt.color != "" && !strings.Contains(result, string(tt.color)) {
				t.Errorf("formatDueDate() = %q, want it to have color %q", result, tt.color)
			}
		})
	}
}
