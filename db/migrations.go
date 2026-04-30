package db

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlserver"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// createMigrator creates a migrate instance based on database type
func createMigrator(db *sql.DB, config *Config, migrationsFS embed.FS, migrationsPath string) (*migrate.Migrate, error) {
	// Create source from embedded filesystem
	source, err := iofs.New(migrationsFS, migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration source: %w", err)
	}

	// Create database driver based on type
	var driver database.Driver

	switch config.Type {
	case Postgres:
		driver, err = postgres.WithInstance(db, &postgres.Config{})
	case MySQL:
		driver, err = mysql.WithInstance(db, &mysql.Config{})
	case MSSQL:
		driver, err = sqlserver.WithInstance(db, &sqlserver.Config{})
	default:
		return nil, fmt.Errorf("unsupported database type for migrations: %s", config.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create database driver: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithInstance("iofs", source, config.Database, driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return m, nil
}

// MigrationVersion returns the current migration version
func (c *Client) MigrationVersion() (uint, bool, error) {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return 0, false, err
	}

	// This is a simplified version check
	// In production, you'd want to use the migrate instance
	var version uint
	var dirty bool

	err = sqlDB.QueryRow("SELECT version, dirty FROM schema_migrations LIMIT 1").Scan(&version, &dirty)
	if err != nil {
		return 0, false, err
	}

	return version, dirty, nil
}
