package listing

import (
	"io/fs"
	"time"
)

// Kind describes the type of a filesystem entry.
type Kind int

const (
	KindFile Kind = iota
	KindDirectory
	KindSymlink
	KindPipe
	KindSocket
	KindDevice
	KindCharDevice
)

func (k Kind) String() string {
	switch k {
	case KindFile:
		return "file"
	case KindDirectory:
		return "directory"
	case KindSymlink:
		return "symbolic link"
	case KindPipe:
		return "pipe"
	case KindSocket:
		return "socket"
	case KindDevice:
		return "block device"
	case KindCharDevice:
		return "character device"
	default:
		return "unknown"
	}
}

// Entry holds metadata for a single filesystem path.
type Entry struct {
	Name string
	Kind Kind
	Mode fs.FileMode
	Size int64
	Links uint64

	Modified time.Time
}

// Permissions returns the 10-char permission string (e.g. "-rw-r--r--").
func (e *Entry) Permissions() string {
	var buf [10]byte
	buf[0] = fileTypeChar(e.Mode)
	formatRwx(buf[1:], e.Mode)
	return string(buf[:])
}

func fileTypeChar(mode fs.FileMode) byte {
	switch {
	case mode&fs.ModeSymlink != 0:
		return 'l'
	case mode&fs.ModeDir != 0:
		return 'd'
	case mode&fs.ModeCharDevice != 0:
		return 'c'
	case mode&fs.ModeDevice != 0:
		return 'b'
	case mode&fs.ModeNamedPipe != 0:
		return 'p'
	case mode&fs.ModeSocket != 0:
		return 's'
	default:
		return '-'
	}
}

func formatRwx(buf []byte, mode fs.FileMode) {
	const rwx = "rwxrwxrwx"
	perms := uint32(mode.Perm())
	for i := 0; i < 9; i++ {
		if perms&(1<<(8-uint(i))) != 0 {
			buf[i] = rwx[i]
		} else {
			buf[i] = '-'
		}
	}
	if mode&fs.ModeSetuid != 0 {
		buf[2] = 's'
	}
	if mode&fs.ModeSetgid != 0 {
		buf[5] = 's'
	}
	if mode&fs.ModeSticky != 0 {
		buf[8] = 't'
	}
}

func kindFromMode(mode fs.FileMode) Kind {
	switch {
	case mode&fs.ModeSymlink != 0:
		return KindSymlink
	case mode&fs.ModeDir != 0:
		return KindDirectory
	case mode&fs.ModeCharDevice != 0:
		return KindCharDevice
	case mode&fs.ModeDevice != 0:
		return KindDevice
	case mode&fs.ModeNamedPipe != 0:
		return KindPipe
	case mode&fs.ModeSocket != 0:
		return KindSocket
	default:
		return KindFile
	}
}
