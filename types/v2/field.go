package types

import (
	option "github.com/adlerhurst/protoc-gen-go-cli/gen/proto/adlerhurst/cli/v1alpha"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

type Field interface {
	Construct() string
}

type field struct {
	*protogen.Field
}

func (f field) name() string {
	name := proto.GetExtension(f.Desc.Options(), option.E_ArgName).(string)
	if name == "" {
		name = string(f.Desc.Name())
	}
	return name
}

func (f field) construct() string {
	return `Name: "` + f.name() + `", Usage: "` + string(f.Comments.Leading) + `"`
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
	field
}

func (f *EnumField) Construct() string {
	name := "EnumField"
	if f.Desc.IsList() {
		name = "EnumSliceField"
	}
	return name + `[` + f.Field.Enum.GoIdent.GoName + `]{` + f.field.construct() + `}`
}
