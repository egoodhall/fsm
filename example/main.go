package main

import (
	"context"
	"log"

	"github.com/egoodhall/fsm"
)

//go:generate go run ../cmd/fsmgen -out . -pkg main create_workspace.yaml

type WorkspaceContext struct {
	RepositoryURL string
	BranchName    string
}

func main() {
	fsm := NewCreateWorkspaceFSM().
		InitialState(func(ctx context.Context, transitions InitialTransitions, i WorkspaceContext) (fsm.Transition, error) {
			return nil, nil
		}).
		CreatingRecordState(func(ctx context.Context, transitions CreatingRecordTransitions, i WorkspaceContext) (fsm.Transition, error) {
			return nil, nil
		}).
		CloningRepositoryState(func(ctx context.Context, transitions CloningRepositoryTransitions, i WorkspaceContext) (fsm.Transition, error) {
			return nil, nil
		}).
		DoneState(func(ctx context.Context, i WorkspaceContext) (fsm.Transition, error) {
			return nil, nil
		}).
		ErrorState(func(ctx context.Context, i WorkspaceContext) (fsm.Transition, error) {
			return nil, nil
		}).
		Build()

	id, err := fsm.SubmitInitial(context.Background(), WorkspaceContext{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(id)
}
