package main

import (
	"testing"
)

func TestInsertTodo(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	tests := []struct {
		name     string
		title    string
		priority Priority
		category string
		dueDate  string
		wantID   int64
	}{
		{
			name:     "basic todo",
			title:    "Test task",
			priority: PriorityMedium,
			category: "",
			dueDate:  "",
			wantID:   1,
		},
		{
			name:     "todo with all fields",
			title:    "Full task",
			priority: PriorityHigh,
			category: "work",
			dueDate:  "2025-12-31",
			wantID:   2,
		},
		{
			name:     "todo with category only",
			title:    "Category task",
			priority: PriorityLow,
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

	insertTestTodo(t, "Test task", PriorityHigh, "work", "")

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
	insertTestTodo(t, "Complete task", PriorityHigh, "work", "2025-12-31")

	todo, err := getTodoByID(1)
	if err != nil {
		t.Fatalf("getTodoByID() error = %v", err)
	}

	if todo.Title != "Complete task" {
		t.Errorf("Title = %v, want %v", todo.Title, "Complete task")
	}
	if todo.Priority != PriorityHigh {
		t.Errorf("Priority = %v, want %v", todo.Priority, PriorityHigh)
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
	insertTestTodo(t, "Task 1", PriorityHigh, "work", "")
	insertTestTodo(t, "Task 2", PriorityLow, "personal", "")
	insertTestTodo(t, "Task 3", PriorityHigh, "work", "")

	// Mark Task 2 as done
	markTodoAsDone(2)

	tests := []struct {
		name      string
		showAll   bool
		showDone  bool
		priority  Priority
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
			priority:  PriorityHigh,
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
			priority:  PriorityLow,
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

func TestMarkTodoAsDone(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	insertTestTodo(t, "Test task", PriorityMedium, "", "")

	err := markTodoAsDone(1)
	if err != nil {
		t.Fatalf("markTodoAsDone() error = %v", err)
	}

	todo, _ := getTodoByID(1)
	if !todo.Done {
		t.Errorf("Todo should be done, but Done = %v", todo.Done)
	}
}

func TestMarkTodoAsUndone(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	insertTestTodo(t, "Test task", PriorityMedium, "", "")
	markTodoAsDone(1)

	err := markTodoAsUndone(1)
	if err != nil {
		t.Fatalf("markTodoAsUndone() error = %v", err)
	}

	todo, _ := getTodoByID(1)
	if todo.Done {
		t.Errorf("Todo should be undone, but Done = %v", todo.Done)
	}
}

func TestMarkTodoAsDone_NotFound(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	err := markTodoAsDone(999)
	if err == nil {
		t.Error("markTodoAsDone() should return error for non-existing todo")
	}
}

func TestDeleteTodo(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	insertTestTodo(t, "Test task", PriorityMedium, "", "")

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
	insertTestTodo(t, "Original title", PriorityLow, "personal", "")

	tests := []struct {
		name         string
		title        string
		priority     Priority
		category     string
		dueDate      string
		wantTitle    string
		wantPriority Priority
		wantCategory string
	}{
		{
			name:         "update title only",
			title:        "New title",
			priority:     "",
			category:     "",
			dueDate:      "",
			wantTitle:    "New title",
			wantPriority: PriorityLow,
			wantCategory: "personal",
		},
		{
			name:         "update priority only",
			title:        "",
			priority:     PriorityHigh,
			category:     "",
			dueDate:      "",
			wantTitle:    "New title",
			wantPriority: PriorityHigh,
			wantCategory: "personal",
		},
		{
			name:         "update category only",
			title:        "",
			priority:     "",
			category:     "work",
			dueDate:      "",
			wantTitle:    "New title",
			wantPriority: PriorityHigh,
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

	insertTestTodo(t, "Test task", PriorityMedium, "", "")

	// Update with empty values should not error
	err := updateTodo(1, "", "", "", "")
	if err != nil {
		t.Errorf("updateTodo() with no changes error = %v", err)
	}
}

func TestCountAllTodos(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	insertTestTodo(t, "Task 1", PriorityMedium, "", "")
	insertTestTodo(t, "Task 2", PriorityMedium, "", "")
	insertTestTodo(t, "Task 3", PriorityMedium, "", "")

	count, err := countAllTodos()
	if err != nil {
		t.Fatalf("countAllTodos() error = %v", err)
	}
	if count != 3 {
		t.Errorf("countAllTodos() = %v, want 3", count)
	}
}

func TestCountCompletedTodos(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	insertTestTodo(t, "Task 1", PriorityMedium, "", "")
	insertTestTodo(t, "Task 2", PriorityMedium, "", "")
	insertTestTodo(t, "Task 3", PriorityMedium, "", "")
	markTodoAsDone(2)

	count, err := countCompletedTodos()
	if err != nil {
		t.Fatalf("countCompletedTodos() error = %v", err)
	}
	if count != 1 {
		t.Errorf("countCompletedTodos() = %v, want 1", count)
	}
}

func TestClearCompletedTodos(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	insertTestTodo(t, "Task 1", PriorityMedium, "", "")
	insertTestTodo(t, "Task 2", PriorityMedium, "", "")
	insertTestTodo(t, "Task 3", PriorityMedium, "", "")

	markTodoAsDone(1)
	markTodoAsDone(2)

	err := clearCompletedTodos()
	if err != nil {
		t.Fatalf("clearCompletedTodos() error = %v", err)
	}

	count, _ := countAllTodos()
	if count != 1 {
		t.Errorf("After clearCompletedTodos(), count = %v, want 1", count)
	}
}

func TestClearAllTodos(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	insertTestTodo(t, "Task 1", PriorityMedium, "", "")
	insertTestTodo(t, "Task 2", PriorityMedium, "", "")
	markTodoAsDone(1)

	err := clearAllTodos()
	if err != nil {
		t.Fatalf("clearAllTodos() error = %v", err)
	}

	count, _ := countAllTodos()
	if count != 0 {
		t.Errorf("After clearAllTodos(), count = %v, want 0", count)
	}
}
