package kmgGoParser

import (
	"strconv"
	"unicode"
	"unicode/utf8"
)

//
type FuncOrMethodDeclaration struct {
	Name         string
	ReceiverType Type
	InParameter  []FuncParameter
	OutParameter []FuncParameter
}

func (t *FuncOrMethodDeclaration) GetKind() Kind {
	if t.ReceiverType == nil {
		return Func
	} else {
		return Method
	}
}

func (t *FuncOrMethodDeclaration) IsExport() bool {
	runeValue, _ := utf8.DecodeRuneInString(t.Name)
	return unicode.IsUpper(runeValue)
}

// TODO finish it.
type FuncType struct {
	InParameter  []FuncParameter
	OutParameter []FuncParameter
}

func (t FuncType) GetKind() Kind {
	return Func
}

type NamedType struct {
	PackagePath string // TODO 在全部ast读取出来之后,还要进行一次变换,把实际的PackagePath搞出来
	Name        string
	UnderType   Type //TODO 第二次扫描AST获取此信息
}

func (t *NamedType) GetKind() Kind {
	return Named
}

//TODO finish it
type StructType struct {
	Field []StructField
}

func (t StructType) GetKind() Kind {
	return Struct
}

type StructField struct {
	Name             string
	Elem             Type
	IsAnonymousField bool
	Tag              string
}

type MapType struct {
	Key   Type
	Value Type
}

func (t MapType) GetKind() Kind {
	return Map
}

type StringType struct {
}

type InterfaceType struct {
}

func (t InterfaceType) GetKind() Kind {
	return Interface
}

type PointerType struct {
	Elem Type
}

func (t PointerType) GetKind() Kind {
	return Ptr
}

func NewPointer(elem Type) PointerType {
	return PointerType{Elem: elem}
}

type Type interface {
	GetKind() Kind
}

type FuncParameter struct {
	Name       string
	Type       Type
	IsVariadic bool //是否有3个点
}

type SliceType struct {
	Elem Type
}

func (t SliceType) GetKind() Kind {
	return Slice
}

// TODO finish it
type ArrayType struct {
	Size int
	Elem Type
}

func (t ArrayType) GetKind() Kind {
	return Array
}

type ChanType struct {
	Dir  ChanDir
	Elem Type
}

func (t ChanType) GetKind() Kind {
	return Chan
}

// 内置的,没有package前缀的类型.
type BuiltinType string

func (t BuiltinType) GetKind() Kind {
	return builtinTypeMap[t]
}

func (t BuiltinType) String() string {
	return string(t)
}

var builtinTypeMap = map[BuiltinType]Kind{
	"bool":       Bool,
	"byte":       Uint8,
	"complex128": Complex128,
	"complex64":  Complex64,
	"error":      Interface, // TODO problem?
	"float32":    Float32,
	"float64":    Float64,
	"int":        Int,
	"int8":       Int8,
	"int16":      Int16,
	"int32":      Int32,
	"int64":      Int64,
	"rune":       Int32,
	"string":     String,
	"uint":       Uint,
	"uint8":      Uint8,
	"uint16":     Uint16,
	"uint32":     Uint32,
	"uint64":     Uint64,
	"uintptr":    Uintptr,
}

type Kind uint

const (
	Invalid Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	Array
	Chan
	Func
	Interface
	Map
	Ptr
	Slice
	String
	Struct
	UnsafePointer
	Method
	Named
)

func (k Kind) String() string {
	if int(k) < len(kindNames) {
		return kindNames[k]
	}
	return "kind" + strconv.Itoa(int(k))
}

var kindNames = []string{
	Invalid:       "invalid",
	Bool:          "bool",
	Int:           "int",
	Int8:          "int8",
	Int16:         "int16",
	Int32:         "int32",
	Int64:         "int64",
	Uint:          "uint",
	Uint8:         "uint8",
	Uint16:        "uint16",
	Uint32:        "uint32",
	Uint64:        "uint64",
	Uintptr:       "uintptr",
	Float32:       "float32",
	Float64:       "float64",
	Complex64:     "complex64",
	Complex128:    "complex128",
	Array:         "array",
	Chan:          "chan",
	Func:          "func",
	Interface:     "interface",
	Map:           "map",
	Ptr:           "ptr",
	Slice:         "slice",
	String:        "string",
	Struct:        "struct",
	UnsafePointer: "unsafe.Pointer",
	Method:        "method",
	Named:         "Named",
}

type ChanDir int

const (
	RecvDir ChanDir             = 1 << iota // <-chan
	SendDir                                 // chan<-
	BothDir = RecvDir | SendDir             // chan
)
