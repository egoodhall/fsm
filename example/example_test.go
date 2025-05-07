package example_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/egoodhall/fsm"
	"github.com/egoodhall/fsm/example"
)

func TestMultistepFSM(t *testing.T) {
	completions := make(chan fsm.TaskID, 3)
	defer close(completions)
	transitions := make(chan fsm.TaskID, 6)
	defer close(transitions)

	// Print debug messages
	slog.SetLogLoggerLevel(slog.LevelDebug)

	f, err := example.NewCreateWorkspaceFSMBuilder().
		CreateRecordState(func(ctx context.Context, transitions example.CreateRecordTransitions, c example.WorkspaceContext) error {
			slog.Info("create record", "attempt", fsm.GetAttempt(ctx))
			if fsm.GetAttempt(ctx) == 0 {
				return errors.New("first attempt")
			}
			return transitions.ToCloneRepo(ctx, c, example.WorkspaceID(1))
		}).
		CloneRepoState(func(ctx context.Context, transitions example.CloneRepoTransitions, c example.WorkspaceContext, i example.WorkspaceID) error {
			return transitions.ToDone(ctx)
		}).
		BuildAndStart(t.Context(),
			fsm.WithLogger(slog.Default()),
			fsm.InMemory(),
			fsm.WithTransitionListener(func(ctx context.Context, id fsm.TaskID, from fsm.State, to fsm.State) {
				transitions <- id
			}),
			fsm.WithCompletionListener(func(ctx context.Context, id fsm.TaskID, state fsm.State) {
				completions <- id
			}),
		)
	if err != nil {
		t.Fatal(err)
	}

	ids := make([]fsm.TaskID, 3)

	ids[0], err = f.Submit(t.Context(), example.WorkspaceContext{})
	if err != nil {
		t.Fatal(err)
	}

	ids[1], err = f.Submit(t.Context(), example.WorkspaceContext{})
	if err != nil {
		t.Fatal(err)
	}

	ids[2], err = f.Submit(t.Context(), example.WorkspaceContext{})
	if err != nil {
		t.Fatal(err)
	}

	for {
		select {
		case id := <-transitions:
			t.Log("transitioned", "id", id)
		case id := <-completions:
			t.Log("completed", "id", id)
		case <-time.After(500 * time.Millisecond):
			t.Log("timeout")
			return
		}
	}
}
