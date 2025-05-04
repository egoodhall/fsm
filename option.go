package fsm

import (
	"context"
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

// TransitionListener is a function that is called when a transition successfully completes.
type TransitionListener func(ctx context.Context, id int64, from, to State)

// OnTransition sets the function to call when a transition successfully completes.
func OnTransition[IN, OUT any](fn TransitionListener) Option[IN, OUT] {
	return func(ctx context.Context, f *fsm[IN, OUT]) error {
		f.onTransition = fn
		return nil
	}
}
