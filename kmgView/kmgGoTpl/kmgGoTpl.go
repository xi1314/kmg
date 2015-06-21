package kmgGoTpl

import (
	"bytes"
	"fmt"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgStrings"
	"go/format"
	"path/filepath"
	"strings"
)

type scope int

const (
	currentScopeStatement  scope = 1
	currentScopeExpression scope = 2
	currentScopeTpl        scope = 3
)

func MustBuildTplOneFile(in []byte) (out []byte) {
	var transformer transformer
	return transformer.msutTransform(in)
}

func MustBuildTplInDir(path string) {
	pathList := kmgFile.MustGetAllFiles(path)
	for _, val := range pathList {
		if filepath.Ext(val) != ".gotpl" {
			continue
		}
		out := MustBuildTplOneFile(kmgFile.MustReadFile(val))
		outFilePath := kmgFile.PathTrimExt(val) + ".go"
		kmgFile.MustWriteFile(outFilePath, out)
	}
}

type transformer struct {
	in           []byte
	pos          int
	currentScope scope
	lastScopeBuf bytes.Buffer
	outBuf       bytes.Buffer
	bracesLevel  int //大括号的层数

	//函数是否横跨?>分析
	isLastFuncTokenWithoutMatchBraces bool
	isFuncOpenInScope                 bool //在当前这个scope里面是否有一个函数在里面打开
	hasFuncOpenBetweenScope           bool //是否在scope之间有一个函数正在打开?
	lastFuncBraceLevel                int
}

func (t *transformer) msutTransform(in []byte) []byte {
	t.in = in
	t.currentScope = currentScopeTpl
	t.lastFuncBraceLevel = -1
	for t.pos = 0; t.pos < len(t.in); t.pos++ {
		if t.isMatchString("<?=") {
			if t.currentScope != currentScopeTpl {
				panic("<? and ?> not match")
			}
			t.endTplScope()
			t.currentScope = currentScopeExpression
			t.pos += 2
			continue
		} else if t.isMatchString("<?") {
			if t.currentScope != currentScopeTpl {
				panic("<? and ?> not match")
			}
			t.endTplScope()
			t.currentScope = currentScopeStatement
			t.pos += 1
			continue
		} else if t.isMatchString("?>") {
			if t.currentScope == currentScopeTpl {
				panic("<? and ?> not match")
			}
			t.endNotTplScope()
			t.currentScope = currentScopeTpl
			t.pos += 1
			continue
		}
		if t.currentScope == currentScopeStatement {
			if t.isMatchString("func ") {
				if t.isLastFuncTokenWithoutMatchBraces {
					panic("func and { not match")
				}
				t.isLastFuncTokenWithoutMatchBraces = true
			} else if in[t.pos] == '{' {
				t.bracesLevel += 1
				if t.isLastFuncTokenWithoutMatchBraces {
					//函数的打开符
					t.isLastFuncTokenWithoutMatchBraces = false
					t.lastFuncBraceLevel = t.bracesLevel
					t.isFuncOpenInScope = true
				}
			} else if in[t.pos] == '}' {
				t.bracesLevel -= 1
				if t.bracesLevel < 0 {
					panic("{ and } not match")
				}
				if t.lastFuncBraceLevel == t.bracesLevel+1 {
					t.lastFuncBraceLevel = -1
					t.isFuncOpenInScope = false
					//函数的关闭符
					if t.hasFuncOpenBetweenScope {
						t.lastScopeBuf.WriteString("return _buf.Bytes()\n")
						t.hasFuncOpenBetweenScope = false
					}
				}
			}
		}
		t.lastScopeBuf.WriteByte(t.in[t.pos])
	}
	if t.currentScope == currentScopeTpl {
		t.endTplScope()
	}
	output := t.outBuf.Bytes()
	output = addImportBytes(output)
	f, err := format.Source(output)
	if err != nil {
		return output
	}
	return f
}

func (t *transformer) endTplScope() {
	if t.lastScopeBuf.Len() > 0 {
		t.outBuf.WriteString("_buf.WriteString(")
		s := t.lastScopeBuf.String()
		if !strings.Contains(s, "`") {
			t.outBuf.WriteString("`")
			t.outBuf.WriteString(s)
			t.outBuf.WriteString("`")
		} else {
			t.outBuf.WriteString(fmt.Sprintf("%#v", t.lastScopeBuf.String()))
		}
		t.outBuf.WriteString(")\n")
	}
	t.lastScopeBuf.Reset()
}

func (t *transformer) endNotTplScope() {
	if t.currentScope == currentScopeStatement {
		t.outBuf.WriteString(strings.TrimSpace(t.lastScopeBuf.String()))
		t.outBuf.WriteString("\n")
		if t.isFuncOpenInScope {
			t.isFuncOpenInScope = false
			t.hasFuncOpenBetweenScope = true
			t.outBuf.WriteString("var _buf bytes.Buffer\n")
		}
	} else if t.currentScope == currentScopeExpression {
		t.outBuf.WriteString("_buf.WriteString(")
		t.outBuf.WriteString(strings.TrimSpace(t.lastScopeBuf.String()))
		t.outBuf.WriteString(")\n")
	}
	t.lastScopeBuf.Reset()
}

func (t *transformer) isMatchString(token string) bool {
	return bytes.HasPrefix(t.in[t.pos:], []byte(token))
}

func addImportBytes(in []byte) (out []byte) {
	var isLastImportToken bool
	var isInImportParentheses bool
	var lastImportParenthesesPos int
	var hasFoundImport bool
	outBuf := &bytes.Buffer{}
	for pos := 0; pos < len(in); pos++ {
		if bytes.HasPrefix(in[pos:], []byte("\nimport")) {
			if isLastImportToken {
				panic("import and ( not match")
			}
			isLastImportToken = true
		} else if in[pos] == '(' {
			if isLastImportToken {
				if isInImportParentheses {
					panic("import ( and ) not match")
				}
				isInImportParentheses = true
				lastImportParenthesesPos = pos
				isLastImportToken = false
				hasFoundImport = true
			}
		} else if in[pos] == ')' {
			if isInImportParentheses {
				var readedImportPathList []string
				importPkgPath := bytes.Split(in[lastImportParenthesesPos+1:pos], []byte{'\n'})
				for _, p := range importPkgPath {
					p = bytes.TrimSpace(p)
					if len(p) == 0 {
						continue
					}
					readedImportPathList = append(readedImportPathList, string(p))
				}
				if !kmgStrings.IsInSlice(readedImportPathList, "\"bytes\"") {
					outBuf.WriteString("\"bytes\"\n")
				}
				isInImportParentheses = false
			}
		}
		outBuf.WriteByte(in[pos])
	}
	if !hasFoundImport {
		outBuf.Reset()
		isLastPackageToken := false
		for pos := 0; pos < len(in); pos++ {
			if bytes.HasPrefix(in[pos:], []byte("package ")) {
				isLastPackageToken = true
			} else if in[pos] == '\n' {
				if isLastPackageToken {
					outBuf.WriteString("\nimport (\n\"bytes\"\n)")
					isLastPackageToken = false
				}
			}
			outBuf.WriteByte(in[pos])
		}
	}
	return outBuf.Bytes()
}
