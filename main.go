package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
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
		category := addCmd.String("category", "", "Category for the todo")

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
		listCmd := flag.NewFlagSet("list", flag.ExitOnError)
		showAll := listCmd.Bool("all", false, "Show all todos")
		showDone := listCmd.Bool("done", false, "Show only completed")
		priority := listCmd.String("priority", "", "Filter by priority")
		category := listCmd.String("category", "", "Filter by category")
		listCmd.Parse(os.Args[2:])

		err := cmdList(*showAll, *showDone, *priority, *category)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo done <id>")
			os.Exit(1)
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: invalid ID")
			os.Exit(1)
		}
		err = cmdDone(id)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "undone":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo undone <id>")
			os.Exit(1)
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: invalid ID")
			os.Exit(1)
		}
		err = cmdUndone(id)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "delete":
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		force := deleteCmd.Bool("force", false, "Skip confirmation")
		deleteCmd.Parse(os.Args[2:])

		args := deleteCmd.Args()
		if len(args) < 1 {
			fmt.Println("Usage: todo delete <id> [--force]")
			os.Exit(1)
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: invalid ID")
			os.Exit(1)
		}
		err = cmdDelete(id, *force)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "show":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo show <id>")
			os.Exit(1)
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: invalid ID")
			os.Exit(1)
		}
		err = cmdShow(id)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "edit":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo edit <id> [--title text] [--priority low|medium|high] [--category name]")
			os.Exit(1)
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: invalid ID")
			os.Exit(1)
		}

		editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
		title := editCmd.String("title", "", "New title")
		priority := editCmd.String("priority", "", "New priority")
		category := editCmd.String("category", "", "New category")
		editCmd.Parse(os.Args[3:])

		err = cmdEdit(id, *title, *priority, *category)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "clear":
		clearCmd := flag.NewFlagSet("clear", flag.ExitOnError)
		clearAll := clearCmd.Bool("all", false, "Clear ALL todos")
		clearCmd.Parse(os.Args[2:])

		err := cmdClear(*clearAll)
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
	fmt.Println("  done    Mark a todo as complete")
	fmt.Println("  undone  Mark a todo as incomplete")
	fmt.Println("  delete  Delete a todo")
	fmt.Println("  show    Show todo details")
	fmt.Println("  edit    Edit a todo")
	fmt.Println("  clear   Remove completed todos")
}
