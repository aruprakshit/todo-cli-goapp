package main

import (
	"database/sql"
	"time"
)

type Todo struct {
	ID        int
	Title     string
	Done      bool
	Priority  string
	Category  string
	CreatedAt time.Time
	DueDate   sql.NullTime
}
