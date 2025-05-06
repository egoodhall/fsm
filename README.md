# FSM - Finite State Machine Generator for Go

A type-safe finite state machine generator for Go that creates code from YAML definitions.

## Installation

```bash
go install github.com/egoodhall/fsm/cmd/fsmgen@latest
```

## Features

- Type-safe state machine definitions
- Parallel processing with configurable worker counts
- In-memory or persistent state storage
- Automatic code generation from YAML definitions

## Usage

1. Define your state machine in YAML:

```yaml
# create_workspace.yaml
name: CreateWorkspace
states:
- 
  name: CreateRecord
  entrypoint: true
  inputs:
  - WorkspaceContext
  transitions:
  - CloneRepo
  - Error
- 
  name: CloneRepo
  workers: 5
  inputs:
  - WorkspaceContext
  - WorkspaceID
  transitions:
  - Done
  - Error
- 
  name: Done
  terminal: true
  inputs:
  - WorkspaceContext
  - WorkspaceID
- 
  name: Error
  terminal: true
  inputs:
  - WorkspaceContext
```

2. Define your custom types:

```go
// example.go
package example

//go:generate go run github.com/egoodhall/fsm/cmd/fsmgen -out . -pkg example create_workspace.yaml

type WorkspaceID int64

type WorkspaceContext struct {
    RepositoryURL string
    BranchName    string
}
```

3. Generate FSM code:

```bash
fsmgen -out ./generated -pkg example create_workspace.yaml
# or use go:generate
go generate ./...
```

4. Use the generated FSM:

```go
fsm, err := example.NewCreateWorkspaceFSMBuilder().
    CreateRecordState(func(ctx context.Context, transitions example.CreateRecordTransitions, c example.WorkspaceContext) error {
        return transitions.ToCloneRepo(ctx, c, example.WorkspaceID(1))
    }).
    CloneRepoState(func(ctx context.Context, transitions example.CloneRepoTransitions, c example.WorkspaceContext, i example.WorkspaceID) error {
        return transitions.ToDone(ctx, c, i)
    }).
    DoneState(func(ctx context.Context, c example.WorkspaceContext, i example.WorkspaceID) error {
        // Terminal state handler
        return nil
    }).
    ErrorState(func(ctx context.Context, c example.WorkspaceContext) error {
        // Error state handler
        return nil
    }).
    BuildAndStart(context.Background())

// Submit a task to the FSM
id, err := fsm.Submit(context.Background(), example.WorkspaceContext{
    RepositoryURL: "https://github.com/example/repo",
    BranchName: "main",
})
```

