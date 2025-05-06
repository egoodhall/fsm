package fsm

import "context"

type State string

type stateKey struct{}

func PutState(ctx context.Context, id State) context.Context {
	return context.WithValue(ctx, stateKey{}, id)
}

func GetState(ctx context.Context) State {
	id, ok := ctx.Value(stateKey{}).(State)
	if !ok {
		panic("state not found")
	}
	return id
}

type TaskID int64

type taskIDKey struct{}

func PutTaskID(ctx context.Context, id TaskID) context.Context {
	return context.WithValue(ctx, taskIDKey{}, id)
}

func GetTaskID(ctx context.Context) TaskID {
	id, ok := ctx.Value(taskIDKey{}).(TaskID)
	if !ok {
		panic("task ID not found")
	}
	return id
}
