package rules

import "drift/internal/core/model"

type EnvRule struct{}

func (r EnvRule) Name() string {
	return "env-drift"
}

func (r EnvRule) Evaluate(ctx model.Context, drifts []model.Drift) []model.Finding {
	var findings []model.Finding

	for _, d := range drifts {
		if d.Type != model.DriftEnv {
			continue
		}

		findings = append(findings, model.Finding{
			Drift:    d,
			Severity: model.SeverityLow,
			Message:  d.Summary,
		})
	}

	return findings
}
