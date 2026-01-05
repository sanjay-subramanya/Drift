package model

type DriftType string

const (
	DriftBranch DriftType = "branch"
	DriftEnv 	DriftType = "env"
	DriftDB 	DriftType = "db"
	DriftAPI 	DriftType = "api"
)

type Drift struct{
	Type 	DriftType
	Summary string
}