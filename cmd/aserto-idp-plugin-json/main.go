package main

import (
	"log"

	"github.com/aserto-dev/aserto-idp-plugin-json/pkg/srv"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered. Error:\n", r)
		}
	}()

	options := &plugin.Options{
		Handler: &srv.JSONPlugin{},
	}

	err := plugin.Serve(options)
	if err != nil {
		log.Println(err.Error())
	}
}
