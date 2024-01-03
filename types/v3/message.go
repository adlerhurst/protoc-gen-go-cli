package types

import (
	_ "embed"
	"log/slog"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var messages = Messages{}

type Messages map[protoreflect.FullName]*Message

type Message struct {
	Message *protogen.Message
	Flags   []*Flag
}

func (msg *Message) NestedFlagNames() string {
	flagNames := make([]string, 0, len(msg.Flags))
	for _, flag := range msg.Flags {
		if flag.Message == nil {
			continue
		}
		flagNames = append(flagNames, `"`+flag.FieldNamePrivate+`"`)
	}

	return strings.Join(flagNames, ", ")
}

func (msg *Message) HasMessageFlags() bool {
	for _, flag := range msg.Flags {
		if flag.Message != nil {
			return true
		}
	}
	return false
}

func SetMessages(service *protogen.Service) {
	for _, method := range service.Methods {
		setMessage(method.Input)
	}
}

var wellknownArgParsers = map[protoreflect.FullName]*Flag{
	"google.protobuf.Timestamp": {
		FieldName: "timestamppb.Timestamp",
		Kind:      "timestamp",
	},
	"google.protobuf.Duration": {
		FieldName: "durationpb.Duration",
		Kind:      "duration",
	},
	"google.protobuf.Struct": {
		FieldName: "structpb.Struct",
		Kind:      "struct",
	},
	"google.protobuf.Any": {
		FieldName: "anypb.Any",
		Kind:      "any",
	},
}

func setMessage(message *protogen.Message) {
	if _, ok := messages[message.Desc.FullName()]; ok {
		return
	}

	msg := &Message{Message: message, Flags: make([]*Flag, len(message.Fields))}
	messages[message.Desc.FullName()] = msg
	for i, field := range message.Fields {
		flag := &Flag{
			IsList:           field.Desc.IsList(),
			FieldName:        field.GoName,
			FieldNamePrivate: field.Desc.JSONName(),
			Name:             field.GoName,
			Kind:             field.Desc.Kind().String(),
			IsPtr:            field.Desc.HasOptionalKeyword(),
		}
		msg.Flags[i] = flag

		if field.Enum != nil {
			flag.Enum = &enumFlag{Type: field.Enum.GoIdent.GoName}
			continue
		}

		if field.Message != nil {
			flag.IsPtr = true
			if wellknownFlag, ok := wellknownArgParsers[field.Message.Desc.FullName()]; ok {
				flag.FieldName = wellknownFlag.FieldName
				flag.Kind = wellknownFlag.Kind
				continue
			}
			flag.Message = &messageFlag{
				Type: field.Message.GoIdent.GoName,
			}
			setMessage(field.Message)
			continue
		}

		if field.Oneof != nil {
			slog.Info("oneof", "field", field.GoName)
			for _, field := range field.Oneof.Fields {
				slog.Info("oneoffield", "name", field.GoIdent)
			}
		}
	}
}

var (
	//go:embed message.go.tmpl
	messageDefinition string
	messageTemplate   = template.Must(template.New("message").Parse(messageDefinition))
)

func GenerateMessages(plugin *protogen.Plugin, file *protogen.File) (err error) {
	if len(messages) == 0 {
		return nil
	}
	gen := plugin.NewGeneratedFile(file.GeneratedFilenamePrefix+"_cli_flags.go", file.GoImportPath)

	header(gen, file)
	messages.imports(gen)

	err = executeTemplate(gen, messageTemplate, messages)
	if err != nil {
		return err
	}

	return nil
}

func (Messages) imports(gen *protogen.GeneratedFile) {
	gen.QualifiedGoIdent(protogen.GoImportPath("os").Ident("os"))
	gen.QualifiedGoIdent(protogen.GoImportPath("github.com/spf13/pflag").Ident("pflag"))
}

type Field struct {
	*protogen.Field
}

func (field *Field) FlagConstructor() string {
	var builder strings.Builder

	builder.WriteString(field.VarName())
	builder.WriteString(" := ")

	builder.WriteString("New")
	builder.WriteString(title.String(field.Desc.Kind().String()))

	if field.Field.Desc.IsList() {
		builder.WriteString("Slice")
	}

	builder.WriteString("Flag")
	if field.Enum != nil {
		builder.WriteString("[")
		builder.WriteString(field.Enum.GoIdent.GoName)
		builder.WriteString("]")
	}
	builder.WriteString(`(set, "`)
	builder.WriteString(field.Desc.JSONName())
	builder.WriteString(`", "")`)

	return builder.String()
}

func (field *Field) FlagAssignment() string {
	var builder strings.Builder

	builder.WriteString("x.")
	builder.WriteString(field.GoName)
	builder.WriteString(" = *")
	builder.WriteString(field.VarName())

	return builder.String()
}

func (field *Field) VarName() string {
	return field.Desc.JSONName() + "Flag"
}
