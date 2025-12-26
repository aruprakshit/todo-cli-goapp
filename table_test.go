package main

import (
	"testing"
)

func TestNewTable(t *testing.T) {
	tests := []struct {
		name        string
		headers     []string
		wantWidths  []int
		wantHeaders []string
	}{
		{
			name:        "basic headers",
			headers:     []string{"ID", "Title", "Status"},
			wantWidths:  []int{2, 5, 6},
			wantHeaders: []string{"ID", "Title", "Status"},
		},
		{
			name:        "single header",
			headers:     []string{"Name"},
			wantWidths:  []int{4},
			wantHeaders: []string{"Name"},
		},
		{
			name:        "empty headers",
			headers:     []string{},
			wantWidths:  []int{},
			wantHeaders: []string{},
		},
		{
			name:        "unicode headers",
			headers:     []string{"✓", "名前"},
			wantWidths:  []int{1, 2},
			wantHeaders: []string{"✓", "名前"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewTable(tt.headers)

			if len(table.Headers) != len(tt.wantHeaders) {
				t.Errorf("NewTable() headers length = %d, want %d", len(table.Headers), len(tt.wantHeaders))
				return
			}

			for i, h := range table.Headers {
				if h != tt.wantHeaders[i] {
					t.Errorf("NewTable() header[%d] = %q, want %q", i, h, tt.wantHeaders[i])
				}
			}

			if len(table.Widths) != len(tt.wantWidths) {
				t.Errorf("NewTable() widths length = %d, want %d", len(table.Widths), len(tt.wantWidths))
				return
			}

			for i, w := range table.Widths {
				if w != tt.wantWidths[i] {
					t.Errorf("NewTable() width[%d] = %d, want %d", i, w, tt.wantWidths[i])
				}
			}

			if len(table.Rows) != 0 {
				t.Errorf("NewTable() rows length = %d, want 0", len(table.Rows))
			}
		})
	}
}

func TestAddRow(t *testing.T) {
	t.Run("adds row and updates widths", func(t *testing.T) {
		table := NewTable([]string{"ID", "Title"})

		table.AddRow([]string{"1", "Short"})
		if len(table.Rows) != 1 {
			t.Errorf("AddRow() rows length = %d, want 1", len(table.Rows))
		}

		// Width should still be header width since "Short" (5) == "Title" (5)
		if table.Widths[1] != 5 {
			t.Errorf("AddRow() width[1] = %d, want 5", table.Widths[1])
		}

		// Add a longer row
		table.AddRow([]string{"2", "Much longer title here"})
		if len(table.Rows) != 2 {
			t.Errorf("AddRow() rows length = %d, want 2", len(table.Rows))
		}

		// Width should now be updated
		if table.Widths[1] != 22 {
			t.Errorf("AddRow() width[1] = %d, want 22", table.Widths[1])
		}
	})

	t.Run("handles ANSI color codes in width calculation", func(t *testing.T) {
		table := NewTable([]string{"Status"})

		// Add row with ANSI color codes
		coloredText := "\033[32m" + "Done" + "\033[0m"
		table.AddRow([]string{coloredText})

		// Width should be 4 (length of "Done"), not including ANSI codes
		if table.Widths[0] != 6 { // "Status" is 6 chars, "Done" is 4
			t.Errorf("AddRow() width[0] = %d, want 6", table.Widths[0])
		}
	})

	t.Run("handles row shorter than headers", func(t *testing.T) {
		table := NewTable([]string{"A", "B", "C"})
		table.AddRow([]string{"1"}) // Only one cell

		if len(table.Rows) != 1 {
			t.Errorf("AddRow() rows length = %d, want 1", len(table.Rows))
		}
	})
}

func TestStripAnsi(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no ANSI codes",
			input:    "plain text",
			expected: "plain text",
		},
		{
			name:     "single color code",
			input:    "\033[31mred text\033[0m",
			expected: "red text",
		},
		{
			name:     "multiple color codes",
			input:    "\033[1m\033[32mbold green\033[0m",
			expected: "bold green",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only ANSI codes",
			input:    "\033[31m\033[0m",
			expected: "",
		},
		{
			name:     "complex ANSI sequence",
			input:    "\033[38;5;196mextended color\033[0m",
			expected: "extended color",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripAnsi(tt.input)
			if result != tt.expected {
				t.Errorf("stripAnsi(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
