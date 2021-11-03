// +build windows

package fileaccess

import (
	"errors"
	"io/fs"
)

func WriteAccess(info fs.FileInfo, file string) error {
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		return errors.New("cannot access")
	}

	return nil
}
