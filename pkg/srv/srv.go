package srv

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/idp-plugin-sdk/pb"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	"github.com/hashicorp/go-multierror"
	"google.golang.org/protobuf/encoding/protojson"
)

var jsonOptions = protojson.MarshalOptions{
	Multiline:       false,
	Indent:          "  ",
	AllowPartial:    true,
	UseProtoNames:   true,
	UseEnumNumbers:  false,
	EmitUnpopulated: false,
}

type JsonPlugin struct {
	Config   *JsonPluginConfig
	decoder  *json.Decoder
	users    bytes.Buffer
	op       plugin.OperationType
	apiUsers []*api.User
	count    int
}

func NewJsonPlugin() *JsonPlugin {
	return &JsonPlugin{
		Config: &JsonPluginConfig{},
	}
}

func (s *JsonPlugin) GetConfig() plugin.PluginConfig {
	return &JsonPluginConfig{}
}

func (s *JsonPlugin) GetVersion() (string, string, string) {
	return GetVersion()
}

func (s *JsonPlugin) Open(cfg plugin.PluginConfig, operation plugin.OperationType) error {
	config, ok := cfg.(*JsonPluginConfig)
	if !ok {
		return errors.New("invalid config")
	}
	s.Config = config
	s.count = 0

	s.op = operation
	switch operation {
	case plugin.OperationTypeWrite:
		{
			s.users.Write([]byte("[\n"))
		}
	case plugin.OperationTypeRead, plugin.OperationTypeDelete:
		{
			file, err := os.Open(s.Config.FromFile)
			if err != nil {
				return err
			}

			s.decoder = json.NewDecoder(file)

			if _, err = s.decoder.Token(); err != nil {
				return err
			}

		}
	}

	return nil
}

func (s *JsonPlugin) Read() ([]*api.User, error) {
	if s.decoder.More() {
		u := api.User{}
		if err := pb.UnmarshalNext(s.decoder, &u); err != nil {
			return nil, err
		}

		return []*api.User{&u}, nil
	} else {
		if _, err := s.decoder.Token(); err != nil {
			return nil, err
		}

		return nil, io.EOF
	}
}

func (s *JsonPlugin) Write(user *api.User) error {
	if s.count != 0 {
		_, _ = s.users.Write([]byte(",\n"))
	}
	b, err := jsonOptions.Marshal(user)
	if err != nil {
		return err
	}

	if _, err := s.users.Write(b); err != nil {
		return err
	}
	s.count++

	return nil
}

func (s *JsonPlugin) Delete(userId string) error {

	var err error
	if len(s.apiUsers) == 0 {
		err = s.readAll()
	}

	for _, user := range s.apiUsers {
		if user.Id == userId {
			user.Deleted = true
		}
	}

	return err
}

func (s *JsonPlugin) Close() (*plugin.Stats, error) {
	switch s.op {
	case plugin.OperationTypeWrite, plugin.OperationTypeDelete:
		{
			file := s.Config.ToFile
			if s.op == plugin.OperationTypeDelete {
				file = s.Config.FromFile
				s.users.Reset()
				s.users.Write([]byte("[\n"))

				for _, user := range s.apiUsers {
					err := s.Write(user)
					if err != nil {
						return nil, err
					}
				}
			}
			_, err := s.users.Write([]byte("\n]\n"))
			if err != nil {
				return nil, err
			}
			f, err := os.Create(file)
			if err != nil {
				return nil, err
			}
			w := bufio.NewWriter(f)
			_, err = s.users.WriteTo(w)
			if err != nil {
				return nil, err
			}
			w.Flush()
		}
	}
	return nil, nil
}

func (s *JsonPlugin) readAll() error {
	var errs error
	users, err := s.Read()
	for err != io.EOF {
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		s.apiUsers = append(s.apiUsers, users...)
		users, err = s.Read()
	}

	return errs
}
