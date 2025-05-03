package fsm

import (
	"context"
)

type State string

const (
	StateInitial State = "__initial__"
	StateDone    State = "__done__"
	StateError   State = "__error__"
)

type Option[IN any, OUT any] func(context.Context, *fsm[IN, OUT]) error

type FSM[IN any, OUT any] interface {
	Start(ctx context.Context) (SubmitFunc[IN], context.CancelFunc, error)
}

type InitialStateBuilder[IN any, OUT any] interface {
	InitialState(transition FirstTransition[IN, OUT]) NthStateBuilder[IN, OUT]
}

type NthStateBuilder[IN any, OUT any] interface {
	AddState(state State, transition NthTransition[IN, OUT]) NthStateBuilder[IN, OUT]
	Build(ctx context.Context, opts ...Option[IN, OUT]) (FSM[IN, OUT], error)
}

type FirstInput[IN any] interface {
	ID() string
	Input() IN
}

type NthInput[IN any, OUT any] interface {
	FirstInput[IN]
	Previous() OUT
}

type Output interface {
	NextState() State
	Data() []byte
}

type SubmitFunc[IN any] func(ctx context.Context, id string, event IN) error

type FirstTransition[IN any, OUT any] func(ctx context.Context, req FirstInput[IN]) (Output, error)

type NthTransition[IN any, OUT any] func(ctx context.Context, req NthInput[IN, OUT]) (Output, error)

func New[IN any, OUT any](name string) InitialStateBuilder[IN, OUT] {
	return &fsm[IN, OUT]{
		name: name,
		rest: make(map[State]NthTransition[IN, OUT]),
	}
}
