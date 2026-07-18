package git

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// represents the git status of a file.
type Status string

const (
	Clean     Status = ""
	Modified  Status = "M"
	Added     Status = "A"
	Deleted   Status = "D"
	Renamed   Status = "R"
	Untracked Status = "?"
	Ignored   Status = "!"
	Conflict  Status = "C"
)

// returns the human-readable label for the status.
func (s Status) String() string {
	switch s {
	case Clean:
		return "clean"
	case Modified:
		return "modified"
	case Added:
		return "added"
	case Deleted:
		return "deleted"
	case Renamed:
		return "renamed"
	case Untracked:
		return "untracked"
	case Ignored:
		return "ignored"
	case Conflict:
		return "conflict"
	default:
		return string(s)
	}
}

// finds the git repo root by running git rev-parse --show-toplevel.
// Returns empty string if not in a git repo.
func RepoRoot(dir string) string {
	cmd := exec.Command("git", "-C", dir, "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// runs git status --porcelain and returns a map of
// absolute path → Status for the given repo root.
func StatusMap(repoRoot string) map[string]Status {
	cmd := exec.Command("git", "-C", repoRoot, "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	return parsePorcelain(out, repoRoot)
}

// parses git status --porcelain output.
// Format: XY filename (XY is 2 chars, space-separated from filename).
func parsePorcelain(data []byte, repoRoot string) map[string]Status {
	m := make(map[string]Status)
	s := strings.TrimRight(string(data), "\n")
	if s == "" {
		return m
	}
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		if len(line) < 3 {
			continue
		}
		xy := line[:2]
		name := strings.TrimSpace(line[3:])

		// Handle renamed files: "R  old -> new"
		if idx := strings.Index(name, " -> "); idx != -1 {
			name = name[idx+4:]
		}

		abs := filepath.Join(repoRoot, name)
		m[abs] = classifyStatus(xy)
	}
	return m
}

// maps the 2-char XY code to a simplified Status.
func classifyStatus(xy string) Status {
	x, y := xy[0], xy[1]

	// Untracked
	if xy == "??" {
		return Untracked
	}
	// Ignored (only shown with -u flag, but handle anyway)
	if xy == "!!" {
		return Ignored
	}
	// Conflict
	if x == 'U' || y == 'U' || (x == 'A' && y == 'A') || (x == 'D' && y == 'D') {
		return Conflict
	}
	// Deleted
	if x == 'D' || y == 'D' {
		return Deleted
	}
	// Renamed
	if x == 'R' || y == 'R' {
		return Renamed
	}
	// Added
	if x == 'A' {
		return Added
	}
	// Modified
	if x == 'M' || y == 'M' {
		return Modified
	}
	return Clean
}
