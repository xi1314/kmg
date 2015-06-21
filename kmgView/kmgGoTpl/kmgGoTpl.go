package kmgGoTpl

import (
	"bytes"
	"fmt"
	"github.com/bronze1man/kmg/kmgCache"
	"github.com/bronze1man/kmg/kmgErr"
	"github.com/bronze1man/kmg/kmgFile"
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

func BuildTplOneFile(in []byte, isHtml bool) (out []byte, err error) {
	var transformer transformer
	transformer.isHtml = isHtml
	err = kmgErr.PanicToError(func() {
		out = transformer.mustTransform(in)
	})
	if err != nil {
		return nil, fmt.Errorf(":%d %s", transformer.lineNum, err)
	}
	return out, err
}

func MustBuildTplInDir(path string) {
	pathList := kmgFile.MustGetAllFiles(path)
	for _, val := range pathList {
		ext := filepath.Ext(val)
		if ext == ".gotpl" {
			out, err := BuildTplOneFile(kmgFile.MustReadFile(val), false)
			if err != nil {
				panic(fmt.Sprintf("%s %s", val, err.Error()))
			}
			outFilePath := kmgFile.PathTrimExt(val) + ".go"
			kmgFile.MustWriteFile(outFilePath, out)
		} else if ext == ".gotplhtml" {
			out, err := BuildTplOneFile(kmgFile.MustReadFile(val), true)
			if err != nil {
				panic(fmt.Sprintf("%s %s", val, err.Error()))
			}
			outFilePath := kmgFile.PathTrimExt(val) + ".go"
			kmgFile.MustWriteFile(outFilePath, out)
		}
	}
}

func MustBuildTplInDirWithCache(path string) {
	kmgCache.MustMd5FileChangeCache("kmgGoTpl_"+path, []string{path}, func() {
		MustBuildTplInDir(path)
	})
}

type transformer struct {
	isHtml bool

	in           []byte
	pos          int
	lineNum      int
	currentScope scope
	lastScopeBuf bytes.Buffer
	outBuf       bytes.Buffer
	bracesLevel  int //大括号的层数

	//函数是否横跨?>分析
	isLastFuncTokenWithoutMatchBraces bool
	isFuncOpenInScope                 bool //在当前这个scope里面是否有一个函数在里面打开
	hasFuncOpenBetweenScope           bool //是否在scope之间有一个函数正在打开?
	lastFuncBraceLevel                int

	hasBytesPackage  bool
	hasKmgXssPackage bool

	// xss 目前简单采用排除法解决这个问题,很少使用的不考虑,并且信任写模板的人. 不是script,不是urlv 就是kmgXss.H
	isInScript        bool
	isLastScriptToken bool
	urlvStatus        urlvStatus
}

type urlvStatus int

const (
	urlvStatusNot      urlvStatus = 0
	urlvStatusQuestion urlvStatus = 1
	urlvStatusKey      urlvStatus = 2
	urlvStatusEqual    urlvStatus = 3
	urlvStatusValue    urlvStatus = 4
	urlvStatusAndSign  urlvStatus = 5
)

func (t *transformer) mustTransform(in []byte) []byte {
	t.in = in
	t.currentScope = currentScopeTpl
	t.lastFuncBraceLevel = -1
	t.lineNum = 1
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
			t.hasBytesPackage = true
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
						t.lastScopeBuf.WriteString("return _buf.String()\n")
						t.hasFuncOpenBetweenScope = false
					}
				}
			}
		} else if t.currentScope == currentScopeTpl {
			if t.isHtml {
				// script
				if t.isMatchString("<script") { //暂时只支持小写
					if t.isLastScriptToken {
						panic("<script and > not match")
					}
					t.isLastScriptToken = true
				} else if t.isMatchString("</script>") {
					t.isInScript = false
				} else if in[t.pos] == '>' {
					if t.isLastScriptToken {
						t.isLastScriptToken = false
						t.isInScript = true
					}
				}
				// urlv
				if in[t.pos] == '?' {
					// 忽略开头的状态错误
					t.urlvStatus = urlvStatusQuestion
				} else if t.urlvStatus == urlvStatusQuestion { //此处恰好是对的,后面不会匹配到?
					if isAlphanum(in[t.pos]) {
						t.urlvStatus = urlvStatusKey
					} else {
						t.urlvStatus = urlvStatusNot //状态计算错误,忽略本次匹配
					}
				} else if t.urlvStatus == urlvStatusKey {
					if isAlphanum(in[t.pos]) {
						t.urlvStatus = urlvStatusKey
					} else if in[t.pos] == '=' {
						t.urlvStatus = urlvStatusEqual
					} else {
						t.urlvStatus = urlvStatusNot //状态计算错误,忽略本次匹配
					}
				} else if t.urlvStatus == urlvStatusEqual {
					if isAlphanum(in[t.pos]) {
						t.urlvStatus = urlvStatusValue
					} else {
						t.urlvStatus = urlvStatusNot //状态计算错误,忽略本次匹配
					}
				} else if t.urlvStatus == urlvStatusValue {
					if isAlphanum(in[t.pos]) {
						t.urlvStatus = urlvStatusValue
					} else if in[t.pos] == '&' {
						t.urlvStatus = urlvStatusAndSign
					} else {
						t.urlvStatus = urlvStatusNot //状态计算错误,忽略本次匹配
					}
				} else if t.urlvStatus == urlvStatusAndSign {
					if isAlphanum(in[t.pos]) {
						t.urlvStatus = urlvStatusKey
					} else {
						t.urlvStatus = urlvStatusNot //状态计算错误,忽略本次匹配
					}
				}
			}
		}
		if in[t.pos] == '\n' {
			t.lineNum++
		}
		t.lastScopeBuf.WriteByte(t.in[t.pos])
	}
	if t.currentScope == currentScopeTpl {
		s := strings.TrimSpace(t.lastScopeBuf.String())
		if s != "" {
			panic("find tpl data after <? } ?>")
		}
	}
	output := t.outBuf.Bytes()
	addPkgList := []string{}
	if t.hasBytesPackage {
		addPkgList = append(addPkgList, "bytes")
	}
	if t.hasKmgXssPackage {
		addPkgList = append(addPkgList, "github.com/bronze1man/kmg/kmgXss")
	}
	output = addImport(output, addPkgList)
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
		if t.isHtml {
			s = strings.Trim(s, "\n")
		}
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
		s := strings.TrimSpace(t.lastScopeBuf.String())
		if t.isHtml {
			// raw
			if strings.HasPrefix(s, "raw(") && strings.HasSuffix(s, ")") {
				t.outBuf.WriteString(s[4 : len(s)-1])
			} else if t.urlvStatus == urlvStatusValue || t.urlvStatus == urlvStatusEqual {
				t.hasKmgXssPackage = true
				t.outBuf.WriteString("kmgXss.Urlv(")
				t.outBuf.WriteString(s)
				t.outBuf.WriteString(")")
			} else if t.isInScript {
				t.hasKmgXssPackage = true
				t.outBuf.WriteString("kmgXss.Jsonv(")
				t.outBuf.WriteString(s)
				t.outBuf.WriteString(")")
			} else {
				t.hasKmgXssPackage = true
				t.outBuf.WriteString("kmgXss.H(")
				t.outBuf.WriteString(s)
				t.outBuf.WriteString(")")
			}
		} else {
			t.outBuf.WriteString(s)
		}
		t.outBuf.WriteString(")\n")
	}
	t.lastScopeBuf.Reset()
}

func (t *transformer) isMatchString(token string) bool {
	return bytes.HasPrefix(t.in[t.pos:], []byte(token))
}

func isAlphanum(b byte) bool {
	return (b >= '0' && b <= '9') || (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}
