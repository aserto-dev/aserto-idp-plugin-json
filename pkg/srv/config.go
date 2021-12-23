package srv

import (
	"os"
	"path/filepath"

	fileaccess "github.com/aserto-dev/aserto-idp-plugin-json/pkg/file-access"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	"github.com/pkg/errors"
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
	FromFile string `description:"Json file path to read or delete from" kind:"attribute" mode:"normal" readonly:"false" name:"from_file"`
	ToFile   string `description:"Json file path to write to" kind:"attribute" mode:"normal" readonly:"false" name:"to_file"`
}

func (c *JsonPluginConfig) Validate(operation plugin.OperationType) error {

	switch operation {
	case plugin.OperationTypeWrite:
		// TODO accept stdout
		if c.ToFile == "" {
			return status.Error(codes.InvalidArgument, "no json file 'to_file' name was provided")
		}
		err := validateWrite(c.ToFile)
		if err != nil {
			return err
		}
	case plugin.OperationTypeRead:
		if c.FromFile == "" {
			return status.Error(codes.InvalidArgument, "no json file 'from_file' name was provided")
		}
		err := validateRead(c.FromFile)
		if err != nil {
			return err
		}
	case plugin.OperationTypeDelete:
		if c.FromFile == "" {
			return status.Error(codes.InvalidArgument, "no json file 'from_file' name was provided")
		}
		err := validateRead(c.FromFile)
		if err != nil {
			return err
		}
		err = validateWrite(c.FromFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *JsonPluginConfig) Description() string {
	return "JSON plugin"
}

func validateRead(file string) error {
	path, err := os.Stat(file)

	if err != nil {
		return errors.Wrapf(err, "'%s' file doesn't exists", file)
	}

	err = fileaccess.ReadAccess(path, file)
	if err != nil {
		return status.Errorf(codes.PermissionDenied, "cannot access '%s' for read", file)
	}

	return nil
}

func validateWrite(file string) error {

	// check if file already exists
	path, err := os.Stat(file)

	if err != nil {

		dir := filepath.Dir(file)

		info, err := os.Stat(dir)
		if err != nil {
			return status.Error(codes.NotFound, err.Error())
		}

		if !info.IsDir() {
			return status.Errorf(codes.InvalidArgument, "'%s' is not a directory", dir)
		}

		err = fileaccess.WriteAccess(info, dir)
		if err != nil {
			return status.Errorf(codes.PermissionDenied, "cannot access '%s' for write", dir)
		}

	} else {
		err = fileaccess.WriteAccess(path, file)
		if err != nil {
			return status.Errorf(codes.PermissionDenied, "cannot access '%s' for write", file)
		}
	}

	return nil
}
