package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// initialiize the database
	err := initDB()

	if err != nil {
		fmt.Println("Error initializing database: ", err)
		os.Exit(1)
	}
	defer db.Close()

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)

		// Define flags
		priority := addCmd.String("priority", "medium", "Prioriy: low, medium, high")
		category := addCmd.String("caegory", "", "Category for the todo")

		addCmd.Parse(os.Args[2:])
		args := addCmd.Args()
		if len(args) < 1 {
			fmt.Println("Error: title is required")
			fmt.Println("Usage: todo add <title> [--priority low|medium|high] [--category name]")
			os.Exit(1)
		}
		title := args[0]

		err := cmdAdd(title, *priority, *category)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "list":
		err := cmdList()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknownn command: %s\n", command)
		printUsage()
		os.Exit(1)
	}

}

func printUsage() {
	fmt.Println("Usage: todo <command> [options]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  add     Add a new todo")
	fmt.Println("  list    List pending todos")
}
