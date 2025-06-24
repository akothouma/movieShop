package database

import (
	"database/sql"
	"fmt"
	"log"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL(Write After Logging) mode for better performance
	if _, err := db.Exec(`
		PRAGMA journal_mode = WAL;
		PRAGMA synchronous = NORMAL;
		PRAGMA foreign_keys = ON;
	`); err != nil {
		return nil, fmt.Errorf("failed to set SQLite pragmas: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(1) // SQLite only allows one writer at a time
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0) // Connections don't timeout

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to SQLite database at", dbPath)
	return db, nil
}

func CloseDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("Error closing SQLite database: %v", err)
		return
	}
	log.Println("Closed SQLite connection")
}