package fsm

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/egoodhall/fsm/gen/sqlc"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

func init() {
	// Register the embedded migrations
	goose.SetBaseFS(migrations)

	// Force SQLite dialect for goose
	goose.SetDialect("sqlite3")

	goose.SetLogger(goose.NopLogger())
}

//go:embed migrations/*.sql
var migrations embed.FS

// InitDB initializes a new SQLite database connection, runs migrations,
// and returns a querier for interacting with the database.
func initDB(ctx context.Context, dbPath string) (sqlc.Querier, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("create data directory: %w", err)
	}

	// Open the SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	return setupDB(db)
}

// setupDB runs migrations and returns a querier for interacting with the database.
func setupDB(db *sql.DB) (sqlc.Querier, error) {

	// Run migrations
	if err := goose.Up(db, "migrations"); err != nil {
		// Check if the error is just that migrations are up-to-date
		if !strings.Contains(err.Error(), "no change") {
			return nil, fmt.Errorf("run migrations: %w", err)
		}
	}

	// Create and return the querier
	return sqlc.New(db), nil
}
