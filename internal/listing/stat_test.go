package listing

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStat_NonExistent(t *testing.T) {
	_, err := Stat("/tmp/mstat_no_such_file_ever", false)
	if err == nil {
		t.Error("Stat() on non-existent file should return error")
	}
}

func TestStat_Directory(t *testing.T) {
	dir := t.TempDir()
	e, err := Stat(dir, false)
	if err != nil {
		t.Fatalf("Stat(dir) error: %v", err)
	}
	if e.Kind != KindDirectory {
		t.Errorf("Kind = %v, want KindDirectory", e.Kind)
	}
	if e.Name != filepath.Base(dir) {
		t.Errorf("Name = %q, want %q", e.Name, filepath.Base(dir))
	}
}

func TestStat_Symlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target.txt")
	if err := os.WriteFile(target, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "link.txt")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	e, err := Stat(link, false)
	if err != nil {
		t.Fatalf("Stat(symlink) error: %v", err)
	}
	if e.Kind != KindSymlink {
		t.Errorf("Kind = %v, want KindSymlink", e.Kind)
	}
}

func TestStat_SymlinkDereference(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target.txt")
	if err := os.WriteFile(target, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "link.txt")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	e, err := Stat(link, true)
	if err != nil {
		t.Fatalf("Stat(symlink, deref) error: %v", err)
	}
	if e.Kind != KindFile {
		t.Errorf("Kind = %v, want KindFile", e.Kind)
	}
	if e.Size != 5 {
		t.Errorf("Size = %d, want 5", e.Size)
	}
}

func TestStat_SpacesInName(t *testing.T) {
	dir := t.TempDir()
	name := "file with spaces.txt"
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}
	e, err := Stat(path, false)
	if err != nil {
		t.Fatalf("Stat(spaces) error: %v", err)
	}
	if e.Name != name {
		t.Errorf("Name = %q, want %q", e.Name, name)
	}
	if e.Kind != KindFile {
		t.Errorf("Kind = %v, want KindFile", e.Kind)
	}
}

func TestStat_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")
	if err := os.WriteFile(path, nil, 0644); err != nil {
		t.Fatal(err)
	}
	e, err := Stat(path, false)
	if err != nil {
		t.Fatalf("Stat(empty) error: %v", err)
	}
	if e.Size != 0 {
		t.Errorf("Size = %d, want 0", e.Size)
	}
}

func TestStat_LongName(t *testing.T) {
	dir := t.TempDir()
	name := "this_is_a_very_long_filename_that_exceeds_normal_expectations_for_a_test_case.txt"
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	e, err := Stat(path, false)
	if err != nil {
		t.Fatalf("Stat(long) error: %v", err)
	}
	if e.Name != name {
		t.Errorf("Name = %q, want %q", e.Name, name)
	}
}
