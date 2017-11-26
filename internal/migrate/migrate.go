// Package migrate performs simple database migrations
package migrate

import (
	"database/sql"
	"fmt"

	"github.com/derWhity/micasa/internal/log"
	"github.com/jmoiron/sqlx"
)

var migrations []dbMigration

type dbMigration struct {
	Version uint
	Queries []string
}

// Execute runs the current DB migration on the given database
func (mig *dbMigration) Execute(db *sqlx.DB, log log.Logger) error {
	// Check if the migration has already run
	query := `SELECT success FROM Migrations WHERE version = $1`
	var success = false
	err := db.QueryRow(query, mig.Version).Scan(&success)
	if err != nil {
		switch {
		case err != sql.ErrNoRows:
			log.Error("Failed to fetch version information", err)
			return err
		}
	}
	if !success {
		// We need to execute this migration
		log.Info(fmt.Sprintf("Executing DB migration #%d", mig.Version))
		for i, query := range mig.Queries {
			log.Info(fmt.Sprintf("Query %d of %d...", (i + 1), len(mig.Queries)))
			if _, err := db.Exec(query); err != nil {
				log.Error(fmt.Sprintf("Query #%d failed", (i+1)), err)
				db.Exec(`REPLACE INTO Migrations(version, success) VALUES($1, 0)`, mig.Version)
				return err
			}
		}
		// Queries executed successfully - save our status
		db.Exec(`REPLACE INTO Migrations(version, success) VALUES($1, 1)`, mig.Version)
	}
	return nil
}

// ExecuteMigrationsOnDb executes the database migrations on the given database instance
func ExecuteMigrationsOnDb(db *sqlx.DB, log log.Logger) error {
	// Create the migrations table if it does not exist, yet
	query := `CREATE TABLE IF NOT EXISTS Migrations (
                version   INTEGER NOT NULL,
                success   INTEGER NOT NULL DEFAULT 0,
                PRIMARY KEY(version)
            )`
	if _, err := db.Exec(query); err != nil {
		log.Error("Failed to create migrations table", err)
		return err
	}
	for _, mig := range migrations {
		if err := mig.Execute(db, log); err != nil {
			log.Error(fmt.Sprintf("Failed to execute migration #%d", mig.Version), err)
			return err
		}
	}
	return nil
}

// For now, the migrations are part of the package...
func init() {
	migrations = []dbMigration{
		{
			Version: 1,
			Queries: []string{
				`CREATE TABLE Users (
					userid	VARCHAR(32) NOT NULL,
					name	VARCHAR(64) NOT NULL UNIQUE ON CONFLICT ABORT,
					passwordHash	VARCHAR(128) NOT NULL DEFAULT '',
					fullName	VARCHAR(128) NOT NULL DEFAULT '',
					createdAt	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updatedAt	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
					PRIMARY KEY(userid)
				);`,
			},
		},
	}
}
