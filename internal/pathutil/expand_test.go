package pathutil

import (
	"os/user"
	"testing"
)

func TestExpand(t *testing.T) {
	u, err := user.Current()
	if err != nil {
		t.Skip("cannot get current user")
	}

	tests := []struct {
		name string
		path string
		want string
	}{
		{"empty path", "", ""},
		{"no tilde", "file.txt", "file.txt"},
		{"tilde only", "~", u.HomeDir},
		{"tilde with slash", "~/Documents", u.HomeDir + "/Documents"},
		{"tilde with nested", "~/a/b/c", u.HomeDir + "/a/b/c"},
		{"tilde in middle", "foo/~", "foo/~"},
		{"double tilde", "~foo", "~foo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Expand(tt.path)
			if got != tt.want {
				t.Errorf("Expand(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
