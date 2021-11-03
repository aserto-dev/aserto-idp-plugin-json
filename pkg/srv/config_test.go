package srv

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateWithEmptyFileName(t *testing.T) {
	assert := require.New(t)
	config := JsonPluginConfig{
		File: "",
	}
	err := config.Validate()

	assert.NotNil(err)
	r, _ := regexp.Compile("InvalidArgument desc = no json file name was provided")
	assert.Regexp(r, err.Error())
}

func TestValidateWithInexistentFileName(t *testing.T) {
	assert := require.New(t)
	config := JsonPluginConfig{
		File: "test",
	}
	err := config.Validate()

	assert.Nil(err)
}

func TestValidateWithInvalidPathToFile(t *testing.T) {
	assert := require.New(t)
	config := JsonPluginConfig{
		File: "testing/test.json",
	}
	err := config.Validate()

	assert.NotNil(err)
	r, _ := regexp.Compile("NotFound desc = stat testing: no such file or directory")
	assert.Regexp(r, err.Error())
}

func TestValidateWithInaccessibleExistingFile(t *testing.T) {
	assert := require.New(t)

	currentDir, err := os.Getwd()
	assert.Nil(err)

	filePath := filepath.Dir(currentDir)
	filePath = filepath.Join(filePath, "testing", "permission-denied.json")
	err = ioutil.WriteFile(filePath, []byte(""), 0444)
	assert.Nil(err)

	config := JsonPluginConfig{
		File: filePath,
	}
	err = config.Validate()

	assert.NotNil(err)
	r, _ := regexp.Compile(".*PermissionDenied desc = cannot access .*permission-denied.json")
	assert.Regexp(r, err.Error())

	err = os.Remove(filePath)
	assert.Nil(err)
}

func TestValidateWithFileInInvalidPath(t *testing.T) {
	assert := require.New(t)

	currentDir, err := os.Getwd()
	assert.Nil(err)

	filePath := filepath.Dir(currentDir)
	filePath = filepath.Join(filePath, "testing", "invalid.json", "test.json")

	config := JsonPluginConfig{
		File: filePath,
	}
	err = config.Validate()

	assert.NotNil(err)
	r, _ := regexp.Compile("InvalidArgument desc = .*invalid.json is not a directory")
	assert.Regexp(r, err.Error())
}

func TestDescription(t *testing.T) {
	assert := require.New(t)
	config := JsonPluginConfig{
		File: "test.json",
	}

	description := config.Description()

	assert.Equal("JSON plugin", description, "should return the description of the plugin")

}
