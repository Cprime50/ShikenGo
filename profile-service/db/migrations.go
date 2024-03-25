package db

import (
	"database/sql"
	"fmt"
)

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS profiles (
            id TEXT PRIMARY KEY,
            user_id TEXT UNIQUE NOT NULL,
            email TEXT UNIQUE NOT NULL,
            username TEXT UNIQUE NOT NULL,
            bio TEXT,
            avatar TEXT,
            score INTEGER DEFAULT 0,
            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        );
    `)
	if err != nil {
		return fmt.Errorf("Error creating table profiles: %w", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS user_id ON profiles (user_id)`)
	if err != nil {
		return fmt.Errorf("Error creating index: %w", err)
	}

	fmt.Println("Migration successful.")
	return nil
}
