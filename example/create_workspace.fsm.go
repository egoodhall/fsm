package main

import (
	"context"

	"github.com/egoodhall/fsm"
)

// Exported interfaces

type CreateWorkspaceFSM interface {
	SubmitInitial(ctx context.Context, i WorkspaceContext) (fsm.TaskID, error)
}

type CreateWorkspaceFSMBuilder_InitialStage interface {
	InitialState(func(ctx context.Context, transitions InitialTransitions, i WorkspaceContext) (fsm.Transition, error)) CreateWorkspaceFSMBuilder_CreatingRecordStage
}

type CreateWorkspaceFSMBuilder_CreatingRecordStage interface {
	CreatingRecordState(func(ctx context.Context, transitions CreatingRecordTransitions, i WorkspaceContext) (fsm.Transition, error)) CreateWorkspaceFSMBuilder_CloningRepositoryStage
}

type CreateWorkspaceFSMBuilder_CloningRepositoryStage interface {
	CloningRepositoryState(func(ctx context.Context, transitions CloningRepositoryTransitions, i WorkspaceContext) (fsm.Transition, error)) CreateWorkspaceFSMBuilder_DoneStage
}

type CreateWorkspaceFSMBuilder_DoneStage interface {
	DoneState(func(ctx context.Context, i WorkspaceContext) (fsm.Transition, error)) CreateWorkspaceFSMBuilder_ErrorStage
}

type CreateWorkspaceFSMBuilder_ErrorStage interface {
	ErrorState(func(ctx context.Context, i WorkspaceContext) (fsm.Transition, error)) CreateWorkspaceFSMBuilder_Build
}

type CreateWorkspaceFSMBuilder_Build interface {
	Build() CreateWorkspaceFSM
}

type InitialTransitions interface {
	ToCreatingRecord() (fsm.Transition, error)
	ToError() (fsm.Transition, error)
}

type CreatingRecordTransitions interface {
	ToCloningRepository() (fsm.Transition, error)
	ToError() (fsm.Transition, error)
}

type CloningRepositoryTransitions interface {
	ToDone() (fsm.Transition, error)
	ToError() (fsm.Transition, error)
}

func NewCreateWorkspaceFSM() CreateWorkspaceFSMBuilder_InitialStage {
	return nil
}
