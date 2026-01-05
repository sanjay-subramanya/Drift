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
	// return strconv.Atoi(strings.TrimSpace(string(out)))
	return strconv.Atoi(out)
}

func DirtyFiles() ([]string, error) {
	out, err := RunGit("status", "--porcelain")
	if err != nil {
		return nil, err
	}

	if out == "" {
		return []string{}, nil
	}

	lines := strings.Split(out, "\n")
	files := make([]string, 0)

	for _, line := range lines {
		if len(line) < 4 {
			continue
		}
		// format: XY <file>
		file := strings.TrimSpace(line[3:])
		files = append(files, file)
	}

	return files, nil
}