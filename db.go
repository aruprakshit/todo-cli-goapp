import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() error {
	var err error

	// open database
	db, err = sql.Open("sqlite3", "todo.db")
	if err != nil {
		return err
	}

	// Ping to verify
	err = db.Ping()
	if err != nil {
		return err
	}

	// Create table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		done INTEGER DEFAULT 0,
		priority TEXT DEFAULT 'medium',
		category TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		due_date DATETIME
	)`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

}
