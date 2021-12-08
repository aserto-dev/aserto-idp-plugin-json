package srv

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FileContainsString(filePath string, word string) (bool, error) {
	content, err := ioutil.ReadFile(filePath)

	if err != nil {
		return false, err
	}

	return strings.Contains(string(content), word), nil
}

func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)

	return err == nil
}

func CreateTestApiUser(id, displayName, email string) *api.User {
	user := api.User{
		Id:          id,
		DisplayName: displayName,
		Email:       email,
		Picture:     "",
		Identities:  make(map[string]*api.IdentitySource),
		Attributes: &api.AttrSet{
			Properties:  &structpb.Struct{Fields: make(map[string]*structpb.Value)},
			Roles:       []string{},
			Permissions: []string{},
		},
		Applications: make(map[string]*api.AttrSet),
		Metadata: &api.Metadata{
			CreatedAt: timestamppb.New(time.Now()),
			UpdatedAt: timestamppb.New(time.Now()),
		},
	}
	return &user
}

func TestOpenForRead(t *testing.T) {
	assert := require.New(t)

	currentDir, err := os.Getwd()
	assert.Nil(err)

	filePath := filepath.Dir(currentDir)
	filePath = filepath.Join(filePath, "testing", "user.json")
	config := JsonPluginConfig{
		File: filePath,
	}
	JSONplugin := NewJsonPlugin()

	err = JSONplugin.Open(&config, plugin.OperationTypeRead)
	assert.Nil(err)
	assert.NotNil(JSONplugin.decoder, "the decoder shouldn't be nil")
}

func TestOpenForReadWithInvalidJson(t *testing.T) {
	assert := require.New(t)

	currentDir, err := os.Getwd()
	assert.Nil(err)

	filePath := filepath.Dir(currentDir)
	filePath = filepath.Join(filePath, "testing", "invalid.json")
	config := JsonPluginConfig{
		File: filePath,
	}
	JSONplugin := NewJsonPlugin()

	err = JSONplugin.Open(&config, plugin.OperationTypeRead)
	assert.NotNil(err)
	r, _ := regexp.Compile("invalid character .* looking for beginning of value")
	assert.Regexp(r, err.Error())
}

func TestOpenForReadWithInexistingFile(t *testing.T) {
	assert := require.New(t)

	config := JsonPluginConfig{
		File: "test.json",
	}
	JSONplugin := NewJsonPlugin()

	err := JSONplugin.Open(&config, plugin.OperationTypeRead)
	assert.NotNil(err)
	r, _ := regexp.Compile("open test.json: no such file or directory")
	assert.Regexp(r, err.Error())
}

func TestOpenForWriteInexistingFile(t *testing.T) {
	assert := require.New(t)

	config := JsonPluginConfig{
		File: "test.json",
	}
	JSONplugin := NewJsonPlugin()

	err := JSONplugin.Open(&config, plugin.OperationTypeWrite)
	assert.Nil(err)
}

func TestReadTwoUsers(t *testing.T) {
	assert := require.New(t)

	currentDir, err := os.Getwd()
	assert.Nil(err)

	filePath := filepath.Dir(currentDir)
	filePath = filepath.Join(filePath, "testing", "user.json")
	config := JsonPluginConfig{
		File: filePath,
	}
	JSONplugin := NewJsonPlugin()

	err = JSONplugin.Open(&config, plugin.OperationTypeRead)
	assert.Nil(err)

	user, err := JSONplugin.Read()
	assert.Nil(err)
	assert.Equal("Euan Garden", user[0].DisplayName)

	user, err = JSONplugin.Read()
	assert.Nil(err)
	assert.Equal("Chris Johnson [SALES]", user[0].DisplayName)

	_, err = JSONplugin.Read()
	assert.NotNil(err)
	assert.Equal(io.EOF, err)
}

func TestReadInvalidApiUser(t *testing.T) {
	assert := require.New(t)

	currentDir, err := os.Getwd()
	assert.Nil(err)

	filePath := filepath.Dir(currentDir)
	filePath = filepath.Join(filePath, "testing", "invalid-user.json")
	config := JsonPluginConfig{
		File: filePath,
	}
	JSONplugin := NewJsonPlugin()

	err = JSONplugin.Open(&config, plugin.OperationTypeRead)
	assert.Nil(err)

	_, err = JSONplugin.Read()
	assert.NotNil(err)
	r, _ := regexp.Compile("proto.* invalid value for enum type.*")
	assert.Regexp(r, err.Error())

	user, err := JSONplugin.Read()
	assert.Nil(err)
	assert.Equal("Chris Johnson [SALES]", user[0].DisplayName, "should successfully read the secound user")

	_, err = JSONplugin.Read()
	assert.NotNil(err)
	assert.Equal(io.EOF, err)
}

func TestDelete(t *testing.T) {
	assert := require.New(t)

	currentDir, err := os.Getwd()
	assert.Nil(err)

	filePath := filepath.Dir(currentDir)
	originalFilePath := filepath.Join(filePath, "testing", "user.json")
	copyFilePath := filepath.Join(filePath, "testing", "copy_user.json")

	bytesRead, err := ioutil.ReadFile(originalFilePath)
	assert.Nil(err)

	//Copy users to copy_user file
	err = ioutil.WriteFile(copyFilePath, bytesRead, 0755)
	assert.Nil(err)

	containDeleted, err := FileContainsString(copyFilePath, "deleted")
	assert.Nil(err)
	assert.False(containDeleted)

	config := JsonPluginConfig{
		File: copyFilePath,
	}
	JSONplugin := NewJsonPlugin()

	err = JSONplugin.Open(&config, plugin.OperationTypeDelete)
	assert.Nil(err)

	err = JSONplugin.Delete("dfdadc39-7335-404d-af66-c77cf13a15f8")
	assert.Nil(err)

	stats, err := JSONplugin.Close()
	assert.Nil(stats)
	assert.Nil(err)
	containDeleted, err = FileContainsString(copyFilePath, "deleted")
	assert.Nil(err)
	assert.True(containDeleted)

	err = os.Remove(copyFilePath)
	assert.Nil(err)
}

func TestWrite(t *testing.T) {
	assert := require.New(t)

	currentDir, err := os.Getwd()
	assert.Nil(err)

	filePath := filepath.Dir(currentDir)
	filePath = filepath.Join(filePath, "testing", "test.json")

	assert.False(FileExists(filePath))

	config := JsonPluginConfig{
		File: filePath,
	}
	JSONplugin := NewJsonPlugin()

	apiUser := CreateTestApiUser("1", "Test Name", "test@email.com")

	err = JSONplugin.Open(&config, plugin.OperationTypeWrite)
	assert.Nil(err)

	err = JSONplugin.Write(apiUser)
	assert.Nil(err)

	stats, err := JSONplugin.Close()
	assert.Nil(stats)
	assert.Nil(err)

	assert.True(FileExists(filePath))

	containName, err := FileContainsString(filePath, "Test Name")
	assert.Nil(err)
	assert.True(containName)

	err = os.Remove(filePath)
	assert.Nil(err)
}
