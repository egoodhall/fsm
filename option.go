package fsm

import (
	"context"
	"log/slog"
)

type TransitionListener func(ctx context.Context, id TaskID, from State, to State, inputs ...any)
type CompletionListener func(ctx context.Context, id TaskID, state State)

type SupportsOptions interface {
	WithContext(update func(ctx context.Context) context.Context)
	WithStore(store Store)
	WithTransitionListener(listener TransitionListener)
	WithCompletionListener(listener CompletionListener)
}

type Option func(SupportsOptions) error

func WithLogger(logger *slog.Logger) Option {
	return func(s SupportsOptions) error {
		s.WithContext(func(ctx context.Context) context.Context {
			return PutLogger(ctx, logger)
		})
		return nil
	}
}

func WithTransitionListener(listener TransitionListener) Option {
	return func(s SupportsOptions) error {
		s.WithTransitionListener(listener)
		return nil
	}
}

func WithCompletionListener(listener CompletionListener) Option {
	return func(s SupportsOptions) error {
		s.WithCompletionListener(listener)
		return nil
	}
}
