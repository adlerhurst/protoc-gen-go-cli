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
		msg := &MessageField{
			importField: importField{
				field: field{parent: req, Field: msgField},
			},
			fields: make([]Field, len(msgField.Message.Fields)),
		}

		if msgField.Desc.Message().FullName().Name() == "Struct" ||
			msgField.Desc.Message().FullName().Name() == "Any" {
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
	return `Name: "` + f.Name() + "\", Usage: `" + string(f.Comments.Leading) + "`"
}

type StringField struct {
	field
}

func (f *StringField) Construct() string {
	name := "StringField"
	if f.Desc.IsList() {
		name = "StringSliceField"
	}
	return name + `{` + f.field.construct() + `}`
}

type BoolField struct {
	field
}

func (f *BoolField) Construct() string {
	name := "BoolField"
	if f.Desc.IsList() {
		name = "BoolSliceField"
	}
	return name + `{` + f.field.construct() + `}`
}

type Int32Field struct {
	field
}

func (f *Int32Field) Construct() string {
	name := "Int32Field"
	if f.Desc.IsList() {
		name = "Int32SliceField"
	}
	return name + `{` + f.field.construct() + `}`
}

type Uint32Field struct {
	field
}

func (f *Uint32Field) Construct() string {
	name := "Int32Field"
	if f.Desc.IsList() {
		name = "Int32SliceField"
	}
	return name + `{` + f.field.construct() + `}`
}

type Int64Field struct {
	field
}

func (f *Int64Field) Construct() string {
	name := "Int64Field"
	if f.Desc.IsList() {
		name = "Int64SliceField"
	}
	return name + `{` + f.field.construct() + `}`
}

type Uint64Field struct {
	field
}

func (f *Uint64Field) Construct() string {
	name := "Uint64Field"
	if f.Desc.IsList() {
		name = "Uint64SliceField"
	}
	return name + `{` + f.field.construct() + `}`
}

type FloatField struct {
	field
}

func (f *FloatField) Construct() string {
	name := "FloatField"
	if f.Desc.IsList() {
		name = "FloatSliceField"
	}
	return name + `{` + f.field.construct() + `}`
}

type DoubleField struct {
	field
}

func (f *DoubleField) Construct() string {
	name := "DoubleField"
	if f.Desc.IsList() {
		name = "DoubleSliceField"
	}
	return name + `{` + f.field.construct() + `}`
}

type BytesField struct {
	field
}

func (f *BytesField) Construct() string {
	name := "BytesField"
	if f.Desc.IsList() {
		name = "BytesSliceField"
	}
	return name + `{` + f.field.construct() + `}`
}

type EnumField struct {
	importField
}

func (f *EnumField) Construct() string {
	name := "EnumField"
	if f.Desc.IsList() {
		name = "EnumSliceField"
	}
	return name + `[` + f.Field.Enum.GoIdent.GoName + `]{` + f.field.construct() + `}`
}

type MessageField struct {
	importField
	fields []Field
}

func (f *MessageField) Construct() string {
	name := "MessageField"
	// TODO: slice
	// if f.Desc.IsList() {
	// 	name = "MessageSliceField"
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
