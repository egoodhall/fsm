package example

//go:generate go run ../cmd/fsmgen -out . -pkg example state_machines.yaml

type WorkspaceID int64

type WorkspaceContext struct {
	RepositoryURL string
	BranchName    string
}
