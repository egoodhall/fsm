package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/egoodhall/fsm"
)

const (
	StateA fsm.State = "a"
	StateB fsm.State = "b"
)

type TaskState struct {
	States []fsm.State
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	machine, err := fsm.New[string, TaskState]("example").
		InitialState(func(ctx context.Context, req fsm.FirstInput[string]) (fsm.Output, error) {
			return fsm.GoTo(StateA, TaskState{States: make([]fsm.State, 0)})
		}).
		AddState(StateA, func(ctx context.Context, req fsm.NthInput[string, TaskState]) (fsm.Output, error) {
			if len(req.Previous().States) == 0 {
				return fsm.GoTo(StateB, TaskState{States: append(req.Previous().States, StateA)})
			}
			req.Logger().Info("Exiting from state A", "id", req.ID(), "states", req.Previous().States)
			return fsm.Done()
		}).
		AddState(StateB, func(ctx context.Context, req fsm.NthInput[string, TaskState]) (fsm.Output, error) {
			return fsm.GoTo(StateA, TaskState{States: append(req.Previous().States, StateB)})
		}).
		Build(ctx, fsm.Logger[string, TaskState](slog.Default()), fsm.DB[string, TaskState]("fsm.db"))

	if err != nil {

		slog.Error("Failed to build machine", "error", err)

	}

	slog.Info("Starting machine")
	submit, stop, err := machine.Start(ctx)
	if err != nil {
		slog.Error("Failed to start machine", "error", err)
	}
	defer stop()

	// Automatically assigned ID 1
	if id, err := submit(ctx, 0, "a"); err != nil {
		slog.Error("Failed to submit task", "id", id, "error", err)
	}

	// Automatically assigned ID 2
	if id, err := submit(ctx, 0, "a"); err != nil {
		slog.Error("Failed to submit task", "id", id, "error", err)
	}

	// Already assigned, this will fail
	if id, err := submit(ctx, 1, "b"); err != nil {
		slog.Error("Failed to submit task", "id", id, "error", err)
	}

	<-ctx.Done()
}
