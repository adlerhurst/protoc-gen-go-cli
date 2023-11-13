package types

import (
	_ "embed"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
)

type Request struct {
	parent *Method
	*protogen.Message
	Args       []Field
	NestedArgs []NestedField
	gen        *protogen.GeneratedFile
}

func RequestFromProto(parent *Method, message *protogen.Message) *Request {
	req := &Request{
		parent:  parent,
		Message: message,
	}

	for _, msgField := range message.Fields {
		field := fieldFromProto(req, msgField)
		if field != nil {
			if nestedField, ok := field.(NestedField); ok {
				DefaultConfig.Logger.Info("found")
				req.NestedArgs = append(req.NestedArgs, nestedField)
				continue
			}
			req.Args = append(req.Args, field)
		}
	}

	return req
}

var (
	//go:embed request.tmpl
	requestDefinition string
	requestTemplate   = template.Must(template.New("request").Parse(requestDefinition))
)

func (request *Request) Generate(plugin *protogen.Plugin, file *protogen.File) ([]*protogen.GeneratedFile, error) {
	request.gen = plugin.NewGeneratedFile(request.filename(file.GeneratedFilenamePrefix), file.GoImportPath)

	header(request.gen, file)
	request.imports(request.gen)
	if err := executeTemplate(request.gen, requestTemplate, request); err != nil {
		return nil, err
	}

	return []*protogen.GeneratedFile{request.gen}, nil
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
	gen.QualifiedGoIdent(protogen.GoImportPath("os").Ident("os"))
}

func (request *Request) Public() string {
	return request.parent.Public() + "Request"
}

func (request *Request) name() string {
	return string(request.Desc.Name())
}
