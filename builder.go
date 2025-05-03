package fsm

import (
	"context"
	"log/slog"
)

var _ InitialStateBuilder[any, any] = new(fsm[any, any])
var _ NthStateBuilder[any, any] = new(fsm[any, any])

func (f *fsm[IN, OUT]) InitialState(transition FirstTransition[IN, OUT]) NthStateBuilder[IN, OUT] {
	f.first = transition
	return f
}

func (f *fsm[IN, OUT]) AddState(state State, transition NthTransition[IN, OUT]) NthStateBuilder[IN, OUT] {
	if f.rest == nil {
		f.rest = make(map[State]NthTransition[IN, OUT])
	}
	f.rest[state] = transition
	return f
}

func (f *fsm[IN, OUT]) Build(ctx context.Context, opts ...Option[IN, OUT]) (FSM[IN, OUT], error) {
	f.logger = slog.Default()

	for _, opt := range opts {
		if err := opt(ctx, f); err != nil {
			return nil, err
		}
	}

	// If the database is not set, set it
	if f.db == nil {
		if err := DB[IN, OUT]("fsm.db")(ctx, f); err != nil {
			return nil, err
		}
	}

	return f, nil
}
