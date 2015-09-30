package kmgCompress

import (
	"fmt"
	"testing"
)

func TestGzip(ot *testing.T) {
	fmt.Println("gzip")
	tsetCompressor(GzipMustCompress, GzipMustUnCompress)
}
