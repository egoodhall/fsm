package example_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/egoodhall/fsm"
	"github.com/egoodhall/fsm/example"
)

func TestMultistepFSM(t *testing.T) {
	var completed int
	fsm, err := example.NewCreateWorkspaceFSMBuilder().
		CreateRecordState(func(ctx context.Context, transitions example.CreateRecordTransitions, c example.WorkspaceContext) error {
			return transitions.ToCloneRepo(ctx, c, example.WorkspaceID(1))
		}).
		CloneRepoState(func(ctx context.Context, transitions example.CloneRepoTransitions, c example.WorkspaceContext, i example.WorkspaceID) error {
			return transitions.ToDone(ctx, c, i)
		}).
		DoneState(func(ctx context.Context, c example.WorkspaceContext, i example.WorkspaceID) error {
			completed++
			return nil
		}).
		ErrorState(func(ctx context.Context, c example.WorkspaceContext) error {
			return nil
		}).
		BuildAndStart(t.Context(), fsm.WithLogger(slog.Default()), fsm.InMemory())
	if err != nil {
		t.Fatal(err)
	}

	_, err = fsm.Submit(t.Context(), example.WorkspaceContext{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = fsm.Submit(t.Context(), example.WorkspaceContext{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = fsm.Submit(t.Context(), example.WorkspaceContext{})
	if err != nil {
		t.Fatal(err)
	}

	<-time.After(500 * time.Millisecond)
	if completed != 3 {
		t.Fatalf("expected 3 completed, got %d", completed)
	}
}
