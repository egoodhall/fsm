# Database Management

This project uses:
- [goose](https://github.com/pressly/goose) for database migrations
- [sqlc](https://sqlc.dev) for type-safe SQL queries

## Prerequisites

```bash
# Install goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# Install sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

## Directory Structure

- `migrations/`: Contains all database migration files
- `queries/`: Contains SQLC query definitions
- `gen/`: Contains generated Go code from SQLC

## Database Operations

### Migrations

```bash
# Create a new migration
make migrate-create

# Run migrations up
make migrate-up

# Roll back migrations
make migrate-down
```

### Generate SQL Code

```bash
# Generate Go code from SQL queries
make sqlc
```

## Configuration

The database connection string can be configured by setting the `DB_URL` environment variable:

```bash
export DB_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable"
```

Default connection string: `postgres://postgres:postgres@localhost:5432/fsm?sslmode=disable` 
