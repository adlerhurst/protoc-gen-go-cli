package types

import "strings"

type Flag struct {
	Name             string
	FieldName        string
	FieldNamePrivate string
	Kind             string

	IsList bool
	IsPtr  bool

	Message *messageFlag
	Enum    *enumFlag
}

func (flag *Flag) FlagConstructor() string {
	var builder strings.Builder

	builder.WriteString(flag.VarName())
	builder.WriteString(" := ")

	builder.WriteString("New")
	builder.WriteString(title.String(flag.Kind))

	if flag.IsList {
		builder.WriteString("Slice")
	}

	builder.WriteString("Flag")
	if flag.Enum != nil {
		builder.WriteString("[")
		builder.WriteString(flag.Enum.Type)
		builder.WriteString("]")
	}
	builder.WriteString(`(set, "`)
	builder.WriteString(flag.FieldNamePrivate)
	builder.WriteString(`", "")`)

	return builder.String()
}

func (flag *Flag) VarName() string {
	return flag.FieldNamePrivate + "Flag"
}

type messageFlag struct {
	// ==.Message.GoIdent.GoName
	Type string
}

type enumFlag struct {
	// ==.Enum.GoIdent.GoName
	Type string
}
