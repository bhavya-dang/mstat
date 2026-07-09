package pathutil

import "os/user"

// Expand resolves ~ at the start of a path to the user's home directory.
func Expand(path string) string {
	if len(path) == 0 || path[0] != '~' {
		return path
	}
	u, err := user.Current()
	if err != nil {
		return path
	}
	if path == "~" {
		return u.HomeDir
	}
	if path[1] != '/' {
		return path
	}
	return u.HomeDir + path[1:]
}
