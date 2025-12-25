package main

import (
	"database/sql"
	"time"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

func (p Priority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	}
	return false
}

type Todo struct {
	ID        int
	Title     string
	Done      bool
	Priority  Priority
	Category  string
	CreatedAt time.Time
	DueDate   sql.NullTime
}
