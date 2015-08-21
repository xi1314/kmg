package kmgRpcJava

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/kmgRpc/testPackage"
	"github.com/bronze1man/kmg/kmgTest"
	"os"
	"path/filepath"
	"testing"
)

func TestJava(ot *testing.T) {
	closer := testPackage.ListenAndServe_Demo(":34895", &testPackage.Demo{}, kmgCrypto.Get32PskFromString("abc psk"))
	defer closer()
	os.Chdir(filepath.Join("java", "src"))
	MustGenerateCode(&GenerateRequest{
		ObjectPkgPath:   "github.com/bronze1man/kmg/kmgRpc/testPackage",
		ObjectName:      "Demo",
		ObjectIsPointer: true,
		OutFilePath:     "testPackage/RpcDemo.java",
		OutPackageName:  "testPackage",
		OutClassName:    "RpcDemo",
	})
	kmgCmd.MustRun("javac -sourcepath . ./testPackage/Main.java")
	out := kmgCmd.MustRunAndReturnOutput("java -classpath . testPackage.Main")
	kmgTest.Equal(out, []byte("Success\n"))
}
