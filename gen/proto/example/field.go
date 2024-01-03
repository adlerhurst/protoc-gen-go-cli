package example

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Field interface {
	AddFlag(set *pflag.FlagSet)
}

type field[V any] struct {
	Name  string
	Usage string
	Value *V

	defaultValue V
}

type PathField struct {
	value *string
}

func (f *PathField) AddFlag(set *pflag.FlagSet) {
	f.value = set.String("path", "", "path to the json file containing the request")
	if err := set.SetAnnotation("path", cobra.BashCompFilenameExt, []string{".json"}); err != nil {
		DefaultConfig.Logger.Error("failed to set path annotation", "cause", err)
	}
}

type StringField field[string]

func (f *StringField) AddFlag(set *pflag.FlagSet) {
	set.StringVar(f.Value, f.Name, f.defaultValue, f.Usage)
}

type StringSliceField field[[]string]

func (f *StringSliceField) AddFlag(set *pflag.FlagSet) {
	set.StringSliceVar(f.Value, f.Name, f.defaultValue, f.Usage)
}

type BoolField field[bool]

func (f *BoolField) AddFlag(set *pflag.FlagSet) {
	set.BoolVar(f.Value, f.Name, f.defaultValue, f.Usage)
}

type BoolSliceField field[[]bool]

func (f *BoolSliceField) AddFlag(set *pflag.FlagSet) {
	set.BoolSliceVar(f.Value, f.Name, f.defaultValue, f.Usage)
}

type Int32Field field[int32]

func (f *Int32Field) AddFlag(set *pflag.FlagSet) {
	set.Int32Var(f.Value, f.Name, f.defaultValue, f.Usage)
}

type Int32SliceField field[[]int32]

func (f *Int32SliceField) AddFlag(set *pflag.FlagSet) {
	set.Int32SliceVar(f.Value, f.Name, f.defaultValue, f.Usage)
}

type Uint32Field field[uint32]

func (f *Uint32Field) AddFlag(set *pflag.FlagSet) {
	set.Uint32Var(f.Value, f.Name, f.defaultValue, f.Usage)
}

type Uint32SliceField field[[]uint]

func (f *Uint32SliceField) AddFlag(set *pflag.FlagSet) {
	set.UintSliceVar(f.Value, f.Name, f.defaultValue, f.Usage)
}

type Int64Field field[int64]

func (f *Int64Field) AddFlag(set *pflag.FlagSet) {
	set.Int64Var(f.Value, f.Name, f.defaultValue, f.Usage)
}

type Int64SliceField field[[]int64]

func (f *Int64SliceField) AddFlag(set *pflag.FlagSet) {
	set.Int64SliceVar(f.Value, f.Name, f.defaultValue, f.Usage)
}

type Uint64Field field[uint64]

func (f *Uint64Field) AddFlag(set *pflag.FlagSet) {
	set.Uint64Var(f.Value, f.Name, f.defaultValue, f.Usage)
}

type Uint64SliceField field[[]uint]

func (f *Uint64SliceField) AddFlag(set *pflag.FlagSet) {
	set.UintSliceVar(f.Value, f.Name, f.defaultValue, f.Usage)
}

type FloatField field[float32]

func (f *FloatField) AddFlag(set *pflag.FlagSet) {
	set.Float32Var(f.Value, f.Name, f.defaultValue, f.Usage)
}

type FloatSliceField field[[]float32]

func (f *FloatSliceField) AddFlag(set *pflag.FlagSet) {
	set.Float32SliceVar(f.Value, f.Name, f.defaultValue, f.Usage)
}

type DoubleField field[float64]

func (f *DoubleField) AddFlag(set *pflag.FlagSet) {
	set.Float64Var(f.Value, f.Name, f.defaultValue, f.Usage)
}

type DoubleSliceField field[[]float64]

func (f *DoubleSliceField) AddFlag(set *pflag.FlagSet) {
	set.Float64SliceVar(f.Value, f.Name, f.defaultValue, f.Usage)
}

type BytesField field[[]byte]

func (f *BytesField) AddFlag(set *pflag.FlagSet) {
	set.BytesBase64Var(f.Value, f.Name, f.defaultValue, f.Usage+" base64 (RFC 4648) encoded")
}

type EnumField[E enum] field[E]

func (enum *EnumField[E]) AddFlag(set *pflag.FlagSet) {
	e := new(E)
	enum.Value = e
	set.Var(enum, enum.Name, enum.Usage)
}

// Set Implements pflag.Value
func (enum *EnumField[E]) Set(s string) (err error) {
	enum.Value, err = parseEnum[E](s)
	return err
}

// String Implements pflag.Value
func (enum *EnumField[E]) String() string {
	return (*enum.Value).String()
}

// Type Implements pflag.Value
func (enum *EnumField[E]) Type() string {
	return "enum"
}

// TODO: implement pflag.SliceValue
type EnumSliceField[E enum] field[[]E]

func (enum *EnumSliceField[E]) AddFlag(set *pflag.FlagSet) {
	e := new([]E)
	enum.Value = e
	set.Var(enum, enum.Name, enum.Usage)
}

// Set Implements pflag.Value
func (enum *EnumSliceField[E]) Set(s string) error {
	if s == "" {
		return nil
	}
	stringReader := strings.NewReader(s)
	csvReader := csv.NewReader(stringReader)
	records, err := csvReader.Read()
	if err != nil {
		return err
	}

	values := make([]E, len(records))

	for i, record := range records {
		e, err := parseEnum[E](record)
		if err != nil {
			return err
		}
		values[i] = *e
	}

	*enum.Value = append(*enum.Value, values...)

	return nil
}

// String Implements pflag.Value
func (enum *EnumSliceField[E]) String() string {
	if len(*enum.Value) == 0 {
		return ""
	}
	list := make([]string, len(*enum.Value))

	for i, e := range *enum.Value {
		list[i] = e.String()
	}

	return "[" + strings.Join(list, ",") + "]"
}

// Type Implements pflag.Value
func (enum *EnumSliceField[E]) Type() string {
	return "enum list"
}

func parseEnum[E enum](s string) (*E, error) {
	e := new(E)
	if desc := (*e).Descriptor().Values().ByName(protoreflect.Name(s)); desc != nil {
		*e = E(desc.Number())
		return e, nil
	}

	if number, err := strconv.Atoi(s); err == nil {
		if desc := (*e).Descriptor().Values().ByNumber(protoreflect.EnumNumber(number)); desc != nil {
			*e = E(desc.Number())
			return e, nil
		}
	}

	return nil, errors.New("unknown enum variable")
}

type message interface {
	String() string
}

type MessageField[M message] struct {
	field[M]
	fields []Field
}

func (message *MessageField[M]) AddFlag(set *pflag.FlagSet) {
	message.Value = new(M)
	// set.Var(message, message.Name, message.Usage)
	for _, field := range message.fields {
		field.AddFlag(set)
	}
	// subFlags := pflag.NewFlagSet(message.Name, pflag.ExitOnError)
	// for _, field := range message.fields {
	// 	field.AddFlag(subFlags)
	// }
	// set.AddFlagSet(subFlags)
}

// Set Implements pflag.Value
func (message *MessageField[M]) Set(s string) (err error) {
	// message.value, err = parseMnum[M](s)
	return err
}

// String Implements pflag.Value
func (message MessageField[M]) String() string {
	return fmt.Sprint(*message.Value)
}

// Type Implements pflag.Value
func (message MessageField[M]) Type() string {
	return "message"
}
