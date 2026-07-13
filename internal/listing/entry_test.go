package listing

import (
	"io/fs"
	"testing"
)

func TestKindString(t *testing.T) {
	tests := []struct {
		kind Kind
		want string
	}{
		{KindFile, "file"},
		{KindDirectory, "directory"},
		{KindSymlink, "symbolic link"},
		{KindPipe, "pipe"},
		{KindSocket, "socket"},
		{KindDevice, "block device"},
		{KindCharDevice, "character device"},
		{Kind(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.kind.String()
			if got != tt.want {
				t.Errorf("Kind.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestKindFromMode(t *testing.T) {
	tests := []struct {
		name string
		mode fs.FileMode
		want Kind
	}{
		{"regular file", 0, KindFile},
		{"directory", fs.ModeDir, KindDirectory},
		{"symlink", fs.ModeSymlink, KindSymlink},
		{"pipe", fs.ModeNamedPipe, KindPipe},
		{"socket", fs.ModeSocket, KindSocket},
		{"device", fs.ModeDevice, KindDevice},
		{"char device", fs.ModeCharDevice, KindCharDevice},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := kindFromMode(tt.mode)
			if got != tt.want {
				t.Errorf("kindFromMode(%v) = %v, want %v", tt.mode, got, tt.want)
			}
		})
	}
}

func TestPermissions(t *testing.T) {
	tests := []struct {
		name string
		mode fs.FileMode
		want string
	}{
		{"regular file rw-r--r--", 0644, "-rw-r--r--"},
		{"regular file rwxr-xr-x", 0755, "-rwxr-xr-x"},
		{"directory", fs.ModeDir | 0755, "drwxr-xr-x"},
		{"all permissions", 0777, "-rwxrwxrwx"},
		{"no permissions", 0000, "----------"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Entry{Mode: tt.mode}
			got := e.Permissions()
			if got != tt.want {
				t.Errorf("Permissions() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFileTypeChar(t *testing.T) {
	tests := []struct {
		mode fs.FileMode
		want byte
	}{
		{0, '-'},
		{fs.ModeDir, 'd'},
		{fs.ModeSymlink, 'l'},
		{fs.ModeNamedPipe, 'p'},
		{fs.ModeSocket, 's'},
		{fs.ModeDevice, 'b'},
		{fs.ModeCharDevice, 'c'},
	}

	for _, tt := range tests {
		t.Run(string(tt.want), func(t *testing.T) {
			got := fileTypeChar(tt.mode)
			if got != tt.want {
				t.Errorf("fileTypeChar(%v) = %c, want %c", tt.mode, got, tt.want)
			}
		})
	}
}
