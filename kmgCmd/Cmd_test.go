package kmgCmd

import (
	"github.com/bronze1man/kmgTest"
	"testing"
)

func TestExist(t *testing.T) {
	kmgTest.Ok(!Exist("absadfklsjdfa"))
	kmgTest.Ok(Exist("ls"))
	kmgTest.Ok(Exist("top"))
}
