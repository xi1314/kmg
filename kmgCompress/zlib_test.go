package kmgCompress

import (
	"fmt"
	"testing"
)

func TestZlib(ot *testing.T) {
	fmt.Println("zlib")
	tsetCompressor(ZlibMustCompress, ZlibMustUnCompress)
}
