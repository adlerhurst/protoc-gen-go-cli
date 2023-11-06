package types

import (
	_ "embed"
	"log"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
)

var (
	//go:embed cli.template
	cliDefinition string
	cliTemplate   = template.Must(template.New("cli").Parse(cliDefinition))
)

// CLI represents an executable
type CLI struct {
	Name             string
	ShortDescription string
	LongDescription  string
	Commands         []*Command
}

func (cli *CLI) ShouldGenerate() bool {
	return cli.Name != ""
}

func (cli *CLI) Generate(plugin *protogen.Plugin, file *protogen.File) error {
	gen := plugin.NewGeneratedFile(cli.fileName(file.GeneratedFilenamePrefix), file.GoImportPath)

	setPackage(gen, file)
	cli.imports(gen)
	if err := executeTemplate(gen, cliTemplate, cli); err != nil {
		log.Fatal(err)
		return err
	}

	for _, command := range cli.Commands {
		if err := command.Generate(plugin, file); err != nil {
			return err
		}
	}

	return nil
}

func (*CLI) imports(gen *protogen.GeneratedFile) {
	gen.QualifiedGoIdent(protogen.GoImportPath("log/slog").Ident("slog"))
	gen.QualifiedGoIdent(protogen.GoImportPath("net").Ident("net"))
	gen.QualifiedGoIdent(protogen.GoImportPath("strconv").Ident("strconv"))
	gen.QualifiedGoIdent(protogen.GoImportPath("github.com/spf13/cobra").Ident("cobra"))
	gen.QualifiedGoIdent(protogen.GoImportPath("github.com/spf13/viper").Ident("viper"))
	gen.QualifiedGoIdent(protogen.GoImportPath("google.golang.org/grpc").Ident("grpc"))
	gen.QualifiedGoIdent(protogen.GoImportPath("google.golang.org/grpc/credentials/insecure").Ident("insecure"))
}

func (cli *CLI) fileName(prefix string) string {
	var builder strings.Builder
	builder.WriteString(prefix)
	if cli.Name != "" {
		builder.WriteRune('_')
		builder.WriteString(strings.ToLower(cli.Name))
	}
	builder.WriteString("_cli.pb.go")
	return builder.String()
}
