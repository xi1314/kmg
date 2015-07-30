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

type addImportStatus int

const (
	addImportStatusNot                addImportStatus = 0
	addImportStatusPackageToken       addImportStatus = 1 //package
	addImportStatusPackageLine        addImportStatus = 2
	addImportStatusImportToken        addImportStatus = 3
	addImportStatusImportParentheses1 addImportStatus = 4
	addImportStatusImportParentheses2 addImportStatus = 6
	addImportStatusImportImpossible   addImportStatus = 7
)

// 添加bytes的import项.
// 这个是近似实现,golang的语法树确实有点复杂.
func addImport(in []byte, pkgList []string) (out []byte) {
	status := addImportStatus(0)
	var hasFoundImport bool
	var lastImportParenthesesPos int
	outBuf := &bytes.Buffer{}
	for pos := 0; pos < len(in); pos++ {
		switch status {
		case addImportStatusNot:
			if bytes.HasPrefix(in[pos:], []byte("package")) {
				status = addImportStatusPackageToken
			}
		case addImportStatusPackageToken:
			if in[pos] == '\n' {
				status = addImportStatusPackageLine
			}
		case addImportStatusPackageLine:
			if bytes.HasPrefix(in[pos:], []byte("import")) {
				status = addImportStatusImportToken
			} else if bytes.HasPrefix(in[pos:], []byte("func ")) {
				status = addImportStatusImportImpossible
			} else if bytes.HasPrefix(in[pos:], []byte("type ")) {
				status = addImportStatusImportImpossible
			} else if bytes.HasPrefix(in[pos:], []byte("var ")) {
				status = addImportStatusImportImpossible
			}
		case addImportStatusImportToken:
			if in[pos] == '(' {
				status = addImportStatusImportParentheses1
				lastImportParenthesesPos = pos
			}
		case addImportStatusImportParentheses1:
			if in[pos] == ')' {
				status = addImportStatusImportParentheses2
				hasFoundImport = true
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
			}
		}
		outBuf.WriteByte(in[pos])
	}
	if !hasFoundImport || status == addImportStatusImportImpossible {
		outBuf.Reset()
		isLastPackageToken := false
		hadAddImport := false
		for pos := 0; pos < len(in); pos++ {
			if !hadAddImport {
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
						hadAddImport = true
					}
				}
			}
			outBuf.WriteByte(in[pos])
		}
	}
	return outBuf.Bytes()
}
