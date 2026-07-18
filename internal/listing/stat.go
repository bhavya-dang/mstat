package listing

import (
	"fmt"
	"os"
	"path/filepath"
)

// collects metadata for a single path.
func Stat(path string, dereference bool) (Entry, error) {
	var info os.FileInfo
	var err error
	if dereference {
		info, err = os.Stat(path)
	} else {
		info, err = os.Lstat(path)
	}
	if err != nil {
		return Entry{}, fmt.Errorf("stat %s: %w", path, err)
	}

	abs, _ := filepath.Abs(path)
	mode := info.Mode()
	return Entry{
		Name:     info.Name(),
		Path:     abs,
		Kind:     kindFromMode(mode),
		Mode:     mode,
		Size:     info.Size(),
		Links:    linksOf(info),
		Modified: info.ModTime(),
	}, nil
}
