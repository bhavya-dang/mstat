package output

import (
	"io"

	"github.com/bhavya-dang/mstat/internal/git"
	"github.com/bhavya-dang/mstat/internal/listing"
)

// rendering table with different views and config options like no icons, simple icons, extended view, etc.
type Options struct {
	Icons        bool
	SimpleIcons  bool
	BriefView    bool
	ExtendedView bool
	NoColor      bool
	NoGit        bool
	Porcelain    bool
	FullPath     bool
	GitMap       map[string]git.Status
}

// writes the output for the given entries.
func Render(w io.Writer, entries []listing.Entry, opts Options) {
	if opts.Porcelain {
		opts.NoColor = true
	}
	RenderTable(w, entries, opts)
}
