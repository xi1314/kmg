package kmgRpc

import (
	"testing"

	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgTest"

	"github.com/bronze1man/kmg/kmgView/kmgGoTpl"
)

func TestMustGenerateCode(t *testing.T) {
	kmgGoTpl.MustBuildTplInDirWithCache("src/github.com/bronze1man/kmg/kmgRpc") // 模板变化需要运行两次,才能看到结果.
	kmgFile.MustDelete("testPackage/generated.go")
	MustGenerateCode(&GenerateRequest{
		ObjectPkgPath:        "github.com/bronze1man/kmg/kmgRpc/testPackage",
		ObjectName:           "Demo",
		ObjectIsPointer:      true,
		OutFilePath:          "testPackage/generated.go",
		OutPackageImportPath: "github.com/bronze1man/kmg/kmgRpc/testPackage",
	})
	kmgCmd.CmdString("kmg go test").SetDir("testPackage").Run()
}

func TestReflectToTplConfig(t *testing.T) {
	conf := reflectToTplConfig(
		&GenerateRequest{
			ObjectPkgPath:        "github.com/bronze1man/kmg/kmgRpc/testPackage",
			ObjectName:           "Demo",
			ObjectIsPointer:      true,
			OutFilePath:          "testPackage/generated.go",
			OutPackageImportPath: "github.com/bronze1man/kmg/kmgRpc/testPackage",
		},
	)
	kmgTest.Equal(len(conf.ApiList), 7)
	for i, name := range []string{
		"DemoFunc2", "DemoFunc3", "DemoFunc4", "DemoFunc5", "DemoFunc7", "DemoFunc8", "PostScoreInt",
	} {
		kmgTest.Equal(conf.ApiList[i].Name, name)
	}
}

func TestTplGenerateCode(t *testing.T) {
	out := tplGenerateCode(&tplConfig{
		OutPackageName: "tplTestPackage",
		ObjectName:     "Demo",
		ObjectTypeStr:  "*Demo",
		ApiList: []Api{
			{
				Name: "PostScoreInt",
				InArgsList: []ArgumentNameTypePair{
					{
						Name:          "LbId",
						ObjectTypeStr: "string",
					},
					{
						Name:          "Score",
						ObjectTypeStr: "int",
					},
				},
				OutArgsList: []ArgumentNameTypePair{
					{
						Name:          "Info",
						ObjectTypeStr: "string",
					},
					{
						Name:          "Err",
						ObjectTypeStr: "error",
					},
				},
			},
		},
		ImportPathMap: map[string]bool{
			"encoding/json": true,
			"errors":        true,
			"fmt":           true,
			"github.com/bronze1man/kmg/kmgCrypto":      true,
			"github.com/bronze1man/kmg/kmgLog":         true,
			"github.com/bronze1man/kmg/kmgNet/kmgHttp": true,
			"net/http": true,
			"bytes":    true,
		},
	})
	kmgFile.MustDeleteFile("tplTestPackage/generated.go")
	kmgFile.MustWriteFileWithMkdir("tplTestPackage/generated.go", []byte(out))
	kmgCmd.CmdString("kmg go test").SetDir("tplTestPackage").Run()
}
