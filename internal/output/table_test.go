package output

import (
	"bytes"
	"io/fs"
	"strings"
	"testing"
	"time"

	"github.com/bhavya-dang/mstat/internal/git"
	"github.com/bhavya-dang/mstat/internal/listing"
)

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{0, "0 B"},
		{1, "1 B"},
		{512, "512 B"},
		{1023, "1023 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
		{1536000, "1.5 MB"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatSize(tt.bytes)
			if got != tt.want {
				t.Errorf("formatSize(%d) = %q, want %q", tt.bytes, got, tt.want)
			}
		})
	}
}

func TestFormatRelativeTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name string
		time time.Time
		want string
	}{
		{"just now", now.Add(-30 * time.Second), "just now"},
		{"5 minutes ago", now.Add(-5 * time.Minute), "5m ago"},
		{"3 hours ago", now.Add(-3 * time.Hour), "3h ago"},
		{"2 days ago", now.Add(-2 * 24 * time.Hour), "2d ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatRelativeTime(tt.time)
			if got != tt.want {
				t.Errorf("formatRelativeTime() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatWithRelative(t *testing.T) {
	now := time.Now()
	got := formatWithRelative(now)
	if got == "" {
		t.Error("formatWithRelative() returned empty string")
	}
}

func TestBuildColumns(t *testing.T) {
	tests := []struct {
		name string
		opts Options
		want int // number of columns
	}{
		{"brief", Options{BriefView: true}, 3},
		{"default", Options{}, 5},
		{"extended", Options{ExtendedView: true}, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cols := buildColumns(tt.opts)
			if len(cols) != tt.want {
				t.Errorf("buildColumns() returned %d columns, want %d", len(cols), tt.want)
			}
		})
	}
}

func TestAlignCell(t *testing.T) {
	tests := []struct {
		name      string
		cell      string
		cellWidth int
		width     int
		align     cellAlign
		want      string
	}{
		{"left align", "hi", 2, 4, alignLeft, "hi  "},
		{"right align", "hi", 2, 4, alignRight, "  hi"},
		{"no padding needed", "hello", 5, 3, alignLeft, "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := alignCell(tt.cell, tt.cellWidth, tt.width, tt.align)
			if got != tt.want {
				t.Errorf("alignCell() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRenderTable_Empty(t *testing.T) {
	var buf bytes.Buffer
	RenderTable(&buf, nil, Options{})
	if buf.Len() != 0 {
		t.Errorf("RenderTable(empty) wrote %d bytes, want 0", buf.Len())
	}
}

func TestRenderTable_SingleFile(t *testing.T) {
	entries := []listing.Entry{{
		Name:     "test.txt",
		Kind:     listing.KindFile,
		Mode:     0644,
		Size:     100,
		Links:    1,
		Modified: time.Now(),
	}}
	var buf bytes.Buffer
	RenderTable(&buf, entries, Options{Icons: false})
	out := buf.String()
	if !strings.Contains(out, "test.txt") {
		t.Error("output missing filename")
	}
	if !strings.Contains(out, "file") {
		t.Error("output missing type")
	}
	if !strings.Contains(out, "100 B") {
		t.Error("output missing size")
	}
}

func TestRenderTable_MixFilesAndDirs(t *testing.T) {
	now := time.Now()
	entries := []listing.Entry{
		{Name: "dir1", Kind: listing.KindDirectory, Mode: fs.ModeDir | 0755, Size: 0, Links: 2, Modified: now},
		{Name: "file.txt", Kind: listing.KindFile, Mode: 0644, Size: 50, Links: 1, Modified: now},
		{Name: "dir2", Kind: listing.KindDirectory, Mode: fs.ModeDir | 0755, Size: 0, Links: 2, Modified: now},
	}
	var buf bytes.Buffer
	RenderTable(&buf, entries, Options{})
	out := buf.String()
	if !strings.Contains(out, "dir1") || !strings.Contains(out, "dir2") || !strings.Contains(out, "file.txt") {
		t.Error("output missing entries")
	}
	if strings.Count(out, "│") < 6 {
		t.Error("table seems malformed")
	}
}

func TestRenderTable_LongName(t *testing.T) {
	long := strings.Repeat("a", 200)
	entries := []listing.Entry{{
		Name:     long,
		Kind:     listing.KindFile,
		Mode:     0644,
		Size:     1,
		Links:    1,
		Modified: time.Now(),
	}}
	var buf bytes.Buffer
	RenderTable(&buf, entries, Options{})
	if !strings.Contains(buf.String(), long) {
		t.Error("output missing long filename")
	}
}

func TestRenderTable_SymlinkEntry(t *testing.T) {
	entries := []listing.Entry{{
		Name:     "link",
		Kind:     listing.KindSymlink,
		Mode:     fs.ModeSymlink | 0777,
		Size:     10,
		Links:    1,
		Modified: time.Now(),
	}}
	var buf bytes.Buffer
	RenderTable(&buf, entries, Options{})
	if !strings.Contains(buf.String(), "symbolic link") {
		t.Error("output missing symlink type")
	}
}

func TestRenderTable_BriefView(t *testing.T) {
	entries := []listing.Entry{{
		Name:     "test.txt",
		Kind:     listing.KindFile,
		Mode:     0644,
		Size:     2048,
		Links:    1,
		Modified: time.Now(),
	}}
	var buf bytes.Buffer
	RenderTable(&buf, entries, Options{BriefView: true})
	out := buf.String()
	if !strings.Contains(out, "2.0 KB") {
		t.Error("output missing size in brief view")
	}
	if strings.Contains(out, "permissions") {
		t.Error("brief view should not contain permissions column")
	}
}

func TestRenderTable_ExtendedView(t *testing.T) {
	entries := []listing.Entry{{
		Name:     "test.txt",
		Kind:     listing.KindFile,
		Mode:     0644,
		Size:     10,
		Links:    1,
		Modified: time.Now(),
	}}
	var buf bytes.Buffer
	RenderTable(&buf, entries, Options{ExtendedView: true})
	out := buf.String()
	if !strings.Contains(out, "permissions octal") {
		t.Error("extended view missing permissions octal column")
	}
	if !strings.Contains(out, "links") {
		t.Error("extended view missing links column")
	}
}

func TestRenderTable_NoIcons(t *testing.T) {
	entries := []listing.Entry{{
		Name:     "test.txt",
		Kind:     listing.KindFile,
		Mode:     0644,
		Size:     1,
		Links:    1,
		Modified: time.Now(),
	}}
	var buf bytes.Buffer
	RenderTable(&buf, entries, Options{Icons: false})
	out := buf.String()
	if !strings.Contains(out, "test.txt") {
		t.Error("output missing filename with icons disabled")
	}
}

func TestRenderTable_DirColor(t *testing.T) {
	entries := []listing.Entry{
		{Name: "mydir", Kind: listing.KindDirectory, Mode: fs.ModeDir | 0755, Size: 0, Links: 2, Modified: time.Now()},
		{Name: "file.txt", Kind: listing.KindFile, Mode: 0644, Size: 1, Links: 1, Modified: time.Now()},
	}
	var buf bytes.Buffer
	RenderTable(&buf, entries, Options{Icons: false})
	out := buf.String()
	if !strings.Contains(out, "\x1b[1;34mmydir\x1b[0m") {
		t.Error("directory name should be wrapped in bold blue ANSI codes")
	}
	if strings.Contains(out, "\x1b[1;34mfile.txt") {
		t.Error("file name should not be colored")
	}
}

func TestRenderTable_NoColor(t *testing.T) {
	entries := []listing.Entry{{
		Name: "mydir", Kind: listing.KindDirectory, Mode: fs.ModeDir | 0755, Size: 0, Links: 2, Modified: time.Now(),
	}}
	var buf bytes.Buffer
	RenderTable(&buf, entries, Options{NoColor: true})
	out := buf.String()
	if strings.Contains(out, "\x1b[") {
		t.Error("output should not contain ANSI codes when NoColor is true")
	}
}

func TestRenderTable_DirColorWithIcon(t *testing.T) {
	entries := []listing.Entry{{
		Name: "mydir", Kind: listing.KindDirectory, Mode: fs.ModeDir | 0755, Size: 0, Links: 2, Modified: time.Now(),
	}}
	var buf bytes.Buffer
	RenderTable(&buf, entries, Options{Icons: true, SimpleIcons: true})
	out := buf.String()
	if !strings.Contains(out, "\x1b[1;34m") {
		t.Error("directory with icon should be colored")
	}
}

func TestBuildColumns_WithGit(t *testing.T) {
	gm := map[string]git.Status{"/a": git.Modified}
	tests := []struct {
		name string
		opts Options
		want int
	}{
		{"brief+git", Options{BriefView: true, GitMap: gm}, 4},
		{"default+git", Options{GitMap: gm}, 6},
		{"extended+git", Options{ExtendedView: true, GitMap: gm}, 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cols := buildColumns(tt.opts)
			if len(cols) != tt.want {
				t.Errorf("buildColumns() returned %d columns, want %d", len(cols), tt.want)
			}
			if cols[1].header != "git" {
				t.Errorf("column[1].header = %q, want %q", cols[1].header, "git")
			}
		})
	}
}

func TestRenderTable_GitColumn(t *testing.T) {
	gm := map[string]git.Status{
		"/repo/modified.txt":  git.Modified,
		"/repo/clean.txt":     git.Clean,
		"/repo/untracked.txt": git.Untracked,
	}
	entries := []listing.Entry{
		{Name: "modified.txt", Path: "/repo/modified.txt", Kind: listing.KindFile, Mode: 0644, Size: 1, Links: 1, Modified: time.Now()},
		{Name: "clean.txt", Path: "/repo/clean.txt", Kind: listing.KindFile, Mode: 0644, Size: 1, Links: 1, Modified: time.Now()},
		{Name: "untracked.txt", Path: "/repo/untracked.txt", Kind: listing.KindFile, Mode: 0644, Size: 1, Links: 1, Modified: time.Now()},
	}
	var buf bytes.Buffer
	RenderTable(&buf, entries, Options{GitMap: gm})
	out := buf.String()
	if !strings.Contains(out, "git") {
		t.Error("output missing git header")
	}
	if !strings.Contains(out, "modified") {
		t.Error("modified file should show 'modified' label")
	}
	if !strings.Contains(out, "clean") {
		t.Error("clean file should show 'clean' label")
	}
	if !strings.Contains(out, "untracked") {
		t.Error("untracked file should show 'untracked' label")
	}
}

func TestRenderTable_FullPath(t *testing.T) {
	entries := []listing.Entry{{
		Name: "test.txt", Path: "/Users/test/project/test.txt",
		Kind: listing.KindFile, Mode: 0644, Size: 1, Links: 1, Modified: time.Now(),
	}}
	var buf bytes.Buffer
	RenderTable(&buf, entries, Options{FullPath: true, Icons: false})
	if !strings.Contains(buf.String(), "/Users/test/project/test.txt") {
		t.Error("output should contain full path")
	}
}

func TestRenderTable_DefaultRelativePath(t *testing.T) {
	entries := []listing.Entry{{
		Name: "test.txt", Path: "/Users/test/project/test.txt",
		Kind: listing.KindFile, Mode: 0644, Size: 1, Links: 1, Modified: time.Now(),
	}}
	var buf bytes.Buffer
	RenderTable(&buf, entries, Options{Icons: false})
	out := buf.String()
	if strings.Contains(out, "│ /Users/test/project/test.txt") {
		t.Error("default should show relative path, not full path")
	}
	if !strings.Contains(out, "test.txt") {
		t.Error("output should contain filename")
	}
}

func TestRenderTable_Porcelain(t *testing.T) {
	gm := map[string]git.Status{
		"/repo/a.txt": git.Modified,
	}
	entries := []listing.Entry{
		{Name: "a.txt", Path: "/repo/a.txt", Kind: listing.KindFile, Mode: 0644, Size: 100, Links: 1, Modified: time.Date(2026, 7, 9, 22, 11, 0, 0, time.UTC)},
		{Name: "b.txt", Path: "/repo/b.txt", Kind: listing.KindFile, Mode: 0644, Size: 50, Links: 1, Modified: time.Date(2026, 7, 9, 20, 59, 0, 0, time.UTC)},
	}
	var buf bytes.Buffer
	Render(&buf, entries, Options{Porcelain: true, GitMap: gm})
	out := buf.String()
	// should be a bordered table
	if !strings.Contains(out, "╭") {
		t.Error("porcelain should still render bordered table")
	}
	// should have git header
	if !strings.Contains(out, "git") {
		t.Error("porcelain output missing git header")
	}
	// modified file should show raw symbol "M"
	if !strings.Contains(out, "M") {
		t.Error("modified file should show raw 'M' symbol in porcelain mode")
	}
	// no ANSI codes
	if strings.Contains(out, "\x1b[") {
		t.Error("porcelain mode should not contain ANSI codes")
	}
	// clean file should have empty git cell (just spaces)
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.Contains(line, "b.txt") && strings.Contains(line, "clean") {
			t.Error("porcelain mode should not show 'clean' label")
		}
	}
}
