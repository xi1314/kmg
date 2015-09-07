package kmgGoParser

import (
	"bytes"
	"fmt"
	"github.com/bronze1man/kmg/kmgGoSource/kmgGoReader"
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
	PackagePath string
	Name        string
	underType   Type //TODO 第二次扫描AST获取此信息
	Pkg         *Package
}

func (t *NamedType) GetKind() Kind {
	return Named
}

func (t *NamedType) GetUnderType() Type {
	if t.underType == nil {
		definer := t.Pkg.Program.GetNamedType(t.PackagePath, t.Name)
		t.underType = definer.underType
	}
	return t.underType
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

/*
第一个字符可能为
	letter -> identifier(单独的类型名,带package的类型的package部分)
	"struct" struct类型开头
	"func" func类型开头
	"interface" interface类型开头
	"*" 指针开头
	"[" (数组,slice) 开头
	"map[" map开头
	"chan " chan开头
	"chan<- " chan<- 开头
	"<-chan" chan<- 开头
*/
func (gofile *File) readType(r *kmgGoReader.Reader) Type {
	id := readIdentifier(r)
	if len(id) == 0 {
		if r.IsMatchAfter([]byte("<-chan")) {
			r.MustReadMatch([]byte("<-chan"))
			r.ReadAllSpace()
			return ChanType{
				Dir:  RecvDir,
				Elem: gofile.readType(r),
			}
		}
		b := r.ReadByte()
		if b == '*' {
			return PointerType{
				Elem: gofile.readType(r),
			}
		} else if b == '[' {
			content := readMatchMiddleParantheses(r)
			if len(content) == 1 {
				return SliceType{
					Elem: gofile.readType(r),
				}
			} else {
				// 仅跳过
				return ArrayType{
					Elem: gofile.readType(r),
				}
			}
		} else if b == '(' {
			typ := gofile.readType(r)
			r.MustReadMatch([]byte(")"))
			return typ
		} else {
			panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
		}
	}
	if bytes.Equal(id, []byte("struct")) {
		return gofile.readStruct(r)
	} else if bytes.Equal(id, []byte("interface")) {
		// 仅跳过
		r.ReadAllSpace()
		b := r.ReadByte()
		if b != '{' {
			panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
		}
		readMatchBigParantheses(r)
		return InterfaceType{}
	} else if bytes.Equal(id, []byte("map")) {
		b := r.ReadByte()
		m := MapType{}
		if b != '[' {
			panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
		}
		m.Key = gofile.readType(r)
		r.MustReadMatch([]byte("]"))
		m.Value = gofile.readType(r)
		return m
	} else if bytes.Equal(id, []byte("func")) {
		// 仅跳过
		r.ReadAllSpace()
		b := r.ReadByte()
		if b != '(' {
			panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
		}
		readMatchSmallParantheses(r) //跳过输入参数
		r.ReadAllSpaceWithoutLineBreak()
		run := r.ReadRune() //跳过输出参数
		if run == '(' {
			readMatchSmallParantheses(r)
		} else if run == '\n' { //换行符可以标识这个函数定义结束了.
			return FuncType{}
		} else if unicode.IsLetter(run) || run == '[' || run == '*' || run == '<' {
			r.UnreadRune() //输出参数只有一个类型
			gofile.readType(r)
		} else {
			r.UnreadRune() //读到了其他东西,退回.
		}
		return FuncType{}
	} else if bytes.Equal(id, []byte("chan")) {
		if r.IsMatchAfter([]byte("<-")) {
			r.MustReadMatch([]byte("<-"))
			r.ReadAllSpace()
			return ChanType{
				Dir:  SendDir,
				Elem: gofile.readType(r),
			}
		} else {
			r.ReadAllSpace()
			return ChanType{
				Dir:  BothDir,
				Elem: gofile.readType(r),
			}
		}
	} else {
		b := r.ReadByte()
		if b == '.' {
			pkgPath := string(id)
			pkgPath, err := gofile.LookupFullPackagePath(pkgPath)
			if err != nil {
				fmt.Println(r.GetFileLineInfo(), err.Error()) //TODO 以目前的复杂度暂时无解.需要把所有相关的package都看一遍才能正确.
			}
			id2 := readIdentifier(r)
			return &NamedType{
				PackagePath: pkgPath,
				Name:        string(id2),
				Pkg:         gofile.Pkg,
			}
		} else {
			r.UnreadByte()
			name := string(id)
			if builtinTypeMap[BuiltinType(name)] != Invalid {
				return BuiltinType(name)
			} else {
				return &NamedType{
					PackagePath: gofile.PackagePath,
					Name:        string(id),
					Pkg:         gofile.Pkg,
				}
			}
		}
	}
	/*
		}else if r.IsMatchAfter([]byte("struct")) { //TODO 解决struct后面必须有一个空格的问题.
			r.ReadAllSpace()
			b := r.ReadByte()
			if b!='{' {
				panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
			}
			readMatchBigParantheses(r)
			return StructType{}
		}else if r.IsMatchAfter([]byte("interface")) {
			r.ReadAllSpace()
			b := r.ReadByte()
			if b!='{' {
				panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
			}
			readMatchBigParantheses(r)
			return InterfaceType{}
		}else if r.IsMatchAfter([]byte("map[")){

		}else if r.IsMatchAfter([]byte("func")){
			//TODO finish it.
		}
	*/
}

func getTypeStructAnonymousName(typ Type) string {
	ntyp, ok := typ.(*NamedType)
	if ok {
		return ntyp.Name
	}
	ptyp, ok := typ.(PointerType)
	if ok {
		return "*" + getTypeStructAnonymousName(ptyp.Elem)
	}
	btyp, ok := typ.(BuiltinType)
	if ok {
		return string(btyp)
	}
	panic(fmt.Errorf("[getTypeStructAnonymousName] unexpect type %T", typ))
}

func (gofile *File) readStruct(r *kmgGoReader.Reader) StructType {
	// 仅跳过
	r.ReadAllSpace()
	b := r.ReadByte()
	if b != '{' {
		panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
	}
	lastReadBuf := []bytesAndType{}
	var lastTag []byte
	out := StructType{}
	for {
		r.ReadAllSpaceWithoutLineBreak()
		b := r.ReadByte()
		if b == '}' {
			return out
		} else if b == '"' || b == '\'' || b == '`' {
			r.UnreadByte()
			lastTag = mustReadGoString(r)
		} else if b == ',' {
			continue
		} else if b == '\n' {
			if len(lastReadBuf) == 0 {
				continue
			} else if len(lastReadBuf) == 1 {
				typ := lastReadBuf[0].typ
				name := getTypeStructAnonymousName(typ)
				out.Field = append(out.Field, StructField{
					Name:             name,
					Elem:             typ,
					IsAnonymousField: true,
					Tag:              string(lastTag),
				})
				lastReadBuf = []bytesAndType{}
			} else if len(lastReadBuf) >= 2 {
				typ := lastReadBuf[len(lastReadBuf)-1].typ
				for i := range lastReadBuf[:len(lastReadBuf)-1] {
					out.Field = append(out.Field, StructField{
						Name:             string(lastReadBuf[i].originByte),
						Elem:             typ,
						IsAnonymousField: false,
						Tag:              string(lastTag),
					})
				}
				lastReadBuf = []bytesAndType{}
			}
		} else {
			r.UnreadByte()
			startPos := r.Pos()
			typ := gofile.readType(r)
			lastReadBuf = append(lastReadBuf, bytesAndType{
				originByte: r.BufToCurrent(startPos),
				typ:        typ,
			})
		}
	}
}

type bytesAndType struct {
	originByte []byte
	typ        Type
}
