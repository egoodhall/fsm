package fsm

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/egoodhall/fsm/gen/sqlc"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Store interface {
	DB() *sql.DB
	Q() sqlc.Querier
}

var _ Store = &store{}

type store struct {
	db *sql.DB
}

func (p *store) DB() *sql.DB {
	return p.db
}

func (p *store) Q() sqlc.Querier {
	return sqlc.New(p.db)
}

func OnDisk(path string) (Store, error) {
	db, err := initDB(context.Background(), path)
	if err != nil {
		return nil, err
	}
	return &store{db: db}, nil
}

func InMemory() (Store, error) {
	return OnDisk("file:fsm.db?mode=memory&cache=shared")
}

// InitDB initializes a new SQLite database connection, runs migrations,
// and returns a querier for interacting with the database.
func initDB(ctx context.Context, dbPath string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("create data directory: %w", err)
	}

	// Open the SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	return setupDB(ctx, db)
}

// setupDB runs migrations and returns a querier for interacting with the database.
func setupDB(ctx context.Context, db *sql.DB) (*sql.DB, error) {

	mdir, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return nil, fmt.Errorf("create migrations subdirectory: %w", err)
	}

	provider, err := goose.NewProvider(goose.DialectSQLite3, db, mdir, goose.WithLogger(goose.NopLogger()))
	if err != nil {
		return nil, fmt.Errorf("create goose provider: %w", err)
	}

	// Run migrations
	if _, err := provider.Up(ctx); err != nil {
		// Check if the error is just that migrations are up-to-date
		if !strings.Contains(err.Error(), "no change") {
			return nil, fmt.Errorf("run migrations: %w", err)
		}
	}

	// Create and return the querier
	return db, nil
}
