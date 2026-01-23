package engine

import (
	"github.com/sanjay-subramanya/drift/internal/analyzers"
	"github.com/sanjay-subramanya/drift/internal/core/model"
	"github.com/sanjay-subramanya/drift/internal/core/rules"
)

type Engine struct {
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Run(ctx model.Context) ([]model.Finding, error) {
	var drifts []model.Drift

	branchDrift, err := analyzers.AnalyzeBranch(ctx)
	if err != nil {
		return nil, err
	}
	drifts = append(drifts, branchDrift...)

	var findings []model.Finding
	findings = append(findings, rules.BranchRule{}.Evaluate(ctx, drifts)...)
	findings = append(findings, rules.EnvRule{}.Evaluate(ctx, drifts)...)

	return findings, nil
}
