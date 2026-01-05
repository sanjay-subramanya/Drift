package engine

import (
	"drift/internal/analyzers"
	"drift/internal/core/model"
	"drift/internal/core/rules"
)

type Engine struct {
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Run(ctx model.Context) ([]model.Finding, error) {
	var drifts []model.Drift

	branchDrift, err := analyzers.AnalyzeBranch(ctx.Base)
	if err != nil {
		return nil, err
	}
	drifts = append(drifts, branchDrift...)

	envDrift, err := analyzers.AnalyzeEnvAndDocker(ctx.Base)
	if err != nil {
		return nil, err
	}
	drifts = append(drifts, envDrift...)

	var findings []model.Finding
	findings = append(findings, rules.BranchRule{}.Evaluate(ctx, drifts)...)
	findings = append(findings, rules.EnvRule{}.Evaluate(ctx, drifts)...)

	return findings, nil
}