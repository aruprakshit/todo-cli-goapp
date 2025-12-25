package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func getTodoByID(id int) (*Todo, error) {
	query := `SELECT id, title, done, priority, category, created_at, due_date FROM todos WHERE id = ?`
	row := db.QueryRow(query, id)

	todo, err := scanTodo(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("todo #%d not found", id)
		}
		return nil, err
	}

	return todo, nil
}

func getAllTodos(showAll, showDone bool, priority, category string) ([]Todo, error) {
	query := `SELECT id, title, done, priority, category, created_at, due_date FROM todos`
	conditions := []string{}
	args := []any{}

	if showDone {
		conditions = append(conditions, "done = 1")
	} else if !showAll {
		conditions = append(conditions, "done = 0")
	}

	if category != "" {
		conditions = append(conditions, "category = ?")
		args = append(args, category)
	}

	if priority != "" {
		conditions = append(conditions, "priority = ?")
		args = append(args, priority)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		var done int

		err := rows.Scan(&todo.ID, &todo.Title, &done, &todo.Priority, &todo.Category, &todo.CreatedAt, &todo.DueDate)
		if err != nil {
			return nil, err
		}

		todo.Done = done == 1
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func scanTodo(row *sql.Row) (*Todo, error) {
	var todo Todo
	var done int

	err := row.Scan(&todo.ID, &todo.Title, &done, &todo.Priority, &todo.Category, &todo.CreatedAt, &todo.DueDate)
	if err != nil {
		return nil, err
	}

	todo.Done = done == 1
	return &todo, nil
}

func insertTodo(title, priority, category, dueDate string) (int64, error) {
	var result sql.Result
	var err error

	if dueDate != "" {
		query := `INSERT INTO todos (title, priority, category, due_date) VALUES (?, ?, ?, ?)`
		result, err = db.Exec(query, title, priority, category, dueDate)
	} else {
		query := `INSERT INTO todos (title, priority, category) VALUES (?, ?, ?)`
		result, err = db.Exec(query, title, priority, category)
	}

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func markTodoAsDone(id int) error {
	return setTodoStatus(id, 1)
}

func markTodoAsUndone(id int) error {
	return setTodoStatus(id, 0)
}

func setTodoStatus(id int, done int) error {
	result, err := db.Exec(`UPDATE todos SET done = ? WHERE id = ?`, done, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo #%d not found", id)
	}

	return nil
}

func deleteTodo(id int) error {
	_, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}

func updateTodo(id int, title, priority, category, dueDate string) error {
	updates := []string{}
	args := []any{}

	if title != "" {
		updates = append(updates, "title = ?")
		args = append(args, title)
	}

	if dueDate != "" {
		updates = append(updates, "due_date = ?")
		args = append(args, dueDate)
	}

	if priority != "" {
		updates = append(updates, "priority = ?")
		args = append(args, priority)
	}

	if category != "" {
		updates = append(updates, "category = ?")
		args = append(args, category)
	}

	if len(updates) == 0 {
		return nil
	}

	query := "UPDATE todos SET " + strings.Join(updates, ", ") + " WHERE id = ?"
	args = append(args, id)

	_, err := db.Exec(query, args...)
	return err
}

func countAllTodos() (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM todos").Scan(&count)
	return count, err
}

func countCompletedTodos() (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM todos WHERE done = 1").Scan(&count)
	return count, err
}

func clearAllTodos() error {
	_, err := db.Exec("DELETE FROM todos")
	return err
}

func clearCompletedTodos() error {
	_, err := db.Exec("DELETE FROM todos WHERE done = 1")
	return err
}

func createTables() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		done INTEGER DEFAULT 0,
		priority TEXT DEFAULT 'medium',
		category TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		due_date DATETIME
	)`

	_, err := db.Exec(createTableSQL)
	return err
}

func initDB() error {
	var err error

	// open database
	db, err = sql.Open("sqlite3", "todo.db")
	if err != nil {
		return err
	}

	// Ping to verify
	err = db.Ping()
	if err != nil {
		return err
	}

	return createTables()
}
