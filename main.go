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
}
