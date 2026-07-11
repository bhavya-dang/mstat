package icons

import (
	"path/filepath"
	"strings"

	"github.com/bhavya-dang/mstat/internal/listing"
)

// returns the Nerd Font icon string for an entry.
// simple = true => show only kind-based icons like file, folder, link, etc
// simple = false => extension and filename matches are attempted as well
func Icon(e listing.Entry, simple bool) string {
	if !simple && e.Kind == listing.KindFile {
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(e.Name), "."))
		if icon, ok := extIcons[ext]; ok {
			return icon
		}
		if icon, ok := nameIcons[strings.ToLower(e.Name)]; ok {
			return icon
		}
	}
	return basicIcon(e.Kind)
}

func basicIcon(kind listing.Kind) string {
	switch kind {
	case listing.KindDirectory:
		return "\uf07b " // nf-fa-folder
	case listing.KindSymlink:
		return "\uf0c1 " // nf-fa-link
	case listing.KindPipe:
		return "\uf016 " // nf-fa-terminal
	case listing.KindSocket:
		return "\uf1e6 " // nf-fa-plug
	case listing.KindDevice, listing.KindCharDevice:
		return "\uf2db " // nf-fa-microchip
	default:
		return "\uf15b " // nf-fa-file
	}
}

// nameIcons maps exact filenames (lowercased) to icons.
var nameIcons = map[string]string{
	"makefile":       "\ue779 ", // nf-dev-gnu
	"dockerfile":     "\ue7b0 ", // nf-dev-docker
	"docker-compose": "\ue7b0 ",
	"license":        "\uf15c ", // nf-fa-file_text_o
	"license.md":     "\uf15c ",
	"license.txt":    "\uf15c ",
	".gitignore":     "\ue702 ", // nf-dev-git
	".gitmodules":    "\ue702 ",
	".gitattributes": "\ue702 ",
	"go.mod":         "\U000f07d3", // nf-custom-go
	"go.sum":         "\U000f07d3",
	"package.json":   "\ue74e ", // nf-dev-javascript
	"package-lock":   "\ue74e ",
	"tsconfig.json":  "\ue628 ", // nf-seti-typescript
	"cargo.toml":     "\ue7a8 ", // nf-dev-rust
	"cargo.lock":     "\ue7a8 ",
	"pyproject.toml": "\ue73c ", // nf-dev-python
	"setup.py":       "\ue73c ",
	"requirements":   "\ue73c ",
	"pom.xml":        "\ue774 ", // nf-dev-java
	"build.gradle":   "\ue774 ",
	".env":           "\ue779 ",
	".env.local":     "\ue779 ",
	"readme.md":      "\uf15c ",
	"readme.txt":     "\uf15c ",
	"changelog.md":   "\uf15c ",
	"contributing":   "\uf15c ",
}

// extIcons maps file extensions (without dot, lowercased) to icons.
var extIcons = map[string]string{
	// Go
	"go": "\U000f07d3",
	// Rust
	"rs": "\ue7a8 ",
	// Python
	"py":  "\ue73c ",
	"pyw": "\ue73c ",
	// JavaScript / TypeScript
	"js":  "\ue74e ",
	"jsx": "\ue74e ",
	"ts":  "\ue628 ",
	"tsx": "\ue628 ",
	"mjs": "\ue74e ",
	// C / C++
	"c":   "\ue61e ",
	"h":   "\ue61e ",
	"cpp": "\ue61d ",
	"hpp": "\ue61d ",
	// Java
	"java": "\ue774 ",
	// Ruby
	"rb":      "\ue739 ",
	"gemspec": "\ue739 ",
	// PHP
	"php": "\ue77d ",
	// Swift
	"swift": "\ue755 ",
	// Kotlin
	"kt":  "\ue721 ",
	"kts": "\ue721 ",
	// Shell
	"sh":   "\ue795 ",
	"bash": "\ue795 ",
	"zsh":  "\ue795 ",
	"fish": "\ue795 ",
	// Web
	"html":   "\ue736 ",
	"htm":    "\ue736 ",
	"css":    "\ue749 ",
	"scss":   "\ue749 ",
	"sass":   "\ue749 ",
	"less":   "\ue749 ",
	"vue":    "\ue7a3 ",
	"svelte": "\ue697 ",
	// Config / Data
	"json": "\ue60b ",
	"yaml": "\ue615 ",
	"yml":  "\ue615 ",
	"toml": "\ue615 ",
	"xml":  "\ue615 ",
	"ini":  "\ue615 ",
	"conf": "\ue615 ",
	"cfg":  "\ue615 ",
	// Markdown / Docs
	"md":   "\uf15c ",
	"rst":  "\uf15c ",
	"txt":  "\uf15c ",
	"tex":  "\uf15c ",
	"pdf":  "\uf1c1 ",
	"doc":  "\uf1c2 ",
	"docx": "\uf1c2 ",
	"xls":  "\uf1c3 ",
	"xlsx": "\uf1c3 ",
	"csv":  "\uf1c3 ",
	"ppt":  "\uf1c4 ",
	"pptx": "\uf1c4 ",
	// Images
	"png":  "\uf1c5 ",
	"jpg":  "\uf1c5 ",
	"jpeg": "\uf1c5 ",
	"gif":  "\uf1c5 ",
	"svg":  "\uf1c5 ",
	"ico":  "\uf1c5 ",
	"webp": "\uf1c5 ",
	"bmp":  "\uf1c5 ",
	// Video / Audio
	"mp4":  "\uf1c8 ",
	"mkv":  "\uf1c8 ",
	"avi":  "\uf1c8 ",
	"mov":  "\uf1c8 ",
	"mp3":  "\uf1c7 ",
	"wav":  "\uf1c7 ",
	"flac": "\uf1c7 ",
	"ogg":  "\uf1c7 ",
	// Archives
	"tar": "\uf1c6 ",
	"gz":  "\uf1c6 ",
	"zip": "\uf1c6 ",
	"bz2": "\uf1c6 ",
	"xz":  "\uf1c6 ",
	"7z":  "\uf1c6 ",
	"rar": "\uf1c6 ",
	// Database
	"sql":    "\uf1c0 ",
	"db":     "\uf1c0 ",
	"sqlite": "\uf1c0 ",
	// Docker
	"dockerignore": "\ue7b0 ",
	// Git
	"gitignore":     "\ue702 ",
	"gitattributes": "\ue702 ",
	"gitmodules":    "\ue702 ",
	// Build / CI
	"makefile": "\ue779 ",
	"cmake":    "\ue779 ",
	"gradle":   "\ue774 ",
	// Misc
	"lock": "\ue72f ",
	"log":  "\uf15c ",
}
