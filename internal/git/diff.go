package git

import "strings"

func UpstreamFiles(local, remote string, isFork bool) ([]string, error) {
	var out string
	var err error

	if !isFork {
		base, err := MergeBase(local, remote)
		if err != nil {
			return nil, err
		}

		out, err = RunGit("diff", "--name-only", base, remote)
		if err != nil {
			return nil, err
		}
	} else {
		out, err = RunGit("diff", "--name-only", "HEAD...FETCH_HEAD")
		if err != nil {
			return nil, err
		}
	}
	if out == "" {
		return nil, nil
	}

	return strings.Split(out, "\n"), nil
}