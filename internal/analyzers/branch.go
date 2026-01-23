package analyzers

import (
	"github.com/sanjay-subramanya/drift/internal/core/model"
	"github.com/sanjay-subramanya/drift/internal/git"
	"github.com/sanjay-subramanya/drift/internal/workspace"
	"github.com/sanjay-subramanya/drift/internal/config"
	"fmt"
	"slices"
	"strings"
	"strconv"
)

func AnalyzeBranch(ctx model.Context) ([]model.Drift, error) {
	baseBranch := ctx.Base
	upstreamURL := ctx.UpstreamURL
	if after, ok := strings.CutPrefix(baseBranch, "origin/"); ok  {
		baseBranch = after
	}
	compareRef := baseBranch
	isFork := false
	
	if upstreamURL == "" {
		compareRef = "origin/" + baseBranch
		if _, err := git.RunGit("rev-parse", "--verify", "--quiet", compareRef); err != nil {
			return nil, fmt.Errorf("Base branch \"%s\" does not exist", compareRef)
		}
		if _, err := git.RunGit("fetch", "origin"); err != nil {
			return nil, fmt.Errorf("Failed to fetch origin: %w", err)
		}
	} else {
		if !strings.HasSuffix(upstreamURL, ".git") {
			upstreamURL += ".git"
		}
		if _, err := git.RunGit("fetch", upstreamURL, compareRef); err != nil {
			return nil, fmt.Errorf("Failed to fetch upstream %s: %w", upstreamURL, err)
		}
		isFork = true
		compareRef = "FETCH_HEAD"
	}

	mergeBase, err := git.MergeBase("HEAD", compareRef)
	if err != nil {
		return nil, err
	}

	behind, err := git.CommitsBehind("HEAD", compareRef)
	if err != nil {
		return nil, err
	}
	if behind == 0 {
		return nil, nil
	}

	upstreamFiles, err := git.UpstreamFiles(mergeBase, compareRef, isFork)
	if err != nil {
		return nil, err
	}

	localDirty, err := git.LocalChanges(mergeBase)
	if err != nil {
		return nil, err
	}

	ignores := config.LoadIgnoreFile()

	// Classification buckets (each file goes into only 1)
	var critical []string
	var high []string
	var low []string

	depHits := DependencyHits(localDirty, upstreamFiles)

	for _, f := range upstreamFiles {
		if config.IsIgnored(f, ignores) {
			continue
		}

		switch {
		// file exists locally (dirty)
		case slices.Contains(localDirty, f) && slices.Contains(upstreamFiles, f):
			critical = append(critical, f)

		// env / docker / deployment files
		case workspace.IsEnvFile(f) || workspace.IsDeploymentFile(f):
			high = append(high, f)

		// dependency you import
		case slices.Contains(depHits, f):
			high = append(high, f)

		// everything else
		default:
			if !isFork {
				low = append(low, f)
			}
		}
	}

	var lines []string
	lines = append(lines, "branch behind by "+ strconv.Itoa(behind) +" commits;")

	if len(critical) > 0 {
		lines = append(lines,
			"[CRITICAL] files YOU are editing changed upstream: "+ stringJoin(critical))
	}
	if len(high) > 0 {
		lines = append(lines,
			"[HIGH] deployment / dependency files changed upstream: "+ stringJoin(high))
	}
	if len(low) > 0 {
		lines = append(lines,
			"[LOW] other files changed upstream: "+ stringJoin(low))
	}

	return []model.Drift{
		{
			Type:    model.DriftBranch,
			Summary: strings.Join(lines, "\n"),
		},
	}, nil
}
