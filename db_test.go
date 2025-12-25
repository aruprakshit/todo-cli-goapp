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

func insertTestTodo(t *testing.T, title, priority, category, dueDate string) int64 {
	id, err := insertTodo(title, priority, category, dueDate)
	if err != nil {
		t.Fatalf("failed to insert test todo: %v", err)
	}
	return id
}

func TestInsertTodo(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	tests := []struct {
		name     string
		title    string
		priority string
		category string
		dueDate  string
		wantID   int64
	}{
		{
			name:     "basic todo",
			title:    "Test task",
			priority: "medium",
			category: "",
			dueDate:  "",
			wantID:   1,
		},
		{
			name:     "todo with all fields",
			title:    "Full task",
			priority: "high",
			category: "work",
			dueDate:  "2025-12-31",
			wantID:   2,
		},
		{
			name:     "todo with category only",
			title:    "Category task",
			priority: "low",
			category: "personal",
			dueDate:  "",
			wantID:   3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := insertTodo(tt.title, tt.priority, tt.category, tt.dueDate)
			if err != nil {
				t.Errorf("insertTodo() error = %v", err)
				return
			}
			if id != tt.wantID {
				t.Errorf("insertTodo() id = %v, want %v", id, tt.wantID)
			}
		})
	}
}

func TestGetTodoByID(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	insertTestTodo(t, "Test task", "high", "work", "")

	tests := []struct {
		name      string
		id        int
		wantTitle string
		wantErr   bool
	}{
		{
			name:      "existing todo",
			id:        1,
			wantTitle: "Test task",
			wantErr:   false,
		},
		{
			name:      "non-existing todo",
			id:        999,
			wantTitle: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo, err := getTodoByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTodoByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && todo.Title != tt.wantTitle {
				t.Errorf("getTodoByID() title = %v, want %v", todo.Title, tt.wantTitle)
			}
		})
	}
}

func TestGetTodoByID_FieldValues(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Insert a todo with all fields
	insertTestTodo(t, "Complete task", "high", "work", "2025-12-31")

	todo, err := getTodoByID(1)
	if err != nil {
		t.Fatalf("getTodoByID() error = %v", err)
	}

	if todo.Title != "Complete task" {
		t.Errorf("Title = %v, want %v", todo.Title, "Complete task")
	}
	if todo.Priority != "high" {
		t.Errorf("Priority = %v, want %v", todo.Priority, "high")
	}
	if todo.Category != "work" {
		t.Errorf("Category = %v, want %v", todo.Category, "work")
	}
	if todo.Done != false {
		t.Errorf("Done = %v, want %v", todo.Done, false)
	}
}

func TestGetAllTodos_Empty(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	todos, err := getAllTodos(true, false, "", "")
	if err != nil {
		t.Fatalf("getAllTodos() error = %v", err)
	}
	if len(todos) != 0 {
		t.Errorf("getAllTodos() returned %d todos, want 0", len(todos))
	}
}

func TestGetAllTodos_Filters(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Insert test data
	insertTestTodo(t, "Task 1", "high", "work", "")
	insertTestTodo(t, "Task 2", "low", "personal", "")
	insertTestTodo(t, "Task 3", "high", "work", "")

	// Mark Task 2 as done
	updateTodoStatus(2, true)

	tests := []struct {
		name      string
		showAll   bool
		showDone  bool
		priority  string
		category  string
		wantCount int
	}{
		{
			name:      "all todos",
			showAll:   true,
			showDone:  false,
			priority:  "",
			category:  "",
			wantCount: 3,
		},
		{
			name:      "pending only (default)",
			showAll:   false,
			showDone:  false,
			priority:  "",
			category:  "",
			wantCount: 2,
		},
		{
			name:      "done only",
			showAll:   false,
			showDone:  true,
			priority:  "",
			category:  "",
			wantCount: 1,
		},
		{
			name:      "filter by priority",
			showAll:   true,
			showDone:  false,
			priority:  "high",
			category:  "",
			wantCount: 2,
		},
		{
			name:      "filter by category",
			showAll:   true,
			showDone:  false,
			priority:  "",
			category:  "work",
			wantCount: 2,
		},
		{
			name:      "filter by priority and category",
			showAll:   true,
			showDone:  false,
			priority:  "low",
			category:  "personal",
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todos, err := getAllTodos(tt.showAll, tt.showDone, tt.priority, tt.category)
			if err != nil {
				t.Errorf("getAllTodos() error = %v", err)
				return
			}
			if len(todos) != tt.wantCount {
				t.Errorf("getAllTodos() returned %d todos, want %d", len(todos), tt.wantCount)
			}
		})
	}
}

func TestUpdateTodoStatus(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	insertTestTodo(t, "Test task", "medium", "", "")

	// Mark as done
	err := updateTodoStatus(1, true)
	if err != nil {
		t.Fatalf("updateTodoStatus(done=true) error = %v", err)
	}

	// Verify it's done
	todo, _ := getTodoByID(1)
	if !todo.Done {
		t.Errorf("Todo should be done, but Done = %v", todo.Done)
	}

	// Mark as undone
	err = updateTodoStatus(1, false)
	if err != nil {
		t.Fatalf("updateTodoStatus(done=false) error = %v", err)
	}

	// Verify it's not done
	todo, _ = getTodoByID(1)
	if todo.Done {
		t.Errorf("Todo should be undone, but Done = %v", todo.Done)
	}
}

func TestUpdateTodoStatus_NotFound(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	err := updateTodoStatus(999, true)
	if err == nil {
		t.Error("updateTodoStatus() should return error for non-existing todo")
	}
}

func TestDeleteTodo(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	insertTestTodo(t, "Test task", "medium", "", "")

	err := deleteTodo(1)
	if err != nil {
		t.Fatalf("deleteTodo() error = %v", err)
	}

	// Verify it's gone
	_, err = getTodoByID(1)
	if err == nil {
		t.Error("getTodoByID() should return error for deleted todo")
	}
}

func TestDeleteTodo_NonExisting(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Deleting non-existing todo should not return error (no rows affected)
	err := deleteTodo(999)
	if err != nil {
		t.Errorf("deleteTodo() error = %v, want nil", err)
	}
}

func TestUpdateTodo(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Insert a test todo
	insertTestTodo(t, "Original title", "low", "personal", "")

	tests := []struct {
		name         string
		title        string
		priority     string
		category     string
		dueDate      string
		wantTitle    string
		wantPriority string
		wantCategory string
	}{
		{
			name:         "update title only",
			title:        "New title",
			priority:     "",
			category:     "",
			dueDate:      "",
			wantTitle:    "New title",
			wantPriority: "low",
			wantCategory: "personal",
		},
		{
			name:         "update priority only",
			title:        "",
			priority:     "high",
			category:     "",
			dueDate:      "",
			wantTitle:    "New title",
			wantPriority: "high",
			wantCategory: "personal",
		},
		{
			name:         "update category only",
			title:        "",
			priority:     "",
			category:     "work",
			dueDate:      "",
			wantTitle:    "New title",
			wantPriority: "high",
			wantCategory: "work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := updateTodo(1, tt.title, tt.priority, tt.category, tt.dueDate)
			if err != nil {
				t.Errorf("updateTodo() error = %v", err)
				return
			}

			todo, _ := getTodoByID(1)
			if todo.Title != tt.wantTitle {
				t.Errorf("Title = %v, want %v", todo.Title, tt.wantTitle)
			}
			if todo.Priority != tt.wantPriority {
				t.Errorf("Priority = %v, want %v", todo.Priority, tt.wantPriority)
			}
			if todo.Category != tt.wantCategory {
				t.Errorf("Category = %v, want %v", todo.Category, tt.wantCategory)
			}
		})
	}
}

func TestUpdateTodo_NoChanges(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	insertTestTodo(t, "Test task", "medium", "", "")

	// Update with empty values should not error
	err := updateTodo(1, "", "", "", "")
	if err != nil {
		t.Errorf("updateTodo() with no changes error = %v", err)
	}
}

func TestCountTodos(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Insert test data
	insertTestTodo(t, "Task 1", "medium", "", "")
	insertTestTodo(t, "Task 2", "medium", "", "")
	insertTestTodo(t, "Task 3", "medium", "", "")

	// Mark one as done
	updateTodoStatus(2, true)

	tests := []struct {
		name          string
		completedOnly bool
		wantCount     int
	}{
		{
			name:          "all todos",
			completedOnly: false,
			wantCount:     3,
		},
		{
			name:          "completed only",
			completedOnly: true,
			wantCount:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := countTodos(tt.completedOnly)
			if err != nil {
				t.Errorf("countTodos() error = %v", err)
				return
			}
			if count != tt.wantCount {
				t.Errorf("countTodos() = %v, want %v", count, tt.wantCount)
			}
		})
	}
}

func TestClearTodos_CompletedOnly(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	insertTestTodo(t, "Task 1", "medium", "", "")
	insertTestTodo(t, "Task 2", "medium", "", "")
	insertTestTodo(t, "Task 3", "medium", "", "")

	// Mark two as done
	updateTodoStatus(1, true)
	updateTodoStatus(2, true)

	// Clear completed only
	err := clearTodos(false)
	if err != nil {
		t.Fatalf("clearTodos(false) error = %v", err)
	}

	// Should have 1 remaining
	count, _ := countTodos(false)
	if count != 1 {
		t.Errorf("After clearTodos(false), count = %v, want 1", count)
	}
}

func TestClearTodos_All(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	insertTestTodo(t, "Task 1", "medium", "", "")
	insertTestTodo(t, "Task 2", "medium", "", "")
	updateTodoStatus(1, true)

	// Clear all
	err := clearTodos(true)
	if err != nil {
		t.Fatalf("clearTodos(true) error = %v", err)
	}

	// Should have 0 remaining
	count, _ := countTodos(false)
	if count != 0 {
		t.Errorf("After clearTodos(true), count = %v, want 0", count)
	}
}
