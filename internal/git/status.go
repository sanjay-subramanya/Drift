package git

import (
	"strconv"
	"strings"
)

func CurrentBranch() (string, error) {
	return RunGit("rev-parse", "--abbrev-ref", "HEAD")
}

func CommitsBehind(local, remote string) (int, error) {
	out, err := RunGit("rev-list", "--count", local+".."+remote)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(out)
}

func LocalChanges(mergeBase string) ([]string, error) {
    out, err := RunGit("diff", "--name-only", mergeBase)
    if err != nil {
        return nil, err
    }
    if out == "" {
        return []string{}, nil
    }
    return strings.Split(out, "\n"), nil
}