package kmgViewResource
import (
	"testing"
	"github.com/bronze1man/kmg/kmgTest"
)

func TestParseImportPath(ot *testing.T){
	importPathList:=parseImportPath("<nil>",[]byte(`
/*
	import (
		"github.com/bronze1man/kmg/kmgView/webResource/bootstrap"
		"boostrap"
	)
*/
`))
	kmgTest.Equal(importPathList,[]string{
		"github.com/bronze1man/kmg/kmgView/webResource/bootstrap",
		"boostrap",
	})

	importPathList=parseImportPath("<nil>",[]byte(`
/*
*/
`))
	kmgTest.Equal(importPathList,nil)

	importPathList=parseImportPath("<nil>",[]byte(`
/*
*/
/*
	import (
		"github.com/bronze1man/kmg/kmgView/webResource/bootstrap"
		"boostrap"
	)
*/
`))
	kmgTest.Equal(importPathList,nil)
}