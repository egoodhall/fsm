package fsm

import (
	"context"
	"log/slog"
)

type stateKey struct{}

func PutState(ctx context.Context, id State) context.Context {
	return context.WithValue(ctx, stateKey{}, id)
}

func GetState(ctx context.Context) State {
	id, ok := ctx.Value(stateKey{}).(State)
	if !ok {
		return State("UNKNOWN")
	}
	return id
}

type taskIDKey struct{}

func PutTaskID(ctx context.Context, id TaskID) context.Context {
	return context.WithValue(ctx, taskIDKey{}, id)
}

func GetTaskID(ctx context.Context) TaskID {
	id, ok := ctx.Value(taskIDKey{}).(TaskID)
	if !ok {
		return TaskID(-1)
	}
	return id
}

type loggerKey struct{}

func PutLogger(ctx context.Context, logger *slog.Logger) context.Context {
	if logger == nil {
		return ctx
	}
	return context.WithValue(ctx, loggerKey{}, logger)
}

func Logger(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerKey{}).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}

type attemptKey struct{}

func PutAttempt(ctx context.Context, attempt int) context.Context {
	return context.WithValue(ctx, attemptKey{}, attempt)
}

func GetAttempt(ctx context.Context) int {
	attempt, ok := ctx.Value(attemptKey{}).(int)
	if !ok {
		return 0
	}
	return attempt
}
