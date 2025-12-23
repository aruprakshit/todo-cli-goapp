package main

import (
	"database/sql"
	"fmt"
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

	fmt.Printf("Added todo #%d: %s\n", id, title)
	return nil
}

func cmdList() error {
	query := "SELECT id, title, priority, category, done FROM todos WHERE done = 0"
	rows, err := db.Query(query)

	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Println("\n Pending Todos:")
	fmt.Println("---------------------------------------")

	for rows.Next() {
		var id, done int
		var title, priority, category string

		err := rows.Scan(&id, &title, &priority, &category, &done)
		if err != nil {
			return err
		}

		fmt.Printf("[%d] %s (priority: %s", id, title, priority)
		if category != "" {
			fmt.Printf(", Category: %s", category)
		}
		fmt.Println(")")
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

	fmt.Printf("Marked todo #%d as done\n", id)
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

	fmt.Printf("Marked todo #%d as not done\n", id)
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

	fmt.Printf("Deleted todo #%d\n", id)
	return nil
}
