package example

import (
	"encoding/csv"
	"errors"
	"fmt"
	os "os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// func (x *CallRequest) ParseFlags(cmd *cobra.Command, args []string) {
// 	set := pflag.NewFlagSet("request", pflag.ContinueOnError)
// 	cmd.Flags().AddFlagSet(set)

// 	_UseFieldName := NewStringFlag(set, "use_field_name", ``, "")
// 	_UseCustomName := NewStringFlag(set, "custom", ``, "")
// 	_IsSomething := NewBoolFlag(set, "is_something", ``, false)
// 	_I32 := NewInt32Flag(set, "i32", ``, 0)
// 	_Ui32 := NewUint32Flag(set, "ui32", ``, 0)
// 	_I64 := NewInt64Flag(set, "i64", ``, 0)
// 	_Ui64 := NewUint64Flag(set, "ui64", ``, 0)
// 	_Fl := NewFloatFlag(set, "fl", ``, 0)
// 	_Dbl := NewDoubleFlag(set, "dbl", ``, 0)
// 	_Beiz := NewBytesFlag(set, "beiz", ``, nil)
// 	_Si32 := NewInt32Flag(set, "si32", ``, 0)
// 	_Si64 := NewInt64Flag(set, "si64", ``, 0)
// 	_F32 := NewUint32Flag(set, "f32", ``, 0)
// 	_F64 := NewUint64Flag(set, "f64", ``, 0)
// 	_Sf32 := NewInt32Flag(set, "sf32", ``, 0)
// 	_Sf64 := NewInt64Flag(set, "sf64", ``, 0)
// 	_RepS := NewDoubleSliceFlag(set, "rep_s", ``, nil)

// 	_Payload := NewStructFlag(set, "payload", ``)
// 	_CreatedAt := NewTimestampFlag(set, "created_at", ``)

// 	_Wat := NewEnumFlag[CallRequest_Wat](set, "wat", ``)
// 	_Some := NewEnumFlag[Some](set, "some", ``)
// 	_RepWat := NewEnumSliceFlag[CallRequest_Wat](set, "rep_wat", ``)

// 	flagIndexes := fieldIndexes(args, "nested", "rep_nest")

// 	// parse primitive flags before first nested
// 	if err := set.Parse(flagIndexes.primitives().args); err != nil {
// 		DefaultConfig.Logger.Error("failed to parse flags", "cause", err)
// 		os.Exit(1)
// 	}

// 	x.UseFieldName = *_UseFieldName.Value
// 	x.UseCustomName = *_UseCustomName.Value
// 	x.IsSomething = *_IsSomething.Value
// 	x.I32 = *_I32.Value
// 	x.Ui32 = *_Ui32.Value
// 	x.I64 = *_I64.Value
// 	x.Ui64 = *_Ui64.Value
// 	x.Fl = *_Fl.Value
// 	x.Dbl = *_Dbl.Value
// 	x.Beiz = *_Beiz.Value
// 	x.Si32 = *_Si32.Value
// 	x.Si64 = *_Si64.Value
// 	x.F32 = *_F32.Value
// 	x.F64 = *_F64.Value
// 	x.Sf32 = *_Sf32.Value
// 	x.Sf64 = *_Sf64.Value
// 	x.RepS = *_RepS.Value
// 	x.Payload = _Payload.Value
// 	x.CreatedAt = _CreatedAt.Value

// 	x.Wat = *_Wat.Value
// 	x.Some = *_Some.Value
// 	x.RepWat = *_RepWat.Value

// 	if flagIdx := flagIndexes.lastByName("nested"); flagIdx != nil {
// 		x.Nested = new(CallRequest_Nested)
// 		x.Nested.ParseFlags(flagIdx.args)
// 	}

// 	for _, idx := range flagIndexes.byName("rep_nest") {
// 		x.RepNest = append(x.RepNest, new(CallRequest_Nested))
// 		x.RepNest[len(x.RepNest)-1].ParseFlags(idx.args)
// 	}
// }

func (x *NestedRequest_Nested) ParseFlags(args []string) {
	set := pflag.NewFlagSet("nested", pflag.ContinueOnError)
	_Id := NewStringFlag(set, "id", ``)
	_Depth := NewInt32Flag(set, "depth", ``)

	if err := set.Parse(args); err != nil {
		DefaultConfig.Logger.Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	x.Id = *_Id.Value
	x.Depth = *_Depth.Value
}

func (x *CallRequest_Nested) ParseFlags(args []string) {
	set := pflag.NewFlagSet("nested", pflag.ContinueOnError)
	Field := NewStringFlag(set, "field", ``)

	if err := set.Parse(args); err != nil {
		DefaultConfig.Logger.Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	x.Field = *Field.Value
}

type argParser[T any] struct {
	primitiveParser[T]
	customParser func(field *T, arg string) error
}

type primitiveParserOpt[T any] func(*primitiveParser[T])

func WithDefaultValue[T any](value T) primitiveParserOpt[T] {
	return func(parser *primitiveParser[T]) {
		parser.defaultValue = value
	}
}

type primitiveParser[T any] struct {
	Value        *T
	defaultValue T
}

func (parser *primitiveParser[T]) applyOpts(opts []primitiveParserOpt[T]) {
	for _, opt := range opts {
		opt(parser)
	}
}

// Set implements pflag.Value.
func (v *argParser[T]) Set(arg string) error {
	if v.customParser != nil {
		return v.customParser(v.Value, arg)
	}

	value, ok := interface{}(v.Value).(protoreflect.ProtoMessage)
	if !ok {
		DefaultConfig.Logger.Error("must implement custom parser", "type", fmt.Sprintf("%T", v.Value))
	}
	return protojson.UnmarshalOptions{
		// AllowPartial: true,
		DiscardUnknown: true,
	}.Unmarshal([]byte(arg), value)
}

// String implements pflag.Value.
func (v *argParser[T]) String() string {
	value, ok := interface{}(v.Value).(protoreflect.ProtoMessage)
	if !ok {
		return fmt.Sprint(v.Value)
	}
	return protojson.Format(value)
}

// Type implements pflag.Value.
func (v *argParser[T]) Type() string {
	value, ok := interface{}(v.Value).(protoreflect.ProtoMessage)
	if !ok {
		return fmt.Sprintf("%T", v.Value)
	}
	return string(value.ProtoReflect().Type().Descriptor().FullName())
}

func NewStructFlag(set *pflag.FlagSet, name, usage string) (parser argParser[structpb.Struct]) {
	parser.Value = new(structpb.Struct)
	set.Var(&parser, name, usage)

	return parser
}

func NewStructSliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]*structpb.Struct]) (parser argParser[[]*structpb.Struct]) {
	parser.applyOpts(opts)
	parser.Value = new([]*structpb.Struct)
	set.Var(&parser, name, usage)
	return parser
}

func NewAnyFlag(set *pflag.FlagSet, name, usage string) (parser argParser[anypb.Any]) {
	// TODO: change to message
	parser.Value = new(anypb.Any)
	set.Var(&parser, name, usage)

	return parser
}

func NewTimestampFlag(set *pflag.FlagSet, name, usage string) (parser argParser[timestamppb.Timestamp]) {
	parser.Value = new(timestamppb.Timestamp)
	parser.customParser = timestampParser
	set.Var(&parser, name, usage)

	return parser
}

func NewTimestampSliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]*timestamppb.Timestamp]) (parser argParser[[]*timestamppb.Timestamp]) {
	parser.applyOpts(opts)
	parser.Value = new([]*timestamppb.Timestamp)
	parser.customParser = slicePtrParser[timestamppb.Timestamp](timestampParser)
	set.Var(&parser, name, usage)
	return parser
}

func timestampParser(field *timestamppb.Timestamp, arg string) error {
	t, err := time.Parse(time.RFC3339, arg)
	if err != nil {
		return err
	}
	*field = *timestamppb.New(t)
	return nil
}

func NewEnumFlag[E enum](set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[E]) (parser argParser[E]) {
	parser.applyOpts(opts)
	parser.Value = new(E)
	parser.customParser = enumParser[E]
	set.Var(&parser, name, usage)
	return parser
}

func NewEnumSliceFlag[E enum](set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]E]) (parser argParser[[]E]) {
	parser.applyOpts(opts)
	parser.Value = new([]E)
	parser.customParser = sliceParser(enumParser[E])
	set.Var(&parser, name, usage)
	return parser
}

func sliceParser[T any](parser func(*T, string) error) func(*[]T, string) error {
	return func(field *[]T, arg string) error {
		stringReader := strings.NewReader(arg)
		csvReader := csv.NewReader(stringReader)
		records, err := csvReader.Read()
		if err != nil {
			return err
		}

		values := make([]T, len(records))
		for i, record := range records {
			value := new(T)
			err := parser(value, record)
			if err != nil {
				return err
			}
			values[i] = *value
		}
		*field = append(*field, values...)

		return nil
	}
}

func slicePtrParser[T any](parser func(*T, string) error) func(*[]*T, string) error {
	return func(field *[]*T, arg string) error {
		stringReader := strings.NewReader(arg)
		csvReader := csv.NewReader(stringReader)
		records, err := csvReader.Read()
		if err != nil {
			return err
		}

		values := make([]*T, len(records))
		for i, record := range records {
			value := new(T)
			err := parser(value, record)
			if err != nil {
				return err
			}
			values[i] = value
		}
		*field = append(*field, values...)

		return nil
	}
}

func enumParser[E enum](field *E, arg string) error {
	if desc := (*field).Descriptor().Values().ByName(protoreflect.Name(arg)); desc != nil {
		*field = E(desc.Number())
		return nil
	}
	if number, err := strconv.Atoi(arg); err == nil {
		if desc := (*field).Descriptor().Values().ByNumber(protoreflect.EnumNumber(number)); desc != nil {
			*field = E(desc.Number())
			return nil
		}
	}

	return errors.New("unknown enum variable")
}

func NewStringFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[string]) (parser primitiveParser[string]) {
	parser.applyOpts(opts)
	parser.Value = new(string)
	set.StringVar(parser.Value, name, parser.defaultValue, usage)

	return parser
}

func NewStringSliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]string]) (parser primitiveParser[[]string]) {
	parser.applyOpts(opts)
	parser.Value = new([]string)
	set.StringSliceVar(parser.Value, name, parser.defaultValue, usage)

	return parser
}

func NewBoolFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[bool]) (parser primitiveParser[bool]) {
	parser.applyOpts(opts)
	parser.Value = new(bool)
	set.BoolVar(parser.Value, name, parser.defaultValue, usage)

	return parser
}

func NewBoolSliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]bool]) (parser primitiveParser[[]bool]) {
	parser.applyOpts(opts)
	parser.Value = new([]bool)
	set.BoolSliceVar(parser.Value, name, parser.defaultValue, usage)

	return parser
}

func NewInt32Flag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[int32]) (parser primitiveParser[int32]) {
	parser.applyOpts(opts)
	parser.Value = new(int32)
	set.Int32Var(parser.Value, name, parser.defaultValue, usage)

	return parser
}

func NewInt32SliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]int32]) (parser primitiveParser[[]int32]) {
	parser.applyOpts(opts)
	parser.Value = new([]int32)
	set.Int32SliceVar(parser.Value, name, parser.defaultValue, usage)

	return parser
}

func NewUint32Flag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[uint32]) (parser primitiveParser[uint32]) {
	parser.applyOpts(opts)
	parser.Value = new(uint32)
	set.Uint32Var(parser.Value, name, parser.defaultValue, usage)

	return parser
}

func NewUint32SliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]uint]) (parser primitiveParser[[]uint]) {
	return newUintSliceFlag(set, name, usage, opts...)
}

func NewInt64Flag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[int64]) (parser primitiveParser[int64]) {
	parser.applyOpts(opts)
	parser.Value = new(int64)
	set.Int64Var(parser.Value, name, parser.defaultValue, usage)

	return parser
}

func NewInt64SliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]int64]) (parser primitiveParser[[]int64]) {
	parser.applyOpts(opts)
	parser.Value = new([]int64)
	set.Int64SliceVar(parser.Value, name, parser.defaultValue, usage)

	return parser
}

func NewUint64Flag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[uint64]) (parser primitiveParser[uint64]) {
	parser.applyOpts(opts)
	parser.Value = new(uint64)
	set.Uint64Var(parser.Value, name, parser.defaultValue, usage)

	return parser
}

func NewUint64SliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]uint]) (parser primitiveParser[[]uint]) {
	return newUintSliceFlag(set, name, usage, opts...)
}

func newUintSliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]uint]) (parser primitiveParser[[]uint]) {
	parser.applyOpts(opts)
	parser.Value = new([]uint)
	set.UintSliceVar(parser.Value, name, parser.defaultValue, usage)

	return parser
}

func NewFloatFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[float32]) (parser primitiveParser[float32]) {
	parser.applyOpts(opts)
	parser.Value = new(float32)
	set.Float32Var(parser.Value, name, parser.defaultValue, usage)

	return parser
}

func NewFloatSliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]float32]) (parser primitiveParser[[]float32]) {
	parser.applyOpts(opts)
	parser.Value = new([]float32)
	set.Float32SliceVar(parser.Value, name, parser.defaultValue, usage)

	return parser
}

func NewDoubleFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[float64]) (parser primitiveParser[float64]) {
	parser.applyOpts(opts)
	parser.Value = new(float64)
	set.Float64Var(parser.Value, name, parser.defaultValue, usage)

	return parser
}

func NewDoubleSliceFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]float64]) (parser primitiveParser[[]float64]) {
	parser.applyOpts(opts)
	parser.Value = new([]float64)
	set.Float64SliceVar(parser.Value, name, parser.defaultValue, usage)

	return parser
}

func NewBytesFlag(set *pflag.FlagSet, name, usage string, opts ...primitiveParserOpt[[]byte]) (parser primitiveParser[[]byte]) {
	parser.applyOpts(opts)
	parser.Value = new([]byte)
	set.BytesBase64Var(parser.Value, name, parser.defaultValue, usage)

	return parser
}
