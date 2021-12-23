// +build windows

package fileaccess

import (
	"errors"
	"io/fs"
)

func WriteAccess(info fs.FileInfo, file string) error {
	if info.Mode().Perm()&(1<<(uint(2))) == 0 &&
		info.Mode().Perm()&(1<<(uint(3))) == 0 &&
		info.Mode().Perm()&(1<<(uint(6))) == 0 &&
		info.Mode().Perm()&(1<<(uint(7))) == 0 {
		return errors.New("cannot access for write")
	}

	return nil
}

func ReadAccess(info fs.FileInfo, file string) error {
	if info.Mode().Perm()&(1<<(uint(4))) == 0 &&
		info.Mode().Perm()&(1<<(uint(5))) == 0 &&
		info.Mode().Perm()&(1<<(uint(6))) == 0 &&
		info.Mode().Perm()&(1<<(uint(7))) == 0 {
		return errors.New("cannot access for read")
	}

	return nil
}
