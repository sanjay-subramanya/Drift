package rules

import "drift/internal/core/model"

type BranchRule struct{}

func (r BranchRule) Name() string {
	return "branch-drift"
}

func (r BranchRule) Evaluate(ctx model.Context, drifts []model.Drift) []model.Finding {
	var findings []model.Finding

	for _, d := range drifts {
		if d.Type != model.DriftBranch {
			continue
		}

		findings = append(findings, model.Finding{
			Drift:    d,
			Severity: model.SeverityHigh,
			Message:  d.Summary,
		})
	}

	return findings
}
