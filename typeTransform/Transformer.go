package typeTransform

import (
	"fmt"
	"reflect"

	"github.com/bronze1man/kmg/kmgTime"
	"github.com/bronze1man/kmg/kmgType"
)

//golang的kind分的太细
type Kind uint

const (
	Invalid Kind = iota
	String
	Int
	Float
	Ptr
	Bool
	Time
	Interface
	Map
	Struct
	Slice
	Array
	Uint
)

func (k Kind) String() string {
	switch k {
	case Invalid:
		return "Invalid"
	case String:
		return "String"
	case Int:
		return "Int"
	case Float:
		return "Float"
	case Ptr:
		return "Ptr"
	case Bool:
		return "Bool"
	case Time:
		return "Time"
	case Interface:
		return "Interface"
	case Map:
		return "Map"
	case Struct:
		return "Struct"
	case Slice:
		return "Slice"
	case Array:
		return "Array"
	case Uint:
		return "Uint"
	default:
		return fmt.Sprintf("Kind: not defined %d", k)
	}
}

type TransformerFunc func(traner Transformer, in reflect.Value, out reflect.Value) (err error)
type Transformer map[Kind]map[Kind]TransformerFunc

func (t Transformer) Transform(in interface{}, out interface{}) (err error) {
	return t.Tran(reflect.ValueOf(in), reflect.ValueOf(out))
}
func (t Transformer) Tran(in reflect.Value, out reflect.Value) (err error) {
	iKind := GetReflectKind(in)
	oKind := GetReflectKind(out)

	//	fmt.Println(iKind,oKind,in.Type().String(),out.Type().String())
	m1, ok := t[iKind]
	if !ok {
		return fmt.Errorf("[typeTransform.tran] not support tran kind: [%s] to [%s]", in.Kind(), out.Kind())
	}
	m2, ok := m1[oKind]
	if !ok {
		return fmt.Errorf("[typeTransform.tran] not support tran kind: [%s] to [%s]", in.Kind(), out.Kind())
	}
	return m2(t, in, out)
}
func (t Transformer) Clone() Transformer {
	out1 := Transformer{}
	for inKind, m1 := range t {
		out2 := map[Kind]TransformerFunc{}
		for outKind, m2 := range m1 {
			out2[outKind] = m2
		}
		out1[inKind] = out2
	}
	return out1
}
func GetReflectKind(in reflect.Value) Kind {
	t := in.Type()
	if t == kmgType.DateTimeReflectType {
		return Time
	}
	switch t.Kind() {
	case reflect.String:
		return String
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Int
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr:
		return Uint
	case reflect.Float32, reflect.Float64:
		return Float
	case reflect.Ptr:
		return Ptr
	case reflect.Bool:
		return Bool
	case reflect.Interface:
		return Interface
	case reflect.Map:
		return Map
	case reflect.Struct:
		return Struct
	case reflect.Slice:
		return Slice
	case reflect.Array:
		return Array
	}
	panic(fmt.Errorf("not implement type %s", t.Kind().String()))
}

var DefaultTransformer = Transformer{
	Map: map[Kind]TransformerFunc{
		Map:       MapToMap,
		Struct:    MapToStruct,
		Ptr:       NonePtrToPtr,
		Interface: NoneInterfaceToInterface,
	},
	String: map[Kind]TransformerFunc{
		String:    StringToString,
		Int:       StringToInt,
		Uint:      StringToUint,
		Float:     StringToFloat,
		Bool:      StringToBool,
		Time:      NewStringToTimeFunc(kmgTime.DefaultTimeZone),
		Ptr:       NonePtrToPtr,
		Interface: NoneInterfaceToInterface,
	},
	Ptr: map[Kind]TransformerFunc{
		Ptr: PtrToPtr, //TODO reference to self..
	},
	Struct: map[Kind]TransformerFunc{
		Map:       StructToMap,
		Ptr:       NonePtrToPtr,
		Interface: NoneInterfaceToInterface,
	},
	Slice: map[Kind]TransformerFunc{
		Slice:     SliceToSlice,
		Ptr:       NonePtrToPtr,
		Interface: NoneInterfaceToInterface,
	},
	Interface: map[Kind]TransformerFunc{
		String: InterfaceToNoneInterface,
		Int:    InterfaceToNoneInterface,
		Float:  InterfaceToNoneInterface,
		Bool:   InterfaceToNoneInterface,
		Time:   InterfaceToNoneInterface,
		Struct: InterfaceToNoneInterface,
		Map:    InterfaceToNoneInterface,
		Ptr:    InterfaceToNoneInterface,
	},
	Int: map[Kind]TransformerFunc{
		String:    IntToString,
		Int:       IntToInt,
		Ptr:       NonePtrToPtr,
		Interface: NoneInterfaceToInterface,
	},
	Float: map[Kind]TransformerFunc{
		Int:       FloatToInt,
		Float:     FloatToFloat,
		Ptr:       NonePtrToPtr,
		String:    FloatToString,
		Interface: NoneInterfaceToInterface,
	},
	Time: map[Kind]TransformerFunc{
		String:    TimeToString,
		Ptr:       NonePtrToPtr,
		Interface: NoneInterfaceToInterface,
	},
}
