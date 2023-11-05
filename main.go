package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	option "github.com/adlerhurst/protoc-gen-go-cli/gen/proto/adlerhurst/cli/v1alpha"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type CLI struct {
	name     string
	commands []*Command
}

type Command struct {
	name    string
	methods []*Call
}

type Call struct {
	name string
	args []*Arg
}

type Arg struct {
	name string
	subs []*Arg
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}
}

func run() error {
	var opts protogen.Options

	if len(os.Args) > 1 {
		return fmt.Errorf("unknown argument %q (this program should be run by protoc, not directly)", os.Args[1])
	}
	in, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	req := &pluginpb.CodeGeneratorRequest{}
	if err := proto.Unmarshal(in, req); err != nil {
		return err
	}
	gen, err := opts.New(req)
	if err != nil {
		return err
	}
	for _, file := range gen.Files {
		if !file.Generate {
			continue
		}
		generateFile(gen, file)
	}

	resp := gen.Response()
	out, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	if _, err := os.Stdout.Write(out); err != nil {
		return err
	}

	return nil
}

const fileSuffix = "_cli.pb.go"

func generateFile(plugin *protogen.Plugin, protoFile *protogen.File) *protogen.GeneratedFile {
	g := plugin.NewGeneratedFile(protoFile.GeneratedFilenamePrefix+fileSuffix, protoFile.GoImportPath)

	g.P("// Code generated by protoc-gen-go-cli. DO NOT EDIT.")
	g.P()
	g.P("package ", protoFile.GoPackageName)
	g.P()

	f := parseFile(protoFile)

	g.P("// cli name: ", f.name)

	for _, svc := range f.commands {
		g.P("//  ", svc.name)
		for _, m := range svc.methods {
			g.P("//    ", m.name)
			for _, field := range m.args {
				g.P("//      --", field.name)
				for _, sub := range field.subs {
					g.P("//        --", sub.name)
					for _, sub := range sub.subs {
						g.P("//          --", sub.name)
						for _, sub := range sub.subs {
							g.P("//            --", sub.name)
							for _, sub := range sub.subs {
								g.P("//              --", sub.name)
								for _, sub := range sub.subs {
									g.P("//                --", sub.name)
								}
							}
						}
					}
				}
			}
		}
	}

	return g
}

func parseFile(protoFile *protogen.File) *CLI {
	opts := protoFile.Desc.Options().(*descriptorpb.FileOptions)

	cli := &CLI{
		commands: make([]*Command, 0, len(protoFile.Services)),
		name:     proto.GetExtension(opts, option.E_CliName).(string),
	}

	for _, svc := range protoFile.Services {
		cli.commands = append(cli.commands, parseService(svc))
	}

	return cli
}

func parseService(svc *protogen.Service) *Command {
	command := &Command{
		methods: make([]*Call, 0, len(svc.Methods)),
	}

	opts := svc.Desc.Options().(*descriptorpb.ServiceOptions)
	command.name = proto.GetExtension(opts, option.E_CommandName).(string)
	if command.name == "" {
		command.name = string(svc.Desc.Name())
	}

	for _, method := range svc.Methods {
		command.methods = append(command.methods, parseMethod(method))
	}

	return command
}

func parseMethod(m *protogen.Method) *Call {
	call := &Call{
		args: parseMessage(m.Input.Desc),
	}

	opts := m.Desc.Options().(*descriptorpb.MethodOptions)
	call.name = proto.GetExtension(opts, option.E_CallName).(string)
	if call.name == "" {
		call.name = string(m.Desc.Name())
	}

	return call
}

func parseMessage(input protoreflect.MessageDescriptor) []*Arg {
	args := make([]*Arg, 0, input.Fields().Len())
	for i := 0; i < input.Fields().Len(); i++ {
		args = append(args, parseField(input.Fields().Get(i)))
	}
	return args
}

func parseField(f protoreflect.FieldDescriptor) *Arg {
	arg := new(Arg)

	opts := f.Options().(*descriptorpb.FieldOptions)
	arg.name = proto.GetExtension(opts, option.E_ArgName).(string)
	if arg.name == "" {
		arg.name = string(f.JSONName())
	}

	if f.Message() == nil {
		return arg
	}

	if f.Message().Messages().Len() == 0 {
		return arg
	}

	for i := 0; i < f.Message().Messages().Len(); i++ {
		msg := f.Message().Messages().Get(i)
		arg.subs = append(arg.subs, parseMessage(msg)...)
	}

	return arg
}

// func parseSubFields(sub *protogen.Message) *Arg {
// 	arg := &Arg{
// 		name: string(sub.Desc.Name()),
// 		subs: make([]*Arg, 0, len(sub.Fields)),
// 	}

// 	for _, field := range sub.Fields {
// 		arg.subs = append(arg.subs, parseField(field))
// 	}

// 	return arg
// }
