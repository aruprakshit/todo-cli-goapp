package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func cmdAdd(title, priority, category string) error {
	if title == "" {
		return fmt.Errorf("title can not be empty")
	}

	insertSQL := `INSERT INTO todos (title, priority, category) VALUES (?, ?, ?)`

	result, err := db.Exec(insertSQL, title, priority, category)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Printf("%s Added todo #%d: %s\n", colorize(Green, "✓"), id, title)
	return nil
}

func cmdList(showAll, showDone bool, priority, category string) error {
	query := `SELECT id, title, priority, category, done FROM todos`
	conditions := []string{}
	args := []any{}
	if showDone {
		conditions = append(conditions, "done = 1")
	} else if !showAll {
		// default - only pending
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
		return err
	}
	defer rows.Close()

	if showDone {
		fmt.Println("\nCompleted Todos:")
	} else if showAll {
		fmt.Println("\nAll Todos:")
	} else {
		fmt.Println("\nPending Todos:")
	}
	fmt.Println("---------------------------------------")

	table := NewTable([]string{"ID", "✓", "Title", "Priority", "Category"})

	for rows.Next() {
		var id, done int
		var title, priority, category string

		err := rows.Scan(&id, &title, &priority, &category, &done)
		if err != nil {
			return err
		}

		statusDisplay := " "
		if done == 1 {
			statusDisplay = colorize(Green, "✓")
		}
		priorityDisplay := colorize(priorityColor(priority), priority)

		table.AddRow([]string{
			fmt.Sprintf("%d", id),
			statusDisplay,
			title,
			priorityDisplay,
			category,
		})
	}

	if len(table.Rows) == 0 {
		fmt.Println("No todos found")
	} else {
		table.Print()
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func cmdDone(id int) error {
	updateSql := `UPDATE todos SET done = 1 WHERE id = ?`
	result, err := db.Exec(updateSql, id)
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

	fmt.Printf("%s Marked todo #%d as done\n", colorize(Green, "✓"), id)
	return nil
}

func cmdUndone(id int) error {
	updateSql := `UPDATE todos SET done = 0 WHERE id = ?`
	result, err := db.Exec(updateSql, id)
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

	fmt.Printf("%s Marked todo #%d as not done\n", colorize(Blue, "x"), id)
	return nil
}

func cmdDelete(id int, force bool) error {
	var title string
	err := db.QueryRow("SELECT title from todos WHERE id = ?", id).Scan(&title)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("todo #%d not found", id)
		}
		return err
	}

	if !force {
		fmt.Printf("Delete todo #%d: \"%s\"? [y/N] ", id, title)
		var response string
		fmt.Scanln(&response)

		if response != "y" && response != "Y" {
			fmt.Println("Cancelled")
			return nil
		}
	}

	_, err = db.Exec("DELETE from todos WHERE id = ?", id)
	if err != nil {
		return err
	}

	fmt.Printf("%s Deleted todo #%d\n", colorize(Red, "✗"), id)
	return nil
}

func cmdShow(id int) error {
	query := `SELECT id, title, done, priority, category, created_at, due_date FROM todos WHERE id = ?`
	var todoId, done int
	var title, priority, category string
	var createdAt time.Time
	var dueDate sql.NullTime

	err := db.QueryRow(query, id).Scan(&todoId, &title, &done, &priority, &category, &createdAt, &dueDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("todo #%d not found", id)
		}
		return err
	}

	fmt.Println()
	fmt.Println("──────────────────────────────────────")
	fmt.Printf("  ID:        %d\n", todoId)
	fmt.Printf("  Title:     %s\n", title)

	// Show status
	if done == 1 {
		fmt.Printf("  Status:    %s\n", colorize(Green, "Done"))
	} else {
		fmt.Printf("  Status:    %s\n", colorize(Yellow, "Pending"))
	}

	fmt.Printf("  Priority:  %s\n", colorize(priorityColor(priority), priority))

	// Only show category if not empty
	if category != "" {
		fmt.Printf("  Category:  %s\n", category)
	}

	fmt.Printf("  Created:   %s\n", createdAt.Format("2006-01-02 15:04"))

	// Only show due date if set
	if dueDate.Valid {
		fmt.Printf("  Due:       %s\n", dueDate.Time.Format("2006-01-02"))
	}

	fmt.Println("──────────────────────────────────────")
	fmt.Println()
	return nil
}

func cmdEdit(id int, title, priority, category string) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 from todos WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("todo #%d not found", id)
	}

	updates := []string{}
	args := []any{}

	if title != "" {
		updates = append(updates, "title = ?")
		args = append(args, title)
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
		return fmt.Errorf("nothing to update. Use --title, --priority, or --category")
	}

	query := "UPDATE todos SET " + strings.Join(updates, ", ") + " WHERE id = ?"
	args = append(args, id)

	_, err = db.Exec(query, args...)
	if err != nil {
		return err
	}

	fmt.Printf("Updated todo #%d\n", id)
	return nil
}

func cmdClear(clearAll bool) error {
	var count int
	var query string

	if clearAll {
		query = "SELECT COUNT(*) FROM todos"
	} else {
		query = "SELECT COUNT(*) FROM todos WHERE done = 1"
	}

	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		if clearAll {
			fmt.Println("No todos to clear")
		} else {
			fmt.Println("No completed todos to clear")
		}
		return nil
	}

	if clearAll {
		fmt.Printf("Delete ALL %d todos? This cannot be undone. [y/N] ", count)
	} else {
		fmt.Printf("Delete %d completed todos? [y/N] ", count)
	}

	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "Y" {
		fmt.Println("Cancelled")
		return nil
	}

	if clearAll {
		_, err = db.Exec("DELETE FROM todos")
	} else {
		_, err = db.Exec("DELETE FROM todos WHERE done = 1")
	}

	if err != nil {
		return err
	}

	if clearAll {
		fmt.Printf("Cleared all %d todos\n", count)
	} else {
		fmt.Printf("Cleared %d completed todos\n", count)
	}

	return nil

}
