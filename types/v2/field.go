package types

import (
	"strings"

	option "github.com/adlerhurst/protoc-gen-go-cli/gen/proto/adlerhurst/cli/v1alpha"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Field interface {
	Construct() string
	Name() string
}

type NestedField interface {
	Field
	IsNested()
	IsRepeated() bool
}

func fieldFromProto(req *Request, msgField *protogen.Field) Field {
	switch msgField.Desc.Kind() {
	// TODO: case protoreflect.GroupKind: when are these fields used?
	case protoreflect.EnumKind:
		return &EnumField{importField: importField{field: field{parent: req, Field: msgField}}}
	case protoreflect.BoolKind:
		return &BoolField{field: field{parent: req, Field: msgField}}
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return &Int32Field{field: field{parent: req, Field: msgField}}
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return &Uint32Field{field: field{parent: req, Field: msgField}}
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return &Int64Field{field: field{parent: req, Field: msgField}}
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return &Uint64Field{field: field{parent: req, Field: msgField}}
	case protoreflect.FloatKind:
		return &FloatField{field: field{parent: req, Field: msgField}}
	case protoreflect.DoubleKind:
		return &DoubleField{field: field{parent: req, Field: msgField}}
	case protoreflect.StringKind:
		return &StringField{field: field{parent: req, Field: msgField}}
	case protoreflect.BytesKind:
		return &BytesField{field: field{parent: req, Field: msgField}}
	case protoreflect.MessageKind:
		switch msgField.Desc.Message().FullName().Name() {
		case "Struct":
			return &StructField{field: field{parent: req, Field: msgField}}
		case "Timestamp":
			return &TimestampField{field: field{parent: req, Field: msgField}}
		case "Any":
			return &AnyField{field: field{parent: req, Field: msgField}}
		}

		msg := &MessageField{
			importField: importField{
				field: field{parent: req, Field: msgField},
			},
			fields: make([]Field, len(msgField.Message.Fields)),
		}

		if msgField.Desc.Message().FullName().Name() == "Any" {
			DefaultConfig.Logger.Info("was recursive type", "name", msgField.Desc.Message().FullName().Name())
			return nil
		}
		for i, field := range msgField.Message.Fields {
			msg.fields[i] = fieldFromProto(req, field)
		}

		return msg
	}
	return nil
}

type importField struct {
	field
}

type field struct {
	*protogen.Field
	parent *Request
}

func (f field) Name() string {
	name := proto.GetExtension(f.Desc.Options(), option.E_ArgName).(string)
	if name == "" {
		name = string(f.Desc.Name())
	}
	return name
}

func (f field) construct() string {
	f.Desc.Default()
	return "(set, \"" + f.Name() + "\", `" + string(f.Comments.Leading) + "`)"
	// return `Name: "` + f.Name() + "\", Usage: `" + string(f.Comments.Leading) + "`" + ", Value: &x." + f.GoName
}

func (f field) IsMessage() bool {
	return f.Desc.Kind() == protoreflect.MessageKind
}

type StringField struct {
	field
}

func (f *StringField) Construct() string {
	name := "StringFlag"
	if f.Desc.IsList() {
		name = "StringSliceFlag"
	}
	return "New" + name + f.field.construct()
}

type BoolField struct {
	field
}

func (f *BoolField) Construct() string {
	name := "BoolFlag"
	if f.Desc.IsList() {
		name = "BoolSliceFlag"
	}
	return "New" + name + f.field.construct()
}

type Int32Field struct {
	field
}

func (f *Int32Field) Construct() string {
	name := "Int32Flag"
	if f.Desc.IsList() {
		name = "Int32SliceFlag"
	}
	return "New" + name + f.field.construct()
}

type Uint32Field struct {
	field
}

func (f *Uint32Field) Construct() string {
	name := "Uint32Flag"
	if f.Desc.IsList() {
		name = "Uint32SliceFlag"
	}
	return "New" + name + f.field.construct()
}

type Int64Field struct {
	field
}

func (f *Int64Field) Construct() string {
	name := "Int64Flag"
	if f.Desc.IsList() {
		name = "Int64SliceFlag"
	}
	return "New" + name + f.field.construct()
}

type Uint64Field struct {
	field
}

func (f *Uint64Field) Construct() string {
	name := "Uint64Flag"
	if f.Desc.IsList() {
		name = "Uint64SliceFlag"
	}
	return "New" + name + f.field.construct()
}

type FloatField struct {
	field
}

func (f *FloatField) Construct() string {
	name := "FloatFlag"
	if f.Desc.IsList() {
		name = "FloatSliceFlag"
	}
	return "New" + name + f.field.construct()
}

type DoubleField struct {
	field
}

func (f *DoubleField) Construct() string {
	name := "DoubleFlag"
	if f.Desc.IsList() {
		name = "DoubleSliceFlag"
	}
	return "New" + name + f.field.construct()
}

type BytesField struct {
	field
}

func (f *BytesField) Construct() string {
	name := "BytesFlag"
	if f.Desc.IsList() {
		name = "BytesSliceFlag"
	}
	return "New" + name + f.field.construct()
}

type EnumField struct {
	importField
}

func (f *EnumField) Construct() string {
	name := "EnumFlag"
	if f.Desc.IsList() {
		name = "EnumSliceFlag"
	}
	return "New" + name + `[` + f.Field.Enum.GoIdent.GoName + `]` + f.field.construct()
}

type MessageField struct {
	importField
	fields []Field
}

func (f *MessageField) Construct() string {
	name := "MessageFlag"
	// TODO: slice
	// if f.Desc.IsList() {
	// 	name = "MessageSliceFlag"
	// }

	ident := f.Field.Message.GoIdent.GoName
	if f.parent.GoIdent.GoImportPath.String() != f.Field.Message.GoIdent.GoImportPath.String() {
		f.parent.gen.Import(f.Field.Message.GoIdent.GoImportPath)
		ident = f.parent.gen.QualifiedGoIdent(f.Field.Message.GoIdent)
	}
	fields := make([]string, len(f.fields))
	DefaultConfig.Logger.Info("fields", "fields", f.fields)
	for i, sub := range f.fields {
		fields[i] = "&" + sub.Construct()
	}
	return name + `[*` + ident + `]{
		field: field[*` + ident + `]{
			` + f.field.construct() + `,
		},
		fields: []Field{
` + strings.Join(fields, ", \n") + `},
	}`
}

func (*MessageField) IsNested() {}

func (f *MessageField) IsRepeated() bool {
	return f.Desc.IsList()
}

type StructField struct {
	field
}

func (f *StructField) Construct() string {
	name := "StructFlag"
	if f.Desc.IsList() {
		name = "StructSliceFlag"
	}
	return "New" + name + f.field.construct()
}

type TimestampField struct {
	field
}

func (f *TimestampField) Construct() string {
	name := "TimestampFlag"
	if f.Desc.IsList() {
		name = "TimestampSliceFlag"
	}
	return "New" + name + f.field.construct()
}

type AnyField struct {
	field
}

func (f *AnyField) Construct() string {
	name := "AnyFlag"
	if f.Desc.IsList() {
		name = "AnySliceFlag"
	}
	return "New" + name + f.field.construct()
}
