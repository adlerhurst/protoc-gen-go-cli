package types

import (
	"bytes"
	"log"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
)

func setPackage(gen *protogen.GeneratedFile, file *protogen.File) {
	gen.P("package ", file.GoPackageName)
}

func executeTemplate(gen *protogen.GeneratedFile, tmpl *template.Template, data any) error {
	var buffer bytes.Buffer

	if err := tmpl.Execute(&buffer, data); err != nil {
		log.Println("failed to execute template", err)
		return err
	}

	_, err := gen.Write(buffer.Bytes())
	if err != nil {
		log.Println("failed to write", err)
	}

	return err
}
