package kmgViewResource

import (
	"github.com/bronze1man/kmg/kmgGoSource/kmgGoParser"
	"github.com/bronze1man/kmg/kmgGoSource/kmgGoReader"
)

// 一个文件只会读第一个import,其他的都会忽略.如果在第一个/**/ 之前出现了其他任何东西,也都忽略.
// /*
// import(
// "webResource/bootstrap"
// "webResource/echart"
// )
// */
// 语法错误panic,看上去不是import的情况直接返回nil

var importToken = []byte("import")

func parseImportPath(filename string, content []byte) []string {
	r := kmgGoReader.NewReaderWithPosFile(filename, content)
	r.ReadAllSpace()
	if r.IsEof() {
		return nil
	}
	var b byte
	// 读入 /*
	b = r.ReadByte()
	if b != '/' {
		return nil
	}
	if r.IsEof() {
		return nil
	}
	b = r.ReadByte()
	if b != '*' {
		return nil
	}
	r.ReadAllSpace()
	// 读入 import (
	if !r.IsMatchAfter(importToken) {
		return nil
	}
	r.MustReadMatch(importToken)
	r.ReadAllSpace()
	if r.IsEof() {
		return nil
	}
	b = r.ReadByte()
	if b != '(' {
		return nil
	}
	output := []string{}
	for {
		r.ReadAllSpace()
		if r.IsEof() {
			panic(r.GetFileLineInfo() + " unexpect EOF")
		}
		b = r.ReadByte()
		if b == ')' {
			return output
		}
		if b == '"' || b == '`' {
			r.UnreadByte()
			importPath := kmgGoParser.MustReadGoString(r)
			output = append(output, string(importPath))
		}
	}
}
