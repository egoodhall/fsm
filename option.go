package fsm

import (
	"context"
	"log/slog"
)

func Logger[IN any, OUT any](logger *slog.Logger) Option[IN, OUT] {
	return func(ctx context.Context, f *fsm[IN, OUT]) error {
		f.logger = logger
		return nil
	}
}

func DB[IN any, OUT any](dbFile string) Option[IN, OUT] {
	return func(ctx context.Context, f *fsm[IN, OUT]) (err error) {
		if f.db, err = initDB(ctx, dbFile); err != nil {
			return
		}
		if f.id, err = f.db.CreateStateMachine(ctx, f.name); err != nil {
			return
		}
		return
	}
}
