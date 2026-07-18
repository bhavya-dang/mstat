package output

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bhavya-dang/mstat/internal/git"
	"github.com/bhavya-dang/mstat/internal/icons"
	"github.com/bhavya-dang/mstat/internal/listing"
	"github.com/mattn/go-runewidth"
)

type cellAlign int

const (
	alignLeft cellAlign = iota
	alignRight
)

// column headings
type column struct {
	header      string
	align       cellAlign
	width       int
	headerWidth int
	render      func(listing.Entry, Options) string
}

// returns the file name based on path display flags.
// Default is relative path from cwd; --full-path shows absolute path.
func displayName(e listing.Entry, opts Options) string {
	if opts.FullPath {
		return e.Path
	}
	wd, err := os.Getwd()
	if err != nil {
		return e.Name
	}
	rel, err := filepath.Rel(wd, e.Path)
	if err != nil {
		return e.Name
	}
	return rel
}

var iconNameCol = column{header: "name", align: alignLeft, render: func(e listing.Entry, opts Options) string {
	icon := ""
	if opts.Icons {
		icon = icons.Icon(e, opts.SimpleIcons)
	}
	label := icon + displayName(e, opts)
	if e.Kind == listing.KindDirectory && !opts.NoColor {
		label = colorBlue(label)
	}
	return label
}}

var gitCol = column{header: "git", align: alignLeft, render: func(e listing.Entry, opts Options) string {
	if opts.GitMap == nil {
		return ""
	}
	s, ok := opts.GitMap[e.Path]
	if !ok || s == git.Clean {
		if opts.Porcelain {
			return "-"
		}
		return colorGreen("clean")
	}
	if opts.Porcelain {
		return string(s)
	}
	return colorGitStatus(s)
}}

const (
	ansiBlue   = "\x1b[1;34m"
	ansiGreen  = "\x1b[32m"
	ansiRed    = "\x1b[31m"
	ansiYellow = "\x1b[33m"
	ansiCyan   = "\x1b[36m"
	ansiReset  = "\x1b[0m"
)

func colorBlue(s string) string {
	return ansiBlue + s + ansiReset
}

func colorGreen(s string) string {
	return ansiGreen + s + ansiReset
}

func colorGitStatus(s git.Status) string {
	label := s.String()
	switch s {
	case git.Modified:
		return ansiYellow + label + ansiReset
	case git.Added:
		return ansiGreen + label + ansiReset
	case git.Deleted:
		return ansiRed + label + ansiReset
	case git.Untracked:
		return ansiCyan + label + ansiReset
	default:
		return label
	}
}

// removes ANSI escape sequences so runewidth measures visible width only.
var ansiReplacer = strings.NewReplacer(
	"\x1b[0m", "",
	"\x1b[1;34m", "",
	"\x1b[34m", "",
	"\x1b[32m", "",
	"\x1b[31m", "",
	"\x1b[33m", "",
	"\x1b[36m", "",
)

func stripAnsi(s string) string {
	return ansiReplacer.Replace(s)
}

// show only name, size, and type
var briefColumns = []column{
	iconNameCol,
	{header: "size", align: alignRight, render: func(e listing.Entry, _ Options) string { return formatSize(e.Size) }},
	{header: "type", align: alignLeft, render: func(e listing.Entry, _ Options) string { return e.Kind.String() }},
}

// default view
var defaultColumns = []column{
	iconNameCol,
	{header: "size", align: alignRight, render: func(e listing.Entry, _ Options) string { return formatSize(e.Size) }},
	{header: "type", align: alignLeft, render: func(e listing.Entry, _ Options) string { return e.Kind.String() }},
	{header: "last modified", align: alignLeft, render: func(e listing.Entry, _ Options) string {
		return formatWithRelative(e.Modified)
	}},
	{header: "permissions", align: alignLeft, render: func(e listing.Entry, _ Options) string { return e.Permissions() }},
}

// show all detail columns
var extendedColumns = []column{
	iconNameCol,
	{header: "size", align: alignRight, render: func(e listing.Entry, _ Options) string { return formatSize(e.Size) }},
	{header: "type", align: alignLeft, render: func(e listing.Entry, _ Options) string { return e.Kind.String() }},
	{header: "last modified", align: alignLeft, render: func(e listing.Entry, _ Options) string {
		return formatWithRelative(e.Modified)
	}},
	{header: "permissions", align: alignLeft, render: func(e listing.Entry, _ Options) string { return e.Permissions() }},
	{header: "permissions octal", align: alignLeft, render: func(e listing.Entry, _ Options) string {
		return fmt.Sprintf("%o", e.Mode.Perm())
	}},
	{header: "links", align: alignRight, render: func(e listing.Entry, _ Options) string { return fmt.Sprintf("%d", e.Links) }},
}

// write bordered table output to w.
func RenderTable(w io.Writer, entries []listing.Entry, opts Options) {
	if len(entries) == 0 {
		return
	}

	cols := buildColumns(opts)
	for i := range cols {
		cols[i].headerWidth = runewidth.StringWidth(stripAnsi(cols[i].header))
	}

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

func buildColumns(opts Options) []column {
	var cols []column
	switch {
	case opts.BriefView:
		cols = briefColumns
	case opts.ExtendedView:
		cols = extendedColumns
	default:
		cols = defaultColumns
	}

	if opts.GitMap != nil {
		// insert git column after name
		newCols := make([]column, 0, len(cols)+1)
		newCols = append(newCols, cols[0]) // name
		newCols = append(newCols, gitCol)
		newCols = append(newCols, cols[1:]...)
		cols = newCols
	}

	return cols
}

// draw the borders for the table

func writeBorderTop(b *strings.Builder, cols []column) {
	b.WriteString("╭")
	for i, col := range cols {
		if i > 0 {
			b.WriteString("┬")
		}
		b.WriteString(strings.Repeat("─", col.width+2))
	}
	b.WriteString("╮\n")
}

func writeHeaderRow(b *strings.Builder, cols []column) {
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

func writeBorderMid(b *strings.Builder, cols []column) {
	b.WriteString("├")
	for i, col := range cols {
		if i > 0 {
			b.WriteString("┼")
		}
		b.WriteString(strings.Repeat("─", col.width+2))
	}
	b.WriteString("┤\n")
}

func writeDataRow(b *strings.Builder, cols []column, row []string, cellWidths []int) {
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

func writeBorderBottom(b *strings.Builder, cols []column) {
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
			rw[ci] = runewidth.StringWidth(stripAnsi(cell))
		}
		widths[ri] = rw
	}
	return widths
}

func computeWidths(cols []column, rows [][]string, widths [][]int) {
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
func formatWithRelative(t time.Time) string {
	return t.Format("Jan 2, 2006 15:04") + " (" + formatRelativeTime(t) + ")"
}

func formatRelativeTime(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 7*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	default:
		return t.Format("Jan 2 15:04")
	}
}

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
