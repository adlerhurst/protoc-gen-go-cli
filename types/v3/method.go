package types

import (
	_ "embed"
	"text/template"

	option "github.com/adlerhurst/protoc-gen-go-cli/gen/proto/adlerhurst/cli/v1alpha"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

type Method struct {
	Parent *Service
	*protogen.Method
	Request string
}

func MethodFromProto(parent *Service, method *protogen.Method) *Method {
	m := &Method{
		Parent:  parent,
		Method:  method,
		Request: method.Input.GoIdent.GoName,
	}

	return m
}

var (
	//go:embed method.tmpl
	methodDefinition string
	methodTemplate   = template.Must(template.New("method").Parse(methodDefinition))
)

func (method *Method) Generate(plugin *protogen.Plugin, gen *protogen.GeneratedFile) error {
	method.imports(gen)

	return executeTemplate(gen, methodTemplate, method)
}

func (method *Method) imports(gen *protogen.GeneratedFile) {
	gen.QualifiedGoIdent(protogen.GoImportPath("github.com/spf13/cobra").Ident("cobra"))
	gen.QualifiedGoIdent(protogen.GoImportPath("github.com/spf13/pflag").Ident("pflag"))
}

func (method *Method) Use() string {
	return lower.String(method.name())
}

func (method *Method) Public() string {
	return method.Parent.Public() + title.String(method.name())
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
