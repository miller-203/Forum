package forum

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	}

	
	if err := db.Ping(); err != nil {
		return fmt.Errorf("unable to reach database: %v", err)
	}

	return createTables()
}

func createTables() error {
	// Table dyal users
	usersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,    -- smiya dyal user
        password TEXT NOT NULL,           -- mot de passe
        email TEXT UNIQUE NOT NULL,       -- email
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP  -- date creation
    );`

	_, err := db.Exec(usersTable)
	if err != nil {
		return err
	}
	// Table for posts
	postsTable := `
	  CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,                                  -- Username of the user who created the post
    title TEXT NOT NULL,                                     -- Title of the post
    content TEXT NOT NULL,                                   -- Content of the post
	Categories VARCHAR(50) NOT NULL,                                -- Categories
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,           -- Creation timestamp
    FOREIGN KEY (username) REFERENCES users(username)       -- Foreign key to the users table (assuming users.username is TEXT)
);
	  `
	_, err = db.Exec(postsTable)
	if err != nil {
		return err
	}
	return nil
}
