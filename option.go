package fsm

import (
	"context"
	"database/sql"
	"log/slog"
)

// Options for the state machine.
type Option[IN any, OUT any] func(context.Context, *fsm[IN, OUT]) error

// Logger sets the logger for the state machine.
func Logger[IN, OUT any](logger *slog.Logger) Option[IN, OUT] {
	return func(ctx context.Context, f *fsm[IN, OUT]) error {
		f.logger = logger
		return nil
	}
}

// DB configures a new database connection for the state machine.
func DB[IN, OUT any](dbFile string) Option[IN, OUT] {
	return func(ctx context.Context, f *fsm[IN, OUT]) error {
		var err error
		if f.db, err = initDB(ctx, dbFile); err != nil {
			return err
		}
		if f.id, err = f.db.CreateStateMachine(ctx, f.name); err != nil {
			return err
		}
		return nil
	}
}

func MemDB[IN, OUT any]() Option[IN, OUT] {
	return DB[IN, OUT]("file:fsm.db?cache=shared&mode=memory")
}

// WithDB configures the existing database connection for the state machine.
func WithDB[IN, OUT any](db *sql.DB) Option[IN, OUT] {
	return func(ctx context.Context, f *fsm[IN, OUT]) error {
		var err error
		if f.db, err = setupDB(ctx, db); err != nil {
			return err
		}
		if f.id, err = f.db.CreateStateMachine(ctx, f.name); err != nil {
			return err
		}
		return nil
	}
}

// OnTransitionFunc is a function that is called when a transition successfully completes.
type OnTransitionFunc func(ctx context.Context, id int64, from, to State)

// OnTransition sets the function to call when a transition successfully completes.
func OnTransition[IN, OUT any](fn OnTransitionFunc) Option[IN, OUT] {
	return func(ctx context.Context, f *fsm[IN, OUT]) error {
		f.onTransition = fn
		return nil
	}
}
