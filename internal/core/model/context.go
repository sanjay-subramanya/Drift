package model

type Context struct {
	RepoPath string
	Base     string
	JSON     bool
	JSONPath string
}

func NewContext() Context {
	return Context{
		RepoPath: ".",
		Base:     "origin/main",
		JSON:     false,
		JSONPath: "drift.json",
	}
}
