package git

func MergeBase(a, b string) (string, error) {
	return RunGit("merge-base", a, b)
}