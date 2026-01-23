package model

type Context struct {
	UpstreamURL string
	Base     string
	JSON     bool
	JSONPath string
}

func NewContext() Context {
	return Context{
		UpstreamURL: "",
		Base:     "origin/main",
		JSON:     false,
		JSONPath: "drift.json",
	}
}
