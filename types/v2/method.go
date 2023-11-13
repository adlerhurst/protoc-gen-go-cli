package types

import (
	_ "embed"
	"strings"
	"text/template"

	option "github.com/adlerhurst/protoc-gen-go-cli/gen/proto/adlerhurst/cli/v1alpha"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

type Method struct {
	parent *Service
	*protogen.Method
	Request *Request

	gen *protogen.GeneratedFile
}

func MethodFromProto(parent *Service, method *protogen.Method) *Method {
	m := &Method{
		parent: parent,
		Method: method,
	}
	m.Request = RequestFromProto(m, method.Input)

	return m
}

var (
	//go:embed method.tmpl
	methodDefinition string
	methodTemplate   = template.Must(template.New("method").Parse(methodDefinition))
)

func (method *Method) Generate(plugin *protogen.Plugin, file *protogen.File) (generated []*protogen.GeneratedFile, err error) {
	method.gen = plugin.NewGeneratedFile(method.filename(file.GeneratedFilenamePrefix), file.GoImportPath)
	generated = append(generated, method.gen)

	header(method.gen, file)
	method.imports()
	if err = executeTemplate(method.gen, methodTemplate, method); err != nil {
		return nil, err
	}

	requestGen, err := method.Request.Generate(plugin, file)
	if err != nil {
		return nil, err
	}
	generated = append(generated, requestGen...)

	return generated, nil
}

func (method *Method) filename(prefix string) string {
	var builder strings.Builder

	builder.WriteString(prefix)
	builder.WriteRune('_')
	builder.WriteString(string(method.parent.Desc.Name()))
	builder.WriteRune('_')
	builder.WriteString(string(method.Desc.Name()))
	builder.WriteString("_cli.pb.go")

	return builder.String()
}

func (method *Method) imports() {
	method.gen.QualifiedGoIdent(protogen.GoImportPath("github.com/spf13/cobra").Ident("cobra"))
	method.gen.QualifiedGoIdent(protogen.GoImportPath("github.com/spf13/pflag").Ident("pflag"))
}

func (method *Method) Use() string {
	return lower.String(method.name())
}

func (method *Method) Public() string {
	return method.parent.Public() + title.String(method.name())
}

func (method *Method) Short() string {
	return string(method.Comments.Leading)
}

func (method *Method) Long() string {
	return string(method.Comments.Leading) + string(method.Comments.Trailing)
}

func (method *Method) VarName() string {
	return method.Public() + "Cmd"
}

func (method *Method) name() string {
	name := proto.GetExtension(method.Desc.Options(), option.E_CallName).(string)
	if name == "" {
		name = string(method.Desc.Name())
	}
	return name
}
