package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/bhavya-dang/mstat/internal/icons"
	"github.com/bhavya-dang/mstat/internal/listing"
	"github.com/mattn/go-runewidth"
)

type cellAlign int

const (
	alignLeft cellAlign = iota
	alignRight
)

type tableColumn struct {
	header      string
	align       cellAlign
	width       int
	headerWidth int
	render      func(e listing.Entry, opts Options) string
}

var compactColumns = []struct {
	name   string
	header string
	align  cellAlign
	render func(listing.Entry, Options) string
}{
	{"name", "name", alignLeft, func(e listing.Entry, opts Options) string {
		if opts.Icons {
			return icons.Icon(e, opts.SimpleIcons) + e.Name
		}
		return e.Name
	}},
	{"size", "size", alignRight, func(e listing.Entry, _ Options) string { return formatSize(e.Size) }},
	{"type", "type", alignLeft, func(e listing.Entry, _ Options) string { return e.Kind.String() }},
	// {"links", "links", alignRight, func(e listing.Entry, _ Options) string { return fmt.Sprintf("%d", e.Links) }},
	{"modified", "modified", alignLeft, func(e listing.Entry, _ Options) string { return e.Modified.Format("Jan 2 15:04") }},
	{"perms", "perms", alignLeft, func(e listing.Entry, _ Options) string { return e.Permissions() }},
}

// writes bordered table output to w.
func RenderTable(w io.Writer, entries []listing.Entry, opts Options) {
	if len(entries) == 0 {
		return
	}
	renderCompactTable(w, entries, opts)
}

func renderCompactTable(w io.Writer, entries []listing.Entry, opts Options) {
	cols := buildCompactColumns(opts)
	rows := make([][]string, len(entries))
	for i, e := range entries {
		row := make([]string, len(cols))
		for j, col := range cols {
			row[j] = col.render(e, opts)
		}
		rows[i] = row
	}

	widths := measureRows(rows)
	computeWidths(cols, rows, widths)

	var b strings.Builder
	writeBorderTop(&b, cols)
	writeHeaderRow(&b, cols)
	writeBorderMid(&b, cols)
	for ri, row := range rows {
		writeDataRow(&b, cols, row, widths[ri])
	}
	writeBorderBottom(&b, cols)
	fmt.Fprint(w, b.String())
}

func buildCompactColumns(opts Options) []tableColumn {
	cols := make([]tableColumn, 0, len(compactColumns))
	for _, c := range compactColumns {
		hw := runewidth.StringWidth(c.header)
		cols = append(cols, tableColumn{
			header:      c.header,
			align:       c.align,
			headerWidth: hw,
			render:      c.render,
		})
	}
	return cols
}

// --- border drawing ---

func writeBorderTop(b *strings.Builder, cols []tableColumn) {
	b.WriteString("╭")
	for i, col := range cols {
		if i > 0 {
			b.WriteString("┬")
		}
		b.WriteString(strings.Repeat("─", col.width+2))
	}
	b.WriteString("╮\n")
}

func writeHeaderRow(b *strings.Builder, cols []tableColumn) {
	b.WriteRune('│')
	for i, col := range cols {
		if i > 0 {
			b.WriteRune('│')
		}
		pad := col.width - col.headerWidth
		left := pad / 2
		right := pad - left
		b.WriteString(" ")
		b.WriteString(strings.Repeat(" ", left))
		b.WriteString(col.header)
		b.WriteString(strings.Repeat(" ", right))
		b.WriteString(" ")
	}
	b.WriteString("│\n")
}

func writeBorderMid(b *strings.Builder, cols []tableColumn) {
	b.WriteString("├")
	for i, col := range cols {
		if i > 0 {
			b.WriteString("┼")
		}
		b.WriteString(strings.Repeat("─", col.width+2))
	}
	b.WriteString("┤\n")
}

func writeDataRow(b *strings.Builder, cols []tableColumn, row []string, cellWidths []int) {
	b.WriteRune('│')
	for i, cell := range row {
		if i > 0 {
			b.WriteRune('│')
		}
		b.WriteString(" ")
		b.WriteString(alignCell(cell, cellWidths[i], cols[i].width, cols[i].align))
		b.WriteString(" ")
	}
	b.WriteString("│\n")
}

func writeBorderBottom(b *strings.Builder, cols []tableColumn) {
	b.WriteString("╰")
	for i, col := range cols {
		if i > 0 {
			b.WriteString("┴")
		}
		b.WriteString(strings.Repeat("─", col.width+2))
	}
	b.WriteString("╯\n")
}

// --- measurement and alignment ---

func measureRows(rows [][]string) [][]int {
	widths := make([][]int, len(rows))
	for ri, row := range rows {
		rw := make([]int, len(row))
		for ci, cell := range row {
			rw[ci] = runewidth.StringWidth(cell)
		}
		widths[ri] = rw
	}
	return widths
}

func computeWidths(cols []tableColumn, rows [][]string, widths [][]int) {
	for i := range cols {
		cols[i].width = cols[i].headerWidth
	}
	for _, rw := range widths {
		for i, w := range rw {
			if w > cols[i].width {
				cols[i].width = w
			}
		}
	}
}

func alignCell(cell string, cellWidth, width int, align cellAlign) string {
	pad := width - cellWidth
	if pad <= 0 {
		return cell
	}
	if align == alignRight {
		return strings.Repeat(" ", pad) + cell
	}
	return cell + strings.Repeat(" ", pad)
}

// --- formatting helpers ---

func formatSize(b int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case b >= GB:
		return fmt.Sprintf("%.1f GB", float64(b)/float64(GB))
	case b >= MB:
		return fmt.Sprintf("%.1f MB", float64(b)/float64(MB))
	case b >= KB:
		return fmt.Sprintf("%.1f KB", float64(b)/float64(KB))
	default:
		return fmt.Sprintf("%d B", b)
	}
}
