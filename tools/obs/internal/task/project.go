package task

import (
	"os/exec"
	"path/filepath"
	"strings"
)

type ProjectInfo struct {
	Name    string
	Dir     string
	Branch  string
	IsImage bool
}

func DetectProjects(tasks []*Task) []ProjectInfo {
	seen := make(map[string]bool)
	var projects []ProjectInfo

	for _, t := range tasks {
		if t.Dir == "" {
			continue
		}
		root := projectRoot(t.Dir)
		if seen[root] {
			continue
		}
		seen[root] = true

		projects = append(projects, ProjectInfo{
			Name:   filepath.Base(root),
			Dir:    root,
			Branch: detectGitBranch(root),
		})
	}

	return projects
}

func projectRoot(dir string) string {
	parts := strings.SplitN(filepath.ToSlash(dir), "/", 3)
	if len(parts) >= 2 {
		return parts[0] + "/" + parts[1]
	}
	return dir
}

func detectGitBranch(dir string) string {
	out, err := exec.Command("git", "-C", dir, "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}
