package main

import (
	"testing"
)

func TestPriorityIsValid(t *testing.T) {
	tests := []struct {
		name     string
		priority Priority
		want     bool
	}{
		{
			name:     "low priority is valid",
			priority: PriorityLow,
			want:     true,
		},
		{
			name:     "medium priority is valid",
			priority: PriorityMedium,
			want:     true,
		},
		{
			name:     "high priority is valid",
			priority: PriorityHigh,
			want:     true,
		},
		{
			name:     "empty string is invalid",
			priority: Priority(""),
			want:     false,
		},
		{
			name:     "random string is invalid",
			priority: Priority("urgent"),
			want:     false,
		},
		{
			name:     "uppercase LOW is invalid",
			priority: Priority("LOW"),
			want:     false,
		},
		{
			name:     "mixed case Medium is invalid",
			priority: Priority("Medium"),
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.priority.IsValid()
			if got != tt.want {
				t.Errorf("Priority(%q).IsValid() = %v, want %v", tt.priority, got, tt.want)
			}
		})
	}
}
