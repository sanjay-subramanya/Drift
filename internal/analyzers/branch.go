package analyzers

import (
	"drift/internal/core/model"
	"drift/internal/git"
	"fmt"
	"strings"
	"strconv"
)

func AnalyzeBranch(base string) ([]model.Drift, error) {
	if _, err := git.RunGit("rev-parse", "--verify", "--quiet", base); err != nil {
		return nil, fmt.Errorf("Base branch \"%s\" does not exist", base)
	}
	_, _ = git.RunGit("fetch", "origin")

	mergeBase, err := git.MergeBase("HEAD", base)
	if err != nil {
		return nil, err
	}

	behind, err := git.CommitsBehind("HEAD", base)
	if err != nil {
		return nil, err
	}

	if behind == 0 {
		return nil, nil
	}

	upstreamFiles, err := git.UpstreamFiles(mergeBase, base)
	if err != nil {
		return nil, err
	}

	localDirty, err := git.DirtyFiles()
	if err != nil {
		return nil, err
	}

	// --- Sprint 3 logic ---
	localFilesOut, err := git.RunGit("ls-files")
	if err != nil {
		return nil, err
	}
	localTracked := strings.Split(localFilesOut, "\n")

	// CRITICAL: upstream files that exist locally (dirty OR clean)
	directHits := intersect(localTracked, upstreamFiles)
	depHits := dependencyHits(localDirty, upstreamFiles)
	otherHits := subtract(upstreamFiles, directHits, depHits)

	var lines []string

	lines = append(lines, "Branch behind by " + strconv.Itoa(behind) + " commits;")

	if len(directHits) > 0 {
		lines = append(lines, "[CRITICAL] Files YOU are editing changed upstream: " + stringJoin(directHits))
	}
	if len(depHits) > 0 {
		lines = append(lines, "[HIGH] Dependencies you import changed upstream: " + stringJoin(depHits))
	}
	if len(otherHits) > 0 {
		lines = append(lines, "[LOW] Other upstream files changed: " + stringJoin(otherHits))
	}
	summary := strings.Join(lines, "\n")

	drifts := []model.Drift{
		{
			Type:    model.DriftBranch,
			Summary: summary,
		},
	}

	return drifts, nil
}

/* ---------------- helpers ---------------- */

func intersect(a, b []string) []string {
	set := make(map[string]bool)
	for _, x := range a {
		set[x] = true
	}
	out := []string{}
	for _, y := range b {
		if set[y] {
			out = append(out, y)
		}
	}
	return out
}

// VERY conservative dependency detection (Sprint-3 safe)
func dependencyHits(localFiles, upstreamFiles []string) []string {
	out := []string{}

	for _, local := range localFiles {
		imports, err := extractImports(local)
		if err != nil {
			continue
		}
		for _, imp := range imports {
			for _, up := range upstreamFiles {
				if imp == up {
					out = append(out, up)
				}
			}
		}
	}

	return unique(out)
}

// crude but safe: regex-free, line-based
func extractImports(file string) ([]string, error) {
	data, err := git.RunGit("show", "HEAD:"+file)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(data, "\n")
	out := []string{}

	for _, l := range lines {
		l = strings.TrimSpace(l)

		// Go / Python / JS style heuristics
		if strings.HasPrefix(l, "import ") ||
			strings.HasPrefix(l, "from ") ||
			strings.Contains(l, "require(") {

			parts := strings.Fields(l)
			for _, p := range parts {
				if strings.Contains(p, "/") && strings.Contains(p, ".") {
					out = append(out, strings.Trim(p, "\"'"))
				}
			}
		}
	}

	return out, nil
}

func subtract(all []string, a []string, b []string) []string {
	exclude := map[string]bool{}
	for _, x := range a {
		exclude[x] = true
	}
	for _, x := range b {
		exclude[x] = true
	}

	out := []string{}
	for _, x := range all {
		if !exclude[x] {
			out = append(out, x)
		}
	}
	return out
}

func unique(in []string) []string {
	m := map[string]bool{}
	out := []string{}
	for _, x := range in {
		if !m[x] {
			m[x] = true
			out = append(out, x)
		}
	}
	return out
}

func stringJoin(files []string) string {
	if len(files) == 0 {
		return "none"
	}
	return strings.Join(files, ", ")
}

// func AnalyzeBranch(base string) ([]model.Drift, error) {
// 	// 1. Ensure remote refs are fresh (read-only)
// 	if err := git.Fetch(); err != nil {
// 		return nil, err
// 	}

// 	// 2. Compute merge-base between HEAD and base (e.g. origin/main)
// 	mergeBase, err := git.MergeBase("HEAD", base)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 3. Count how many commits HEAD is behind base
// 	behind, err := git.CommitsBehind("HEAD", base)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if behind == 0 {
// 		return nil, nil
// 	}

// 	// 4. List files changed upstream since branch point
// 	files, err := git.UpstreamFiles(mergeBase, base)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return []model.Drift{
// 		{
// 			Type: model.DriftBranch,
// 			Summary: "Branch is behind base by " +
// 				strconv.Itoa(behind) +
// 				" commits; upstream files changed: " +
// 				stringJoin(files),
// 		},
// 	}, nil
// }

// func stringJoin(files []string) string {
// 	if len(files) == 0 {
// 		return "none"
// 	}
// 	out := ""
// 	for _, f := range files {
// 		out += f + ", "
// 	}
// 	return out[:len(out)-2]
// }
