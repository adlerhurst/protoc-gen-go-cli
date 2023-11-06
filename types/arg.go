package types

import (
	"time"

	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	_ Arg = (*floatArg)(nil)
	_ Arg = (*doubleArg)(nil)
	_ Arg = (*int64Arg)(nil)
	_ Arg = (*uint64Arg)(nil)
	_ Arg = (*int32Arg)(nil)
	_ Arg = (*uint32Arg)(nil)
	_ Arg = (*boolArg)(nil)
	_ Arg = (*stringArg)(nil)
	_ Arg = (*messageArg)(nil)
	_ Arg = (*bytesArg)(nil)
	_ Arg = (*enumArg[int32])(nil)
	_ Arg = (*enumArg[string])(nil)
	_ Arg = (*durationArg)(nil)
	_ Arg = (*optionalArg)(nil)
	_ Arg = (*requiredArg)(nil)
	_ Arg = (*repeatedArg)(nil)
	_ Arg = (*oneOfArg)(nil)
)

type Arg interface {
	Name() string
}

type arg struct {
	ProtoName string
	name      string
}

func (a *arg) Name() string {
	return a.name
}

type floatArg struct {
	arg
	value float32
}

func NewFloatArg(name string, value float32) *floatArg {
	return &floatArg{
		arg: arg{
			name: name,
		},
		value: value,
	}
}

type doubleArg struct {
	arg
	value float64
}

func NewDoubleArg(name string, value float64) *doubleArg {
	return &doubleArg{
		arg: arg{
			name: name,
		},
		value: value,
	}
}

type int64Arg struct {
	arg
	value int64
}

func NewInt64Arg(name string, value int64) *int64Arg {
	return &int64Arg{
		arg: arg{
			name: name,
		},
		value: value,
	}
}

type uint64Arg struct {
	arg
	value uint64
}

func NewUint64Arg(name string, value uint64) *uint64Arg {
	return &uint64Arg{
		arg: arg{
			name: name,
		},
		value: value,
	}
}

type int32Arg struct {
	arg
	value int32
}

func NewInt32Arg(name string, value int32) *int32Arg {
	return &int32Arg{
		arg: arg{
			name: name,
		},
		value: value,
	}
}

type uint32Arg struct {
	arg
	value uint32
}

func NewUint32Arg(name string, value uint32) *uint32Arg {
	return &uint32Arg{
		arg: arg{
			name: name,
		},
		value: value,
	}
}

type boolArg struct {
	arg
	value bool
}

func NewBoolArg(name string, value bool) *boolArg {
	return &boolArg{
		arg: arg{
			name: name,
		},
		value: value,
	}
}

type stringArg struct {
	arg
	value string
}

func NewStringArg(name string, value string) *stringArg {
	return &stringArg{
		arg: arg{
			name: name,
		},
		value: value,
	}
}

type messageArg struct {
	arg
	value map[string]Arg
}

func NewMessageArg(name string, value map[string]Arg) *messageArg {
	return &messageArg{
		arg: arg{
			name: name,
		},
		value: value,
	}
}

type bytesArg struct {
	arg
	value []byte
}

func NewBytesArg(name string, value []byte) *bytesArg {
	return &bytesArg{
		arg: arg{
			name: name,
		},
		value: value,
	}
}

type enumValue interface {
	~int32 | string
}

// TODO: handle name and number
type enumArg[V enumValue] struct {
	arg
	descriptor protoreflect.Descriptor
	value      V
}

func NewEnumArg[V enumValue](name string, descriptor protoreflect.EnumDescriptor, value V) *enumArg[V] {
	return &enumArg[V]{
		arg: arg{
			name: name,
		},
		descriptor: descriptor,
		value:      value,
	}
}

type durationArg struct {
	arg
	value time.Duration
}

func NewDurationArg(name string, value time.Duration) *durationArg {
	return &durationArg{
		arg: arg{
			name: name,
		},
		value: value,
	}
}

type optionalArg struct {
	arg
	value Arg
}

func NewOptionalArg(name string, value Arg) *optionalArg {
	return &optionalArg{
		arg: arg{
			name: name,
		},
		value: value,
	}
}

type requiredArg struct {
	arg
	value Arg
}

func NewRequiredArg(name string, value Arg) *requiredArg {
	return &requiredArg{
		arg: arg{
			name: name,
		},
		value: value,
	}
}

type repeatedArg struct {
	arg
	value []Arg
}

func NewRepeatedArg(name string, value ...Arg) *repeatedArg {
	return &repeatedArg{
		arg: arg{
			name: name,
		},
		value: value,
	}
}

type oneOfArg struct {
	arg
	value Arg
}

func NewOneOfArg(name string, value Arg) *oneOfArg {
	return &oneOfArg{
		arg: arg{
			name: name,
		},
		value: value,
	}
}
