package srv

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/aserto-dev/aserto-idp-plugin-json/pkg/config"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/idp-plugin-sdk/pb"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	"github.com/hashicorp/go-multierror"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var jsonOptions = protojson.MarshalOptions{
	Multiline:       false,
	Indent:          "  ",
	AllowPartial:    true,
	UseProtoNames:   true,
	UseEnumNumbers:  false,
	EmitUnpopulated: false,
}

type JSONPlugin struct {
	Config   *config.JSONPluginConfig
	decoder  *json.Decoder
	users    bytes.Buffer
	op       plugin.OperationType
	apiUsers []*api.User
	count    int
}

func NewJSONPlugin() *JSONPlugin {
	return &JSONPlugin{
		Config: &config.JSONPluginConfig{},
	}
}

func (s *JSONPlugin) GetConfig() plugin.Config {
	return &config.JSONPluginConfig{}
}

func (s *JSONPlugin) GetVersion() (string, string, string) {
	return config.GetVersion()
}

func (s *JSONPlugin) Open(cfg plugin.Config, operation plugin.OperationType) error {
	conf, ok := cfg.(*config.JSONPluginConfig)
	if !ok {
		return errors.New("invalid config")
	}
	s.Config = conf
	s.count = 0

	s.op = operation
	switch operation {
	case plugin.OperationTypeWrite:

		s.users.Write([]byte("[\n"))

	case plugin.OperationTypeRead, plugin.OperationTypeDelete:

		file, err := os.Open(s.Config.FromFile)
		if err != nil {
			return err
		}

		s.decoder = json.NewDecoder(file)

		if _, err = s.decoder.Token(); err != nil {
			return err
		}
	}

	return nil
}

func (s *JSONPlugin) Read() ([]*api.User, error) {
	if s.decoder.More() {
		u := api.User{}
		if err := pb.UnmarshalNext(s.decoder, &u); err != nil {
			return nil, err
		}

		return []*api.User{&u}, nil
	}
	if _, err := s.decoder.Token(); err != nil {
		return nil, err
	}

	return nil, io.EOF
}

func (s *JSONPlugin) Write(user *api.User) error {
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

func (s *JSONPlugin) Delete(userID string) error {

	var err error
	if len(s.apiUsers) == 0 {
		err = s.readAll()
	}

	for _, user := range s.apiUsers {
		if user.Id == userID {
			user.Deleted = true
			user.Metadata.DeletedAt = timestamppb.New(time.Now())
		}
	}

	return err
}

func (s *JSONPlugin) Close() (*plugin.Stats, error) {
	switch s.op {
	case plugin.OperationTypeWrite, plugin.OperationTypeDelete:

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
	return nil, nil
}

func (s *JSONPlugin) readAll() error {
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
