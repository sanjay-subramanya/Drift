package rules

import "drift/internal/core/model"

type Rule interface {
	Evaluate (ctx model.Context) ([]model.Finding, error)
	Name() string
}
