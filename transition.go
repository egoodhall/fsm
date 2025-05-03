package fsm

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log/slog"
)

// FirstTransition is a function that transitions the task from the initial state to a new state.
type FirstTransition[IN any, OUT any] func(ctx context.Context, req FirstInput[IN]) (Output, error)

// NthTransition is a function that transitions the task from a non-initial state to a new state.
type NthTransition[IN any, OUT any] func(ctx context.Context, req NthInput[IN, OUT]) (Output, error)

// FirstInput is an input for the initial state transition.
type FirstInput[IN any] interface {
	ID() int64
	Input() IN
	Logger() *slog.Logger
}

var _ FirstInput[any] = new(firstInput[any])

type firstInput[IN any] struct {
	id     int64
	event  IN
	logger *slog.Logger
}

func (r *firstInput[IN]) ID() int64 {
	return r.id
}

func (r *firstInput[IN]) Input() IN {
	return r.event
}

func (r *firstInput[IN]) Logger() *slog.Logger {
	return r.logger
}

// NthInput is an input for state transitions other than the initial state.
type NthInput[IN any, OUT any] interface {
	FirstInput[IN]
	// The output from the previous state transition.
	Previous() OUT
}

var _ NthInput[any, any] = new(nthInput[any, any])

type nthInput[IN any, OUT any] struct {
	id     int64
	event  IN
	output OUT
	logger *slog.Logger
}

func (r *nthInput[IN, OUT]) ID() int64 {
	return r.id
}

func (r *nthInput[IN, OUT]) Input() IN {
	return r.event
}

func (r *nthInput[IN, OUT]) Previous() OUT {
	return r.output
}

func (r *nthInput[IN, OUT]) Logger() *slog.Logger {
	return r.logger
}

// Output is the output from a state transition.
type Output interface {
	// The next state that this task should transition to.
	NextState() State
	// The data that should be passed to the next state.
	Data() []byte
}

func GoTo[T any](state State, data T) (Output, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(data); err != nil && err.Error() == "gob: cannot encode nil value" {
		return &output{state: state, data: buf.Bytes()}, nil
	} else if err != nil {
		return nil, fmt.Errorf("marshal output: %w", err)
	}
	return &output{state: state, data: buf.Bytes()}, nil
}

func Done() (Output, error) {
	return GoTo[any](StateDone, nil)
}

func Error[T any](err error) (Output, error) {
	return &output{state: StateError, data: []byte(err.Error())}, nil
}

var _ Output = new(output)

type output struct {
	state State
	data  []byte
}

func (r *output) NextState() State {
	return r.state
}

func (r *output) Data() []byte {
	return r.data
}
