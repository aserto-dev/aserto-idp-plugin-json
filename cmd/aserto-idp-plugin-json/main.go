package main

import (
	"log"

	"github.com/aserto-dev/aserto-idp-plugin-json/pkg/srv"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
)

func main() {

	options := &plugin.Options{
		Handler: &srv.JSONPlugin{},
	}

	err := plugin.Serve(options)
	if err != nil {
		log.Println(err.Error())
	}
}
