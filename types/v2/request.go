package types

import (
	_ "embed"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Request struct {
	parent *Method
	*protogen.Message
	Args []Field
}

// func (req *Request) UnmarshalArgs(cmd *cobra.Command, args []string) {
// 	set := pflag.NewFlagSet("request", pflag.ContinueOnError)
// 	for _, field := range req.Args {
// 		field.AddFlag(set)
// 	}

// 	cmd.Flags().AddFlagSet(set)
// 	cmd.DisableFlagParsing = false
// 	if err := cmd.ParseFlags(args); err != nil {
// 		DefaultConfig.Logger.Error("failed to parse flags", "cause", err)
// 		os.Exit(1)
// 	}

// }

func RequestFromProto(parent *Method, message *protogen.Message) *Request {
	req := Request{
		parent:  parent,
		Message: message,
	}

	for _, msgField := range message.Fields {
		var f Field
		switch msgField.Desc.Kind() {
		// TODO: case protoreflect.GroupKind: when are these fields used?
		case protoreflect.EnumKind:
			f = &EnumField{field: field{Field: msgField}}
		case protoreflect.BoolKind:
			f = &BoolField{field: field{Field: msgField}}
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			f = &Int32Field{field: field{Field: msgField}}
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			f = &Uint32Field{field: field{Field: msgField}}
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			f = &Int64Field{field: field{Field: msgField}}
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			f = &Uint64Field{field: field{Field: msgField}}
		case protoreflect.FloatKind:
			f = &FloatField{field: field{Field: msgField}}
		case protoreflect.DoubleKind:
			f = &DoubleField{field: field{Field: msgField}}
		case protoreflect.StringKind:
			f = &StringField{field: field{Field: msgField}}
		case protoreflect.BytesKind:
			f = &BytesField{field: field{Field: msgField}}
		case protoreflect.MessageKind:

		}
		if f != nil {
			req.Args = append(req.Args, f)
		}
	}

	return &req
}

var (
	//go:embed request.tmpl
	requestDefinition string
	requestTemplate   = template.Must(template.New("request").Parse(requestDefinition))
)

func (request *Request) Generate(plugin *protogen.Plugin, file *protogen.File) ([]*protogen.GeneratedFile, error) {
	gen := plugin.NewGeneratedFile(request.filename(file.GeneratedFilenamePrefix), file.GoImportPath)

	header(gen, file)
	request.imports(gen)
	if err := executeTemplate(gen, requestTemplate, request); err != nil {
		return nil, err
	}

	return []*protogen.GeneratedFile{gen}, nil
}

func (request *Request) filename(prefix string) string {
	var builder strings.Builder

	builder.WriteString(prefix)
	builder.WriteRune('_')
	builder.WriteString(string(request.parent.parent.Desc.Name()))
	builder.WriteRune('_')
	builder.WriteString(string(request.parent.Desc.Name()))
	builder.WriteRune('_')
	builder.WriteString(string(request.Desc.Name()))
	builder.WriteString("_cli.pb.go")

	return builder.String()
}

func (*Request) imports(gen *protogen.GeneratedFile) {
	gen.QualifiedGoIdent(protogen.GoImportPath("github.com/spf13/cobra").Ident("cobra"))
	gen.QualifiedGoIdent(protogen.GoImportPath("github.com/spf13/pflag").Ident("pflag"))
	gen.QualifiedGoIdent(protogen.GoImportPath("os").Ident("os"))
}

func (request *Request) Public() string {
	return request.parent.Public() + "Request"
}

func (request *Request) name() string {
	return string(request.Desc.Name())
}
