package db

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
)

func RunMigrations() error {
	// Specify the migration directory and SQLite database path
	dbPath := "pkg/db/database.db"
	migrationDir := "pkg/db/migrations/sqlite"

	// Run migrations using golang-migrate
	return migrateDB("sqlite3", dbPath, migrationDir)
}

func migrateDB(driver, dbPath, migrationDir string) error {
	// Open a database connection
	db, err := sql.Open(driver, dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Run migrations
	instance, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return err
	}

	srcDriver, err := (&file.File{}).Open(migrationDir)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("file", srcDriver, "sqlite3", instance)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Migration error: %v", err)
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}
