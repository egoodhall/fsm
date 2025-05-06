package example

//go:generate go run ../cmd/fsmgen -out . -pkg example create_workspace.yaml

type WorkspaceID int64

type WorkspaceContext struct {
	RepositoryURL string
	BranchName    string
}
