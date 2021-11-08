package srv

import (
	"os"
	"path/filepath"

	fileaccess "github.com/aserto-dev/aserto-idp-plugin-json/pkg/file-access"
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

	// check if file already exists
	path, err := os.Stat(c.File)

	if err != nil {

		dir := filepath.Dir(c.File)

		info, err := os.Stat(dir)
		if err != nil {
			return status.Error(codes.NotFound, err.Error())
		}

		if !info.IsDir() {
			return status.Errorf(codes.InvalidArgument, "%s is not a directory", dir)
		}

		err = fileaccess.WriteAccess(info, dir)
		if err != nil {
			return status.Errorf(codes.PermissionDenied, "cannot access %s", dir)
		}

	} else {
		err = fileaccess.WriteAccess(path, c.File)
		if err != nil {
			return status.Errorf(codes.PermissionDenied, "cannot access %s", c.File)
		}
	}

	return nil
}

func (c *JsonPluginConfig) Description() string {
	return "JSON plugin"
}
