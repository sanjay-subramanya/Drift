package analyzers

import (
	"drift/internal/config"
	"drift/internal/core/model"
	"drift/internal/git"
	"drift/internal/workspace"
)

func AnalyzeEnvAndDocker(base string) ([]model.Drift, error) {
	_, _ = git.RunGit("fetch", "origin")

	mergeBase, err := git.MergeBase("HEAD", base)
	if err != nil {
		return nil, err
	}

	upstreamFiles, err := git.UpstreamFiles(mergeBase, base)
	if err != nil {
		return nil, err
	}

	ignores := config.LoadIgnoreFile()

	var affected []string
	for _, f := range upstreamFiles {
		if config.IsIgnored(f, ignores) {
			continue
		}
		if workspace.IsEnvFile(f) || workspace.IsDeploymentFile(f) {
			affected = append(affected, f)
		}
	}

	if len(affected) == 0 {
		return nil, nil
	}

	return []model.Drift{
		{
			Type: model.DriftEnv,
			Summary: "[HIGH] Env/Docker files changed upstream: " +
				stringJoin(affected),
		},
	}, nil
}

func join(items []string) string {
	out := ""
	for _, i := range items {
		out += i + ", "
	}
	return out[:len(out)-2]
}
