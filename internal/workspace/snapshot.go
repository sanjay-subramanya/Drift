package workspace

import "drift/internal/git"

type Snapshot struct {
	Branch string
}

func Capture() (Snapshot, error) {
	branch, err := git.CurrentBranch()
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{Branch: branch}, nil
}