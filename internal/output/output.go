package output

import (
	"io"

	"github.com/bhavya-dang/mstat/internal/listing"
)

// Render writes the appropriate output for the given entries.
func Render(w io.Writer, entries []listing.Entry) {
	RenderTable(w, entries)
}
