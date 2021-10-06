package srv

import (
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/sys/unix"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// values set by linker using ldflag -X
var (
	ver    string // nolint:gochecknoglobals // set by linker
	date   string // nolint:gochecknoglobals // set by linker
	commit string // nolint:gochecknoglobals // set by linker
)

func GetVersion() (string, string, string) {
	return ver, date, commit
}

type JsonPluginConfig struct {
	File string `description:"Json file path" kind:"attribute" mode:"normal" readonly:"false"`
}

func (c *JsonPluginConfig) Validate() error {

	if c.File == "" {
		return status.Error(codes.InvalidArgument, "no json file name was provided")
	}

	dir := filepath.Dir(c.File)

	info, err := os.Stat(dir)
	if err != nil {
		return status.Error(codes.NotFound, err.Error())
	}

	if !info.IsDir() {
		return status.Errorf(codes.InvalidArgument, "%s is not a directory", dir)
	}

	if runtime.GOOS == "windows" {
		if info.Mode().Perm()&(1<<(uint(7))) == 0 {
			return status.Errorf(codes.PermissionDenied, "cannot access %s", dir)
		}
	} else {
		err = unix.Access(dir, unix.W_OK)
		if err != nil {
			return status.Errorf(codes.PermissionDenied, "cannot access %s: %s", dir, err.Error())
		}
	}

	return nil
}

func (c *JsonPluginConfig) Description() string {
	return "JSON plugin"
}
