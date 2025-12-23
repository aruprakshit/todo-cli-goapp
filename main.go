package main

import "fmt"

func main() {
	err := initDB()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Database initialized successfully!")

	defer db.Close()

	// Add a test todo
	err = cmdAdd("Buy groceries", "low", "shopping")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// List all todos
	err = cmdList()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
