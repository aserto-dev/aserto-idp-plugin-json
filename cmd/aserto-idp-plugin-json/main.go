package main

import (
	"log"

	"github.com/aserto-dev/idp-plugin-sdk/examples/json/pkg/srv"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
)

func main() {

	options := &plugin.PluginOptions{
		PluginHandler: &srv.JsonPlugin{},
	}

	err := plugin.Serve(options)
	if err != nil {
		log.Println(err.Error)
	}
}
