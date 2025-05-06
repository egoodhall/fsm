package fsm

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/egoodhall/fsm/gen/sqlc"
	"github.com/egoodhall/fsm/migrations"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type Q interface {
	sqlc.Querier
}

type Store interface {
	DB() *sql.DB
	Q() Q
}

var _ Store = &store{}

type store struct {
	db *sql.DB
}

func (p *store) DB() *sql.DB {
	return p.db
}

func (p *store) Q() Q {
	return sqlc.New(p.db)
}

func OnDisk(path string) Option {
	return func(options SupportsOptions) error {
		db, err := initDB(context.Background(), path)
		if err != nil {
			return err
		}
		options.WithStore(&store{db})
		return nil
	}
}

func InMemory() Option {
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

	provider, err := goose.NewProvider(goose.DialectSQLite3, db, migrations.FS, goose.WithLogger(goose.NopLogger()))
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
