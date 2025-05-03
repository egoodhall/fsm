package fsm

import (
	"context"
	"log/slog"
)

// Represents a state in the state machine. This identifies
// where a given task is in the state machine. The default
// states are:
// - __initial__: The task has been submitted but not yet processed
// - __done__: The task has completed successfully
// - __error__: A transition failed
type State string

const (
	StateInitial State = "__initial__"
	StateDone    State = "__done__"
	StateError   State = "__error__"
)

// Options for the state machine.
type Option[IN any, OUT any] func(context.Context, *fsm[IN, OUT]) error

// FSM is a finite state machine. It tracks the state of a task
// and transitions it between states based on the input. State is
// persisted in a SQLite database.
type FSM[IN any, OUT any] interface {
	Start(ctx context.Context) (SubmitFunc[IN], context.CancelFunc, error)
}

// InitialStateBuilder is a builder for the initial state transition of the state machine.
type InitialStateBuilder[IN any, OUT any] interface {
	InitialState(transition FirstTransition[IN, OUT]) NthStateBuilder[IN, OUT]
}

// NthStateBuilder is a builder for state transitions other than the initial state.
type NthStateBuilder[IN any, OUT any] interface {
	AddState(state State, transition NthTransition[IN, OUT]) NthStateBuilder[IN, OUT]
	Build(ctx context.Context, opts ...Option[IN, OUT]) (FSM[IN, OUT], error)
}

// FirstInput is an input for the initial state transition.
type FirstInput[IN any] interface {
	ID() int64
	Input() IN
	Logger() *slog.Logger
}

// NthInput is an input for state transitions other than the initial state.
type NthInput[IN any, OUT any] interface {
	FirstInput[IN]
	// The output from the previous state transition.
	Previous() OUT
}

// Output is the output from a state transition.
type Output interface {
	// The next state that this task should transition to.
	NextState() State
	// The data that should be passed to the next state.
	Data() []byte
}

// SubmitFunc is a function that submits a task to the state machine.
type SubmitFunc[IN any] func(ctx context.Context, id int64, event IN) (int64, error)

// FirstTransition is a function that transitions the task from the initial state to a new state.
type FirstTransition[IN any, OUT any] func(ctx context.Context, req FirstInput[IN]) (Output, error)

// NthTransition is a function that transitions the task from a non-initial state to a new state.
type NthTransition[IN any, OUT any] func(ctx context.Context, req NthInput[IN, OUT]) (Output, error)

// New creates a new state machine
// The name is used to identify the state machine in the database.
func New[IN any, OUT any](name string) InitialStateBuilder[IN, OUT] {
	return &fsm[IN, OUT]{
		name: name,
		rest: make(map[State]NthTransition[IN, OUT]),
	}
}
