//go:build darwin

package listing

import (
	"os"
	"syscall"
)

type syscallStatT = syscall.Stat_t

func linksOf(info os.FileInfo) uint64 {
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return 1
	}
	return uint64(st.Nlink)
}
