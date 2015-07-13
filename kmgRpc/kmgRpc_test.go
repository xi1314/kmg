package kmgRpc

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgRpc/testPackage"
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestMustGenerateCode(t *testing.T) {
	kmgFile.MustDeleteFile("testPackage/generated.go")
	MustGenerateCode(GenerateRequest{
		Object:               &testPackage.Demo{},
		ObjectName:           "Demo",
		OutFilePath:          "testPackage/generated.go",
		OutPackageImportPath: "github.com/bronze1man/kmg/kmgRpc/testPackage",
	})
	kmgCmd.CmdString("kmg go test").SetDir("testPackage").Run()
}

func TestReflectToTplConfig(t *testing.T) {
	conf := reflectToTplConfig(
		GenerateRequest{
			Object:               &testPackage.Demo{},
			ObjectName:           "Demo",
			OutFilePath:          "testPackage/generated.go",
			OutPackageImportPath: "github.com/bronze1man/kmg/kmgRpc/testPackage",
		},
	)
	kmgTest.Equal(len(conf.ApiList), 6)
	for i, name := range []string{
		"DemoFunc2", "DemoFunc3", "DemoFunc4", "DemoFunc5", "DemoFunc7", "PostScoreInt",
	} {
		kmgTest.Equal(conf.ApiList[i].Name, name)
	}
}

func TestTplGenerateCode(t *testing.T) {
	out := tplGenerateCode(tplConfig{
		OutPackageName: "testPackage",
		OutKeyByteList: "1,2",
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
						Name:          "err",
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
	kmgFile.MustDeleteFile("testPackage/generated.go")
	kmgFile.MustWriteFileWithMkdir("testPackage/generated.go", []byte(out))
	kmgCmd.CmdString("kmg go test").SetDir("testPackage").Run()
}
