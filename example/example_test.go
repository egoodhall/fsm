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

	f, err := example.NewTestMachineFSMBuilder().
		FromState1(func(ctx context.Context, transitions example.State1Transitions, c int) error {
			if fsm.GetAttempt(ctx) == 0 {
				return errors.New("first attempt")
			}
			return transitions.ToState2(ctx, c)
		}).
		FromState2(func(ctx context.Context, transitions example.State2Transitions, c int) error {
			return transitions.ToDone(ctx)
		}).
		BuildAndStart(t.Context(),
			fsm.WithLogger(slog.Default()),
			fsm.WithStore(fsm.InMemory()),
			fsm.WithBackoff(fsm.ExponentialBackoff(10*time.Millisecond, 1*time.Second)),
			fsm.WithTransitionListener(func(ctx context.Context, id fsm.TaskID, from fsm.State, to fsm.State) {
				fsm.Logger(ctx).Info("transitioned", "id", id, "from", from, "to", to, "attempt", fsm.GetAttempt(ctx))
			}),
			fsm.WithCompletionListener(func(ctx context.Context, id fsm.TaskID, state fsm.State) {
				fsm.Logger(ctx).Info("completed", "id", id, "state", state, "attempt", fsm.GetAttempt(ctx))
				completions <- id
			}),
		)
	if err != nil {
		t.Fatal(err)
	}

	ids := make([]fsm.TaskID, 3)

	ids[0], err = f.Submit(t.Context(), 0)
	if err != nil {
		t.Fatal(err)
	}

	ids[1], err = f.Submit(t.Context(), 1)
	if err != nil {
		t.Fatal(err)
	}

	ids[2], err = f.Submit(t.Context(), 2)
	if err != nil {
		t.Fatal(err)
	}

	completed := make([]bool, 3)
	for {
		select {
		case id := <-completions:
			t.Log("completed", "id", id)
			completed[id-1] = true
		case <-time.After(500 * time.Millisecond):
			for i, id := range ids {
				if !completed[i] {
					t.Fatal("timeout", "id", id)
				}
			}
			return
		}
	}
}
