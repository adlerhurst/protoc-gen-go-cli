package types

import (
	_ "embed"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
)

type Command struct {
	Name    string
	Service string

	ShortDescription string
	LongDescription  string

	Calls []*Call
}

var (
	//go:embed command.template
	commandDefinition string
	commandTemplate   = template.Must(template.New("command").Parse(commandDefinition))
)

func (command *Command) Generate(plugin *protogen.Plugin, file *protogen.File) error {
	gen := plugin.NewGeneratedFile(file.GeneratedFilenamePrefix+"_"+strings.ToLower(command.Name)+"_cli.pb.go", file.GoImportPath)

	setPackage(gen, file)
	command.imports(gen)
	if err := executeTemplate(gen, commandTemplate, command); err != nil {
		return err
	}

	for _, call := range command.Calls {
		if err := call.Generate(plugin, file); err != nil {
			return err
		}
	}

	return nil
}

func (*Command) imports(gen *protogen.GeneratedFile) {
	gen.QualifiedGoIdent(protogen.GoImportPath("os").Ident("os"))
	gen.QualifiedGoIdent(protogen.GoImportPath("github.com/spf13/cobra").Ident("cobra"))
	gen.QualifiedGoIdent(protogen.GoImportPath("github.com/spf13/viper").Ident("viper"))
	gen.QualifiedGoIdent(protogen.GoImportPath("google.golang.org/grpc").Ident("grpc"))
}
