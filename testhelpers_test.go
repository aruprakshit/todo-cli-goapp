package main

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB creates an in-memory database for testing
func setupTestDB(t *testing.T) {
	var err error
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	err = createTables()
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}
}

// teardownTestDB closes the test database
func teardownTestDB() {
	if db != nil {
		db.Close()
	}
}

func insertTestTodo(t *testing.T, title string, priority Priority, category, dueDate string) int64 {
	id, err := insertTodo(title, priority, category, dueDate)
	if err != nil {
		t.Fatalf("failed to insert test todo: %v", err)
	}
	return id
}
