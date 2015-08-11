package kmgGoParser

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgFile"
	"path/filepath"
	"strings"
	"unicode"
)

// 暂时忽略任何测试包. 此处不做编译检查,认为所有输入的go文件都是正常的.
func MustParsePackage(gopath string, pkgPath string) *Package {
	dirPath := filepath.Join(gopath, "src", pkgPath)
	pkg := &Package{
		ImportMap: map[string]bool{},
		PkgPath:   pkgPath,
	}
	for _, path := range kmgFile.MustReadDirFileOneLevel(dirPath) {
		if strings.HasSuffix(path, ".go") {
			pkg.mustAddFile(filepath.Join(dirPath, path))
		}
	}
	return pkg
}

//再多解析一个文件,path 是绝对路径
func (pkg *Package) mustAddFile(path string) {
	//fmt.Println(path)
	file := parseFile(pkg.PkgPath, path)
	for imp := range file.ImportMap {
		pkg.AddImport(imp)
	}
	for _, funcDecl := range file.FuncList {
		pkg.FuncList = append(pkg.FuncList, funcDecl)
	}
	for _, funcDecl := range file.MethodList {
		pkg.MethodList = append(pkg.MethodList, funcDecl)
	}
	for _, namedType := range file.NamedTypeList {
		pkg.NamedTypeList = append(pkg.NamedTypeList, namedType)
	}
}

func (pkg *Package) LookupNamedType(name string) *NamedType {
	for i := range pkg.NamedTypeList {
		if pkg.NamedTypeList[i].Name == name {
			return pkg.NamedTypeList[i]
		}
	}
	return nil
}

func parseFile(pkgPath string, path string) *File {
	gofile := &File{
		PackagePath:    pkgPath,
		ImportMap:      map[string]bool{},
		AliasImportMap: map[string]string{},
	}
	content := kmgFile.MustReadFile(path)
	posFile := newPosFile(path, content)
	content = goSourceRemoveComment(content, posFile)
	r := newReader(content, posFile)

	r.ReadAllSpace()
	r.MustReadMatch([]byte("package"))
	r.ReadUntilByte('\n')
	for {
		if r.IsEof() {
			return gofile //没有 import 正确情况
		}
		r.ReadAllSpace()
		if r.IsMatchAfter([]byte("import")) {
			gofile.readImport(r)
			continue
		}
		break
	}
	for {
		switch {
		case r.IsEof():
			return gofile
		case r.IsMatchAfter([]byte("func")):
			funcDecl := gofile.readGoFunc(r)
			if funcDecl.GetKind() == Func {
				gofile.FuncList = append(gofile.FuncList, funcDecl)
			} else {
				gofile.MethodList = append(gofile.MethodList, funcDecl)
			}
		case r.IsMatchAfter([]byte("type")):
			//r.ReadUntilByte('\n')
			gofile.readGoType(r)
		case r.IsMatchAfter([]byte("var")):
			gofile.readGoVar(r)
		case r.IsMatchAfter([]byte("const")):
			gofile.readGoConst(r)
		// 有一些没有分析的代码,里面可能包含import,此处先简单绕过.
		case r.IsMatchAfter([]byte(`"`)) || r.IsMatchAfter([]byte("`")):
			mustReadGoString(r)
			//fmt.Println(string(ret))
		case r.IsMatchAfter([]byte(`'`)):
			mustReadGoChar(r)
		default:
			r.ReadByte()
		}
	}
}

func readIdentifier(r *reader) []byte {
	buf := &bytes.Buffer{}
	if r.IsEof() {
		panic(r.GetFileLineInfo() + " unexcept EOF")
	}
	b := r.ReadRune()
	if b == '_' || unicode.IsLetter(b) {
		buf.WriteRune(b)
	} else {
		r.UnreadRune()
		return nil
	}
	for {
		if r.IsEof() {
			return buf.Bytes()
		}
		b := r.ReadRune()
		if b == '_' || unicode.IsLetter(b) || unicode.IsDigit(b) {
			buf.WriteRune(b)
		} else {
			r.UnreadRune()
			return buf.Bytes() // 不是Identifier的东西留个调用者处理
		}
	}
}

// 跳过 "{" "}",默认当前已经有第一层了(已经读入一个"{"了)
func readMatchBigParantheses(r *reader) []byte {
	return readMatchChar(r, '{', '}')
}

// 跳过 "[" "]",默认当前已经有第一层了(已经读入一个"["了)
func readMatchMiddleParantheses(r *reader) []byte {
	return readMatchChar(r, '[', ']')
}

// 跳过 "(" ")",默认当前已经有第一层了(已经读入一个"("了)
func readMatchSmallParantheses(r *reader) []byte {
	return readMatchChar(r, '(', ')')
}

func readMatchChar(r *reader, starter byte, ender byte) []byte {
	startPos := r.pos
	level := 1
	for {
		if r.IsEof() {
			panic(r.GetFileLineInfo() + " unexcept EOF")
		}
		b := r.ReadByte()
		if b == '"' || b == '`' {
			r.UnreadByte()
			mustReadGoString(r)
		} else if b == '\'' {
			r.UnreadByte()
			mustReadGoChar(r)
		} else if b == starter {
			level++
		} else if b == ender {
			level--
			if level == 0 {
				return r.buf[startPos:r.pos]
			}
		}
	}
}
