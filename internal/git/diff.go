package git

import "strings"

// func CommitsBehind(base string) (string, error) {
// 	return RunGit("rev-list", "--count", "HEAD.."+base)
// }

func UpstreamFiles(local, remote string) ([]string, error) {
	base, err := MergeBase(local, remote)
	if err != nil {
		return nil, err
	}

	out, err := RunGit("diff", "--name-only", base, remote)
	if err != nil {
		return nil, err
	}

	if out == "" {
		return nil, nil
	}

	return strings.Split(out, "\n"), nil
}