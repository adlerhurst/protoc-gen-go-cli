package main

import (
	"flag"
	"log"

	"github.com/adlerhurst/protoc-gen-go-cli/types/v3"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/pluginpb"
)

var registry protoregistry.Files

func main() {
	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		for _, file := range plugin.Files {
			if err := registry.RegisterFile(file.Desc); err != nil {
				log.Println("register failed", err)
			}

			if !file.Generate {
				continue
			}

			for _, svc := range file.Services {
				types.SetMessages(svc)

				service := types.NewService(svc)
				if err := service.Generate(plugin, file); err != nil {
					return err
				}
			}
			err := types.GenerateMessages(plugin, file)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
