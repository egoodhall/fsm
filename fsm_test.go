package fsm_test

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/egoodhall/fsm"
)

const (
	StateA fsm.State = "a"
	StateB fsm.State = "b"
)

type TaskState struct {
	States []fsm.State
}

func TestFSM(t *testing.T) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	store, err := fsm.InMemory()
	if err != nil {
		t.Fatalf("Failed to create in-memory store: %v", err)
	}

	machine, err := fsm.New[string, TaskState]("example", store).
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
		Build(ctx,
			// fsm.MemDB[string, TaskState](),
			fsm.Logger[string, TaskState](slog.Default()),
		)

	if err != nil {
		t.Fatalf("Failed to build machine: %v", err)
	}

	slog.Info("Starting machine")
	submit, stop, err := machine.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start machine: %v", err)
	}
	defer stop()

	// Automatically assigned ID 1
	if id, err := submit(ctx, 0, "a"); err != nil {
		t.Fatalf("Failed to submit task %d: %v", id, err)
	}

	// Automatically assigned ID 2
	if id, err := submit(ctx, 0, "a"); err != nil {
		t.Fatalf("Failed to submit task %d: %v", id, err)
	}

	// Already assigned, this will fail
	if id, err := submit(ctx, 1, "b"); err == nil {
		t.Fatalf("Expected error for task %d", id)
	}

	select {
	case <-ctx.Done():
	case <-time.After(1 * time.Second):
	}
}
