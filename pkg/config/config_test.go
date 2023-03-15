package config_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/aserto-dev/aserto-idp-plugin-json/pkg/config"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	"github.com/stretchr/testify/require"
)

func TestValidateWriteWithEmptyFileName(t *testing.T) {
	assert := require.New(t)
	cfg := config.JSONPluginConfig{
		ToFile: "",
	}
	err := cfg.Validate(plugin.OperationTypeWrite)

	assert.NotNil(err)
	r := regexp.MustCompile("InvalidArgument desc = no json file 'to_file' name was provided")
	assert.Regexp(r, err.Error())
}

func TestValidateReadWithEmptyFileName(t *testing.T) {
	assert := require.New(t)
	cfg := config.JSONPluginConfig{
		FromFile: "",
	}
	err := cfg.Validate(plugin.OperationTypeRead)

	assert.NotNil(err)
	r := regexp.MustCompile("InvalidArgument desc = no json file 'from_file' name was provided")
	assert.Regexp(r, err.Error())
}

func TestValidateDeleteWithEmptyFileName(t *testing.T) {
	assert := require.New(t)
	cfg := config.JSONPluginConfig{
		FromFile: "",
		ToFile:   "test.txt",
	}
	err := cfg.Validate(plugin.OperationTypeDelete)

	assert.NotNil(err)
	r := regexp.MustCompile("InvalidArgument desc = no json file 'from_file' name was provided")
	assert.Regexp(r, err.Error())
}

func TestValidateWriteWithInexistentFileName(t *testing.T) {
	assert := require.New(t)
	cfg := config.JSONPluginConfig{
		ToFile: "test",
	}
	err := cfg.Validate(plugin.OperationTypeWrite)

	assert.Nil(err)
}

func TestValidateReadWithInexistentFileName(t *testing.T) {
	assert := require.New(t)
	cfg := config.JSONPluginConfig{
		FromFile: "test",
	}
	err := cfg.Validate(plugin.OperationTypeRead)

	assert.NotNil(err)
	assert.Equal("'test' file doesn't exists: stat test: no such file or directory", err.Error())
}

func TestValidateWriteWithInvalidPathToFile(t *testing.T) {
	assert := require.New(t)
	cfg := config.JSONPluginConfig{
		ToFile: "testing/test.json",
	}
	err := cfg.Validate(plugin.OperationTypeWrite)

	assert.NotNil(err)
	r := regexp.MustCompile("NotFound desc = stat testing: no such file or directory")
	assert.Regexp(r, err.Error())
}

func TestValidateWriteWithInaccessibleExistingFile(t *testing.T) {
	assert := require.New(t)

	currentDir, err := os.Getwd()
	assert.Nil(err)

	filePath := filepath.Dir(currentDir)
	filePath = filepath.Join(filePath, "testing", "permission-denied.json")
	err = os.WriteFile(filePath, []byte(""), 0444) //nolint:gosec
	assert.Nil(err)

	cfg := config.JSONPluginConfig{
		ToFile: filePath,
	}
	err = cfg.Validate(plugin.OperationTypeWrite)

	assert.NotNil(err)
	r := regexp.MustCompile(".*PermissionDenied desc = cannot access .*permission-denied.json")
	assert.Regexp(r, err.Error())

	err = os.Remove(filePath)
	assert.Nil(err)
}

func TestValidateReadWithInaccessibleExistingFile(t *testing.T) {
	assert := require.New(t)

	currentDir, err := os.Getwd()
	assert.Nil(err)

	filePath := filepath.Dir(currentDir)
	filePath = filepath.Join(filePath, "testing", "permission-denied.json")
	err = os.WriteFile(filePath, []byte(""), 0222) //nolint:gosec
	assert.Nil(err)

	cfg := config.JSONPluginConfig{
		FromFile: filePath,
	}
	err = cfg.Validate(plugin.OperationTypeRead)

	assert.NotNil(err)
	r := regexp.MustCompile(".*PermissionDenied desc = cannot access .*permission-denied.json")
	assert.Regexp(r, err.Error())

	err = os.Remove(filePath)
	assert.Nil(err)
}

func TestValidateWriteWithFileInInvalidPath(t *testing.T) {
	assert := require.New(t)

	currentDir, err := os.Getwd()
	assert.Nil(err)

	filePath := filepath.Dir(currentDir)
	filePath = filepath.Join(filePath, "testing", "invalid.json", "test.json")

	cfg := config.JSONPluginConfig{
		ToFile: filePath,
	}
	err = cfg.Validate(plugin.OperationTypeWrite)

	assert.NotNil(err)
	r := regexp.MustCompile("InvalidArgument desc = '.*invalid.json' is not a directory")
	assert.Regexp(r, err.Error())
}

func TestDescription(t *testing.T) {
	assert := require.New(t)
	cfg := config.JSONPluginConfig{}

	description := cfg.Description()

	assert.Equal("JSON plugin", description, "should return the description of the plugin")

}
