package main

import (
	"fmt"
	"time"
)

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func cmdAdd(title string, priority Priority, category, dueDate string) error {
	if title == "" {
		return fmt.Errorf("title can not be empty")
	}

	if !priority.IsValid() {
		return fmt.Errorf("invalid priority: %s. Use low, medium, or high", priority)
	}

	if dueDate != "" {
		_, err := parseDate(dueDate)
		if err != nil {
			return fmt.Errorf("invalid date format. Use YYYY-MM-DD")
		}
	}

	id, err := insertTodo(title, priority, category, dueDate)
	if err != nil {
		return err
	}

	fmt.Printf("%s Added todo #%d: %s\n", colorize(Green, "✓"), id, title)
	return nil
}

func cmdList(showAll, showDone bool, priority Priority, category string) error {
	todos, err := getAllTodos(showAll, showDone, priority, category)
	if err != nil {
		return err
	}

	if showDone {
		fmt.Println("\nCompleted Todos:")
	} else if showAll {
		fmt.Println("\nAll Todos:")
	} else {
		fmt.Println("\nPending Todos:")
	}
	fmt.Println("---------------------------------------")

	table := NewTable([]string{"ID", "✓", "Title", "Priority", "Category", "Due"})

	for _, todo := range todos {
		statusDisplay := " "
		if todo.Done {
			statusDisplay = colorize(Green, "✓")
		}
		priorityDisplay := colorize(priorityColor(todo.Priority), string(todo.Priority))
		dueDateDisplay := formatDueDate(todo.DueDate)

		table.AddRow([]string{
			fmt.Sprintf("%d", todo.ID),
			statusDisplay,
			todo.Title,
			priorityDisplay,
			todo.Category,
			dueDateDisplay,
		})
	}

	if len(table.Rows) == 0 {
		fmt.Println("No todos found")
	} else {
		table.Print()
	}

	return nil
}

func cmdDone(id int) error {
	err := markTodoAsDone(id)
	if err != nil {
		return err
	}

	fmt.Printf("%s Marked todo #%d as done\n", colorize(Green, "✓"), id)
	return nil
}

func cmdUndone(id int) error {
	err := markTodoAsUndone(id)
	if err != nil {
		return err
	}

	fmt.Printf("%s Marked todo #%d as not done\n", colorize(Blue, "x"), id)
	return nil
}

func cmdDelete(id int, force bool) error {
	todo, err := getTodoByID(id)
	if err != nil {
		return err
	}

	if !force {
		fmt.Printf("Delete todo #%d: \"%s\"? [y/N] ", todo.ID, todo.Title)
		var response string
		fmt.Scanln(&response)

		if response != "y" && response != "Y" {
			fmt.Println("Cancelled")
			return nil
		}
	}

	err = deleteTodo(id)
	if err != nil {
		return err
	}

	fmt.Printf("%s Deleted todo #%d\n", colorize(Red, "✗"), id)
	return nil
}

func cmdShow(id int) error {
	todo, err := getTodoByID(id)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("──────────────────────────────────────")
	fmt.Printf("  ID:        %d\n", todo.ID)
	fmt.Printf("  Title:     %s\n", todo.Title)

	// Show status
	if todo.Done {
		fmt.Printf("  Status:    %s\n", colorize(Green, "Done"))
	} else {
		fmt.Printf("  Status:    %s\n", colorize(Yellow, "Pending"))
	}

	fmt.Printf("  Priority:  %s\n", colorize(priorityColor(todo.Priority), string(todo.Priority)))

	// Only show category if not empty
	if todo.Category != "" {
		fmt.Printf("  Category:  %s\n", todo.Category)
	}

	fmt.Printf("  Created:   %s\n", todo.CreatedAt.Format("2006-01-02 15:04"))

	// Only show due date if set
	if todo.DueDate.Valid {
		fmt.Printf("  Due:       %s\n", formatDueDate(todo.DueDate))
	}

	fmt.Println("──────────────────────────────────────")
	fmt.Println()
	return nil
}

func cmdEdit(id int, title string, priority Priority, category, dueDate string) error {
	_, err := getTodoByID(id)
	if err != nil {
		return err
	}

	if dueDate != "" {
		_, err := parseDate(dueDate)
		if err != nil {
			return fmt.Errorf("invalid date format. Use YYYY-MM-DD")
		}
	}

	if priority != "" && !priority.IsValid() {
		return fmt.Errorf("invalid priority: %s. Use low, medium, or high", priority)
	}

	if title == "" && priority == "" && category == "" && dueDate == "" {
		return fmt.Errorf("nothing to update. Use --title, --priority, --category, or --due")
	}

	err = updateTodo(id, title, priority, category, dueDate)
	if err != nil {
		return err
	}

	fmt.Printf("Updated todo #%d\n", id)
	return nil
}

func cmdClear(clearAll bool) error {
	var count int
	var err error

	if clearAll {
		count, err = countAllTodos()
	} else {
		count, err = countCompletedTodos()
	}
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
		err = clearAllTodos()
	} else {
		err = clearCompletedTodos()
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
