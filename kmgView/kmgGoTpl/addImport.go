package kmgGoTpl

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgStrings"
)

/*
type importStatus int

const (
	importStatusNot          importStatus = 0
	importStatusSpace1       importStatus = 1
	importStatusImportToken  importStatus = 2
	importStatusParentheses1 importStatus = 3
	importStatusImportList   importStatus = 4
	importStatusParentheses2 importStatus = 5
)
*/

// 添加bytes的import项.
// 这个是近似实现,golang的语法树确实有点复杂.
func addImport(in []byte, pkgList []string) (out []byte) {
	var isLastImportToken bool
	var isInImportParentheses bool
	var lastImportParenthesesPos int
	var hasFoundImport bool
	var hadAddImport bool
	outBuf := &bytes.Buffer{}
	for pos := 0; pos < len(in); pos++ {
		if !hadAddImport {
			if bytes.HasPrefix(in[pos:], []byte("import")) {
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
					for _, pkg := range pkgList {
						if !kmgStrings.IsInSlice(readedImportPathList, "\""+pkg+"\"") {
							outBuf.WriteString("\"" + pkg + "\"\n")
						}
					}
					hadAddImport = true
					isInImportParentheses = false
				}
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
					outBuf.WriteString("\nimport (")
					for _, pkg := range pkgList {
						outBuf.WriteString("\"" + pkg + "\"\n")
					}
					outBuf.WriteString(")")
					isLastPackageToken = false
				}
			}
			outBuf.WriteByte(in[pos])
		}
	}
	return outBuf.Bytes()
}
