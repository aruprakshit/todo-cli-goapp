# Todo CLI

A command-line todo application built with Go and SQLite.

## Features

- Create, read, update, and delete todos
- Mark todos as done/undone
- Prioritize tasks (low, medium, high)
- Categorize tasks
- Filter by status, priority, or category
- Bulk clear completed todos
- Persistent storage with SQLite

## Prerequisites

- Go 1.21 or higher
- GCC (C compiler) - required for SQLite driver

### Install GCC

```bash
# Ubuntu/Debian
sudo apt install build-essential

# Fedora
sudo dnf install gcc
```

## Getting Started

### 1. Clone the repository

```bash
git clone <repository-url>
cd todo-cli
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Build the application

```bash
go build -o todo
```

### 4. Run the application

```bash
./todo
```

## Usage

### Add a todo

```bash
./todo add "Buy groceries"
./todo add --priority high "Urgent task"
./todo add --priority low --category personal "Read a book"
./todo add --priority high --category work "Finish report"
```

**Flags:**
- `--priority` - Set priority: low, medium (default), high
- `--category` - Set category name

### List todos

```bash
./todo list                        # Show pending todos only
./todo list --all                  # Show all todos
./todo list --done                 # Show completed todos only
./todo list --priority high        # Filter by priority
./todo list --category work        # Filter by category
./todo list --all --category work  # Combine filters
```

**Flags:**
- `--all` - Show all todos (pending and completed)
- `--done` - Show only completed todos
- `--priority` - Filter by priority level
- `--category` - Filter by category name

### Show todo details

```bash
./todo show 1
```

### Mark as done/undone

```bash
./todo done 1      # Mark todo #1 as complete
./todo undone 1    # Mark todo #1 as incomplete
```

### Edit a todo

```bash
./todo edit 1 --title "Updated title"
./todo edit 1 --priority high
./todo edit 1 --category work
./todo edit 1 --title "New title" --priority low --category personal
```

**Flags:**
- `--title` - New title
- `--priority` - New priority
- `--category` - New category

### Delete a todo

```bash
./todo delete 1           # Asks for confirmation
./todo delete --force 1   # Skip confirmation
```

**Flags:**
- `--force` - Skip confirmation prompt

### Clear todos

```bash
./todo clear        # Remove all completed todos
./todo clear --all  # Remove ALL todos (with confirmation)
```

**Flags:**
- `--all` - Clear all todos, not just completed ones

## Command Reference

| Command | Description |
|---------|-------------|
| `add <title>` | Add a new todo |
| `list` | List todos |
| `show <id>` | Show todo details |
| `done <id>` | Mark todo as complete |
| `undone <id>` | Mark todo as incomplete |
| `edit <id>` | Edit a todo |
| `delete <id>` | Delete a todo |
| `clear` | Remove completed todos |

## Project Structure

```
todo-cli/
├── main.go       # Entry point, command routing
├── db.go         # Database operations
├── models.go     # Data structures
├── commands.go   # Command handlers
├── go.mod        # Go module file
├── go.sum        # Dependency checksums
└── todo.db       # SQLite database (created on first run)
```

## Database

The application uses SQLite for data persistence. The database file `todo.db` is created automatically in the current directory on first run.

### Schema

```sql
CREATE TABLE todos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    done INTEGER DEFAULT 0,
    priority TEXT DEFAULT 'medium',
    category TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    due_date DATETIME
)
```

## Testing

### Run all tests

```bash
go test ./...           # Run all tests
go test -v ./...        # Run with verbose output
```

### Run specific tests

```bash
go test -v -run TestFunctionName ./...     # Run tests matching a pattern
go test -v -run TestInsertTodo ./...       # Example: run TestInsertTodo
go test -v -run TestInsertTodo/basic ./... # Run a specific subtest
```

### Test coverage

```bash
# Basic coverage percentage
go test -cover ./...

# Detailed per-function coverage
mkdir -p coverage && go test -coverprofile=coverage/coverage.out ./... && go tool cover -func=coverage/coverage.out

# HTML report (opens in browser)
mkdir -p coverage && go test -coverprofile=coverage/coverage.out ./... && go tool cover -html=coverage/coverage.out -o coverage/coverage.html

# Coverage with race detection
go test -race -cover ./...
```

> **Note:** Add `coverage/` to your `.gitignore` file.

## Dependencies

- [go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite driver for Go

## License

MIT
