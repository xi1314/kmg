package JsTpl

import (
	"bytes"
	"path/filepath"

	"github.com/bronze1man/kmg/encoding/kmgJson"
	"github.com/bronze1man/kmg/kmgFile"
)

func MustBuildTplOneFile(in []byte) (out []byte) {
	//把所有的`xxx` 转成 "xxx\xxx" 之类的,严格保留里面的所有字符串,暂时没有任何办法可以打出 `
	outbuf := &bytes.Buffer{}
	thisHereDocBuf := &bytes.Buffer{}
	isInHereDoc := false
	for i := range in {
		//进入heredoc
		if !isInHereDoc && in[i] == '`' {
			isInHereDoc = true
			continue
		}
		//出heredoc
		if isInHereDoc && in[i] == '`' {
			isInHereDoc = false
			outbuf.Write(kmgJson.MustMarshal(thisHereDocBuf.String()))
			thisHereDocBuf.Reset()
			continue
		}
		//不是heredoc的部分
		if !isInHereDoc {
			outbuf.WriteByte(in[i])
			continue
		}
		//是heredoc的部分
		if isInHereDoc {
			thisHereDocBuf.WriteByte(in[i])
			continue
		}
	}
	if isInHereDoc {
		panic("end with heredoc")
	}
	return outbuf.Bytes()
}

func MustBuildTplInDir(path string) {
	pathList, err := kmgFile.GetAllFiles(path)
	if err != nil {
		panic(err)
	}
	for _, val := range pathList {
		if filepath.Ext(val) != ".jst" {
			continue
		}
		out := MustBuildTplOneFile(kmgFile.MustReadFile(val))
		outFilePath := kmgFile.PathTrimExt(val) + ".js"
		kmgFile.MustWriteFile(outFilePath, out)
	}
}
