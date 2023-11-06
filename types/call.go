package types

import (
	_ "embed"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
)

type Call struct {
	Name             string
	ProtoName        string
	ShortDescription string
	LongDescription  string
	Args             []Arg
}

func (c *Call) NamePrivate() string {
	return strings.ToLower(c.Name[:1]) + c.Name[1:]
}

var (
	//go:embed call.template
	callDefinition string
	callTemplate   = template.Must(template.New("call").Parse(callDefinition))
)

func (c *Call) Generate(plugin *protogen.Plugin, file *protogen.File) error {
	gen := plugin.NewGeneratedFile(file.GeneratedFilenamePrefix /*+c.parent.name*/ +"_"+strings.ToLower(c.Name)+"_cli_call.pb.go", file.GoImportPath)

	setPackage(gen, file)
	c.imports(gen)

	if err := executeTemplate(gen, callTemplate, c); err != nil {
		return err
	}
	return nil
}

func (c *Call) imports(gen *protogen.GeneratedFile) {
	gen.QualifiedGoIdent(protogen.GoImportPath("os").Ident("os"))
	gen.QualifiedGoIdent(protogen.GoImportPath("github.com/spf13/cobra").Ident("cobra"))
	gen.QualifiedGoIdent(protogen.GoImportPath("google.golang.org/protobuf/encoding/protojson").Ident("protojson"))
}
