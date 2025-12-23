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
