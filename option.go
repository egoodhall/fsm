package fsm

import "log/slog"

type TransitionListener func(id TaskID, from State, to State, inputs ...any)
type CompletionListener func(id TaskID, state State)

type SupportsOptions interface {
	WithStore(store Store)
	WithLogger(logger *slog.Logger)
	WithTransitionListener(listener TransitionListener)
	WithCompletionListener(listener CompletionListener)
}

type Option func(SupportsOptions) error

func WithLogger(logger *slog.Logger) Option {
	return func(s SupportsOptions) error {
		s.WithLogger(logger)
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
