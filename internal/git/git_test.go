package git

import (
	"path/filepath"
	"testing"
)

func TestParsePorcelain(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect map[string]Status
	}{
		{
			name:  "modified file",
			input: "M  file.txt",
			expect: map[string]Status{
				filepath.Join("/repo", "file.txt"): Modified,
			},
		},
		{
			name:  "untracked file",
			input: "?? newfile.txt",
			expect: map[string]Status{
				filepath.Join("/repo", "newfile.txt"): Untracked,
			},
		},
		{
			name:  "deleted file",
			input: " D deleted.txt",
			expect: map[string]Status{
				filepath.Join("/repo", "deleted.txt"): Deleted,
			},
		},
		{
			name:  "added file",
			input: "A  staged.txt",
			expect: map[string]Status{
				filepath.Join("/repo", "staged.txt"): Added,
			},
		},
		{
			name: "multiple files",
			input: "M  a.txt\n" +
				"?? b.txt\n" +
				" D c.txt",
			expect: map[string]Status{
				filepath.Join("/repo", "a.txt"): Modified,
				filepath.Join("/repo", "b.txt"): Untracked,
				filepath.Join("/repo", "c.txt"): Deleted,
			},
		},
		{
			name:  "renamed file",
			input: "R  old.txt -> new.txt",
			expect: map[string]Status{
				filepath.Join("/repo", "new.txt"): Renamed,
			},
		},
		{
			name:   "empty input",
			input:  "",
			expect: map[string]Status{},
		},
		{
			name:  "both modified index and worktree",
			input: "MM both.txt",
			expect: map[string]Status{
				filepath.Join("/repo", "both.txt"): Modified,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsePorcelain([]byte(tt.input), "/repo")
			if len(got) != len(tt.expect) {
				t.Fatalf("got %d entries, want %d", len(got), len(tt.expect))
			}
			for path, want := range tt.expect {
				if got[path] != want {
					t.Errorf("got[%q] = %q, want %q", path, got[path], want)
				}
			}
		})
	}
}

func TestStatusString(t *testing.T) {
	tests := []struct {
		status Status
		want   string
	}{
		{Clean, "clean"},
		{Modified, "modified"},
		{Added, "added"},
		{Deleted, "deleted"},
		{Renamed, "renamed"},
		{Untracked, "untracked"},
		{Ignored, "ignored"},
		{Conflict, "conflict"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.status.String()
			if got != tt.want {
				t.Errorf("Status(%q).String() = %q, want %q", string(tt.status), got, tt.want)
			}
		})
	}
}

func TestClassifyStatus(t *testing.T) {
	tests := []struct {
		xy   string
		want Status
	}{
		{"??", Untracked},
		{"!!", Ignored},
		{"M ", Modified},
		{" M", Modified},
		{"MM", Modified},
		{"A ", Added},
		{" D", Deleted},
		{"D ", Deleted},
		{"R ", Renamed},
		{"UU", Conflict},
		{"AU", Conflict},
		{"UD", Conflict},
		{"  ", Clean},
	}

	for _, tt := range tests {
		t.Run(string(tt.want), func(t *testing.T) {
			got := classifyStatus(tt.xy)
			if got != tt.want {
				t.Errorf("classifyStatus(%q) = %q, want %q", tt.xy, got, tt.want)
			}
		})
	}
}
