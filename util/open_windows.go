// +build windows

package util

import (
	"syscall"
)

func Safe_Open(fname string, mode, perm int) syscall.Handle {
	fd, e := syscall.Open(fname, mode, uint32(perm))
	if e != nil {
		panic(e)
	}
	return fd
}
