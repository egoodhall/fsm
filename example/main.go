package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/egoodhall/fsm"
)

//go:generate go run ../cmd/fsmgen -out . -pkg main create_workspace.yaml

type WorkspaceID int64

type WorkspaceContext struct {
	RepositoryURL string
	BranchName    string
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	fsm, err := NewCreateWorkspaceFSMBuilder().
		CreateRecordState(func(ctx context.Context, transitions CreateRecordTransitions, c WorkspaceContext) error {
			slog.Info("create record", "id", fsm.GetTaskID(ctx))
			return transitions.ToCloneRepo(ctx, c, WorkspaceID(1))
		}).
		CloneRepoState(func(ctx context.Context, transitions CloneRepoTransitions, c WorkspaceContext, i WorkspaceID) error {
			slog.Info("clone repo", "id", fsm.GetTaskID(ctx))
			return transitions.ToDone(ctx, c, i)
		}).
		DoneState(func(ctx context.Context, c WorkspaceContext, i WorkspaceID) error {
			slog.Info("done", "id", fsm.GetTaskID(ctx))
			return nil
		}).
		ErrorState(func(ctx context.Context, c WorkspaceContext) error {
			slog.Info("error", "id", fsm.GetTaskID(ctx))
			return nil
		}).
		BuildAndStart(ctx, fsm.WithLogger(slog.Default()), fsm.OnDisk("fsm.db"))
	if err != nil {
		log.Fatal(err)
	}

	id, err := fsm.Submit(ctx, WorkspaceContext{})
	if err != nil {
		log.Fatal(err)
	}

	id2, err := fsm.Submit(ctx, WorkspaceContext{})
	if err != nil {
		log.Fatal(err)
	}

	id3, err := fsm.Submit(ctx, WorkspaceContext{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(id, id2, id3)
	<-ctx.Done()
}
