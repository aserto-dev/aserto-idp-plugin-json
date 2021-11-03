// +build !windows

package fileaccess

import (
	"io/fs"

	"golang.org/x/sys/unix"
)

func WriteAccess(info fs.FileInfo, file string) error {
	return unix.Access(file, unix.W_OK)
}
