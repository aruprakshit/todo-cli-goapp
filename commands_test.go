package main

import (
	"strings"
	"testing"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid date",
			input:   "2024-12-25",
			wantErr: false,
		},
		{
			name:    "valid date with different values",
			input:   "2023-01-15",
			wantErr: false,
		},
		{
			name:    "invalid format MM-DD-YYYY",
			input:   "12-25-2024",
			wantErr: true,
		},
		{
			name:    "invalid format with slashes",
			input:   "2024/12/25",
			wantErr: true,
		},
		{
			name:    "invalid month",
			input:   "2024-13-25",
			wantErr: true,
		},
		{
			name:    "invalid day",
			input:   "2024-12-32",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "random text",
			input:   "not-a-date",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDate(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestCmdAdd_Validation(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		priority    Priority
		category    string
		dueDate     string
		wantErr     bool
		errContains string
	}{
		{
			name:        "empty title",
			title:       "",
			priority:    PriorityMedium,
			category:    "",
			dueDate:     "",
			wantErr:     true,
			errContains: "title can not be empty",
		},
		{
			name:        "invalid priority",
			title:       "Test task",
			priority:    Priority("urgent"),
			category:    "",
			dueDate:     "",
			wantErr:     true,
			errContains: "invalid priority",
		},
		{
			name:        "invalid date format",
			title:       "Test task",
			priority:    PriorityMedium,
			category:    "",
			dueDate:     "25-12-2024",
			wantErr:     true,
			errContains: "invalid date format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cmdAdd(tt.title, tt.priority, tt.category, tt.dueDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("cmdAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("cmdAdd() error = %q, want it to contain %q", err.Error(), tt.errContains)
			}
		})
	}
}

func TestCmdEdit_Validation(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Insert a test todo to edit
	testID := int(insertTestTodo(t, "Test task", PriorityMedium, "work", ""))

	tests := []struct {
		name        string
		id          int
		title       string
		priority    Priority
		category    string
		dueDate     string
		wantErr     bool
		errContains string
	}{
		{
			name:        "non-existent todo",
			id:          999,
			title:       "New title",
			priority:    "",
			category:    "",
			dueDate:     "",
			wantErr:     true,
			errContains: "not found",
		},
		{
			name:        "invalid date format",
			id:          testID,
			title:       "",
			priority:    "",
			category:    "",
			dueDate:     "25-12-2024",
			wantErr:     true,
			errContains: "invalid date format",
		},
		{
			name:        "invalid priority",
			id:          testID,
			title:       "",
			priority:    Priority("urgent"),
			category:    "",
			dueDate:     "",
			wantErr:     true,
			errContains: "invalid priority",
		},
		{
			name:        "nothing to update",
			id:          testID,
			title:       "",
			priority:    "",
			category:    "",
			dueDate:     "",
			wantErr:     true,
			errContains: "nothing to update",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cmdEdit(tt.id, tt.title, tt.priority, tt.category, tt.dueDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("cmdEdit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("cmdEdit() error = %q, want it to contain %q", err.Error(), tt.errContains)
			}
		})
	}
}
