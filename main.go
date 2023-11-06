package main

import (
	"flag"
	"log"

	"github.com/adlerhurst/protoc-gen-go-cli/types/v2"
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
			_, err := generateFile(plugin, file)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func generateFile(plugin *protogen.Plugin, file *protogen.File) (gen []*protogen.GeneratedFile, err error) {
	for _, svc := range file.Services {
		service := types.ServiceFromProto(svc)
		serviceGen, err := service.Generate(plugin, file)
		if err != nil {
			return nil, err
		}
		gen = append(gen, serviceGen...)
	}

	return gen, nil
}
