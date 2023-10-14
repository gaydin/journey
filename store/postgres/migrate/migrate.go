package migrate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gaydin/journey/store/postgres/migrate/migrations"
)

const (
	migrationTableCreate = `
CREATE TABLE IF NOT EXISTS migrations (
 name VARCHAR(255)
,UNIQUE(name)
)
`

	migrationInsert = `
INSERT INTO migrations (name) VALUES ($1)
`

	migrationSelect = `
SELECT name FROM migrations
`
)

// Migrate performs the database migration. If the migration fails
// and error is returned.
func Migrate(db *pgxpool.Pool) error {
	if err := createTable(db); err != nil {
		return fmt.Errorf("migrate error createTable: %w", err)
	}

	completed, err := selectCompleted(db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("migrate error selectCompleted: %w", err)
	}

	statements, err := getStatements(migrations.FS)
	if err != nil {
		return fmt.Errorf("migrate error getstatements: %w", err)
	}

	for _, statement := range statements {
		if _, ok := completed[statement.Name]; ok {
			continue
		}

		if _, err := db.Exec(context.TODO(), statement.Value); err != nil {
			return fmt.Errorf("statement %s exec error: %w", statement.Name, err)
		}

		if err := insertMigration(db, statement.Name); err != nil {
			return err
		}
	}

	return nil
}

func createTable(db *pgxpool.Pool) error {
	_, err := db.Exec(context.TODO(), migrationTableCreate)
	return err
}

func insertMigration(db *pgxpool.Pool, name string) error {
	_, err := db.Exec(context.TODO(), migrationInsert, name)
	return err
}

func selectCompleted(db *pgxpool.Pool) (map[string]struct{}, error) {
	completedMigrations := map[string]struct{}{}
	rows, err := db.Query(context.TODO(), migrationSelect)
	if err != nil {
		return nil, fmt.Errorf("query migrationSelect %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string

		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("rows.Scan %w", err)
		}

		completedMigrations[name] = struct{}{}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err %w", err)
	}

	return completedMigrations, nil
}
