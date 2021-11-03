package srv

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
)

type JsonPlugin struct {
	Config      *JsonPluginConfig
	op          plugin.OperationType
	jsonContent []map[string]interface{}
}

func NewJsonPlugin() *JsonPlugin {
	return &JsonPlugin{
		Config: &JsonPluginConfig{},
	}
}

func (s *JsonPlugin) GetConfig() plugin.PluginConfig {
	return &JsonPluginConfig{}
}

func (s *JsonPlugin) Open(cfg plugin.PluginConfig, operation plugin.OperationType) error {
	config, ok := cfg.(*JsonPluginConfig)
	if !ok {
		return errors.New("invalid config")
	}
	s.Config = config
	s.op = operation
	switch operation {
	case plugin.OperationTypeWrite:
		{
			s.jsonContent = make([]map[string]interface{}, 0)
		}
	case plugin.OperationTypeRead, plugin.OperationTypeDelete:
		{
			file, err := ioutil.ReadFile(s.Config.File)
			if err != nil {
				return err
			}
			err = json.Unmarshal(file, &s.jsonContent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *JsonPlugin) Read() ([]*api.User, error) {
	var user api.User
	if len(s.jsonContent) == 0 {
		return nil, io.EOF
	}

	u, err := json.Marshal(s.jsonContent[0])
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(u, &user)
	if err != nil {
		return nil, err
	}

	s.jsonContent = s.jsonContent[1:]

	return []*api.User{&user}, nil
}

func (s *JsonPlugin) Write(user *api.User) error {
	var userInterface map[string]interface{}
	u, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = json.Unmarshal(u, &userInterface)
	if err != nil {
		return err
	}

	s.jsonContent = append(s.jsonContent, userInterface)

	return nil
}

func (s *JsonPlugin) Delete(userId string) error {
	for i := len(s.jsonContent) - 1; i >= 0; i-- {
		if s.jsonContent[i]["id"].(string) == userId {
			s.jsonContent[i]["deleted"] = true
			break
		}
	}
	return nil
}

func (s *JsonPlugin) Close() error {
	switch s.op {
	case plugin.OperationTypeWrite, plugin.OperationTypeDelete:
		{
			fileContent, err := json.Marshal(s.jsonContent)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(s.Config.File, fileContent, 0644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
