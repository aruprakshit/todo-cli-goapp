package main

import (
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
