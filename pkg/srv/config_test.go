package srv

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	"github.com/stretchr/testify/require"
)

func TestValidateWriteWithEmptyFileName(t *testing.T) {
	assert := require.New(t)
	config := JsonPluginConfig{
		ToFile: "",
	}
	err := config.Validate(plugin.OperationTypeWrite)

	assert.NotNil(err)
	r, _ := regexp.Compile("InvalidArgument desc = no json file 'to_file' name was provided")
	assert.Regexp(r, err.Error())
}

func TestValidateReadWithEmptyFileName(t *testing.T) {
	assert := require.New(t)
	config := JsonPluginConfig{
		FromFile: "",
	}
	err := config.Validate(plugin.OperationTypeRead)

	assert.NotNil(err)
	r, _ := regexp.Compile("InvalidArgument desc = no json file 'from_file' name was provided")
	assert.Regexp(r, err.Error())
}

func TestValidateDeleteWithEmptyFileName(t *testing.T) {
	assert := require.New(t)
	config := JsonPluginConfig{
		FromFile: "",
		ToFile:   "test.txt",
	}
	err := config.Validate(plugin.OperationTypeDelete)

	assert.NotNil(err)
	r, _ := regexp.Compile("InvalidArgument desc = no json file 'from_file' name was provided")
	assert.Regexp(r, err.Error())
}

func TestValidateWriteWithInexistentFileName(t *testing.T) {
	assert := require.New(t)
	config := JsonPluginConfig{
		ToFile: "test",
	}
	err := config.Validate(plugin.OperationTypeWrite)

	assert.Nil(err)
}

func TestValidateReadWithInexistentFileName(t *testing.T) {
	assert := require.New(t)
	config := JsonPluginConfig{
		FromFile: "test",
	}
	err := config.Validate(plugin.OperationTypeRead)

	assert.NotNil(err)
	assert.Equal("'test' file doesn't exists: stat test: no such file or directory", err.Error())
}

func TestValidateWriteWithInvalidPathToFile(t *testing.T) {
	assert := require.New(t)
	config := JsonPluginConfig{
		ToFile: "testing/test.json",
	}
	err := config.Validate(plugin.OperationTypeWrite)

	assert.NotNil(err)
	r, _ := regexp.Compile("NotFound desc = stat testing: no such file or directory")
	assert.Regexp(r, err.Error())
}

func TestValidateWriteWithInaccessibleExistingFile(t *testing.T) {
	assert := require.New(t)

	currentDir, err := os.Getwd()
	assert.Nil(err)

	filePath := filepath.Dir(currentDir)
	filePath = filepath.Join(filePath, "testing", "permission-denied.json")
	err = ioutil.WriteFile(filePath, []byte(""), 0444)
	assert.Nil(err)

	config := JsonPluginConfig{
		ToFile: filePath,
	}
	err = config.Validate(plugin.OperationTypeWrite)

	assert.NotNil(err)
	r, _ := regexp.Compile(".*PermissionDenied desc = cannot access .*permission-denied.json")
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
	err = ioutil.WriteFile(filePath, []byte(""), 0222)
	assert.Nil(err)

	config := JsonPluginConfig{
		FromFile: filePath,
	}
	err = config.Validate(plugin.OperationTypeRead)

	assert.NotNil(err)
	r, _ := regexp.Compile(".*PermissionDenied desc = cannot access .*permission-denied.json")
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

	config := JsonPluginConfig{
		ToFile: filePath,
	}
	err = config.Validate(plugin.OperationTypeWrite)

	assert.NotNil(err)
	r, _ := regexp.Compile("InvalidArgument desc = '.*invalid.json' is not a directory")
	assert.Regexp(r, err.Error())
}

func TestDescription(t *testing.T) {
	assert := require.New(t)
	config := JsonPluginConfig{}

	description := config.Description()

	assert.Equal("JSON plugin", description, "should return the description of the plugin")

}
