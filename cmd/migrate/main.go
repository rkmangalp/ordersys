package main

import (
	"database/sql"
	"embed"
	"log"
	"os"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: migrate [up|down]")
	}

	// Database connection
	dsn := "root:@tcp(localhost:3307)/ordersys?parseTime=true&multiStatements=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Ensure database exists
	if _, err := db.Exec(`CREATE DATABASE IF NOT EXISTS ordersys`); err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	// Create migrations tracking table
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
		)
	`); err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	switch os.Args[1] {
	case "up":
		runMigrationsUp(db)
	case "down":
		runMigrationsDown(db)
	default:
		log.Fatalf("Unknown command: %s", os.Args[1])
	}
}

func runMigrationsUp(db *sql.DB) {
	// Read all .up.sql files from embedded migrations
	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		log.Fatalf("Failed to read migrations dir: %v", err)
	}

	var upFiles []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".up.sql") {
			upFiles = append(upFiles, e.Name())
		}
	}
	sort.Strings(upFiles)

	for _, file := range upFiles {
		version := strings.TrimSuffix(file, ".up.sql")

		// Check if already applied
		var exists bool
		err := db.QueryRow(`SELECT 1 FROM schema_migrations WHERE version = ?`, version).Scan(&exists)
		if err == sql.ErrNoRows {
			// Not applied, run it
			content, err := migrationFiles.ReadFile("migrations/" + file)
			if err != nil {
				log.Fatalf("Failed to read migration %s: %v", file, err)
			}

			// Split by semicolon and execute each statement
			statements := splitStatements(string(content))
			tx, err := db.Begin()
			if err != nil {
				log.Fatalf("Failed to begin transaction for %s: %v", file, err)
			}

			for _, stmt := range statements {
				stmt = strings.TrimSpace(stmt)
				if stmt == "" {
					continue
				}
				if _, err := tx.Exec(stmt); err != nil {
					tx.Rollback()
					log.Fatalf("Failed to execute statement in %s: %v\nStatement: %s", file, err, stmt)
				}
			}

			// Record migration
			if _, err := tx.Exec(`INSERT INTO schema_migrations (version) VALUES (?)`, version); err != nil {
				tx.Rollback()
				log.Fatalf("Failed to record migration %s: %v", version, err)
			}

			if err := tx.Commit(); err != nil {
				log.Fatalf("Failed to commit migration %s: %v", version, err)
			}

			log.Printf("Applied migration: %s", version)
		} else if err != nil {
			log.Fatalf("Failed to check migration status: %v", err)
		} else {
			log.Printf("Skipping already applied: %s", version)
		}
	}
	log.Println("All migrations applied successfully")
}

func runMigrationsDown(db *sql.DB) {
	// Get applied migrations in reverse order
	rows, err := db.Query(`SELECT version FROM schema_migrations ORDER BY applied_at DESC`)
	if err != nil {
		log.Fatalf("Failed to query migrations: %v", err)
	}
	defer rows.Close()

	var versions []string
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			log.Fatalf("Failed to scan version: %v", err)
		}
		versions = append(versions, v)
	}

	for _, version := range versions {
		downFile := version + ".down.sql"
		content, err := migrationFiles.ReadFile("migrations/" + downFile)
		if err != nil {
			log.Fatalf("Failed to read down migration %s: %v", downFile, err)
		}

		statements := splitStatements(string(content))
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("Failed to begin transaction for %s: %v", downFile, err)
		}

		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := tx.Exec(stmt); err != nil {
				tx.Rollback()
				log.Fatalf("Failed to execute statement in %s: %v\nStatement: %s", downFile, err, stmt)
			}
		}

		if _, err := tx.Exec(`DELETE FROM schema_migrations WHERE version = ?`, version); err != nil {
			tx.Rollback()
			log.Fatalf("Failed to remove migration record %s: %v", version, err)
		}

		if err := tx.Commit(); err != nil {
			log.Fatalf("Failed to commit down migration %s: %v", version, err)
		}

		log.Printf("Rolled back migration: %s", version)
	}
	log.Println("All migrations rolled back successfully")
}

func splitStatements(sql string) []string {
	// Simple split by semicolon - works for our migration files
	var statements []string
	var current strings.Builder
	inString := false
	stringChar := byte(0)

	for i := 0; i < len(sql); i++ {
		c := sql[i]
		current.WriteByte(c)

		if (c == '\'' || c == '"' || c == '`') && (i == 0 || sql[i-1] != '\\') {
			if !inString {
				inString = true
				stringChar = c
			} else if c == stringChar {
				inString = false
			}
		}

		if c == ';' && !inString {
			statements = append(statements, current.String())
			current.Reset()
		}
	}

	if current.Len() > 0 {
		statements = append(statements, current.String())
	}

	return statements
}