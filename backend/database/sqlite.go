package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the SQLite database and creates all necessary tables
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create tables
	if err := createTables(db); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}

func createTables(db *sql.DB) error {
	schemas := []string{
		// Processes table
		`CREATE TABLE IF NOT EXISTS processes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			time TEXT,
			pid TEXT,
			comm TEXT,
			args TEXT
		);`,

		// Network connections table
		`CREATE TABLE IF NOT EXISTS network_connections (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			pid TEXT,
			comm TEXT,
			ip_version TEXT,
			source_addr TEXT,
			source_port TEXT,
			dest_addr TEXT,
			dest_port TEXT
		);`,

		// Disk latency table
		`CREATE TABLE IF NOT EXISTS disk_latency (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			range_min INTEGER,
			range_max INTEGER,
			count INTEGER
		);`,
	}

	for _, schema := range schemas {
		if _, err := db.Exec(schema); err != nil {
			return err
		}
	}

	return nil
}
