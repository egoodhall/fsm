package main

import (
	"context"

	"github.com/egoodhall/fsm"
)

//go:generate go run ../cmd/fsmgen -out generated -pkg generated create_workspace.yaml

type WorkspaceContext struct {
	RepositoryURL string
	BranchName    string
}

func main() {
	fsm := generated.NewCreateWorkspaceFSM().
		InitialState(func(ctx context.Context, transitions generated.InitialTransitions, i WorkspaceContext) (fsm.Transition, error) {
			return nil, nil
		}).
		CreatingRecordState(func(ctx context.Context, transitions generated.CreatingRecordTransitions, i WorkspaceContext) (fsm.Transition, error) {
			return nil, nil
		}).
		CloningRepositoryState(func(ctx context.Context, transitions generated.CloningRepositoryTransitions, i WorkspaceContext) (fsm.Transition, error) {
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
}
