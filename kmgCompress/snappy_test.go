package kmgCompress

import (
	"fmt"
	"github.com/golang/snappy"
	"testing"
)

func SnappyMustCompress(inb []byte) (outb []byte) {
	return snappy.Encode(nil, inb)
}

func SnappyMustUnCompress(inb []byte) (outb []byte) {
	outb, err := snappy.Decode(nil, inb)
	if err != nil {
		panic(err)
	}
	return outb
}

func TestSnappy(ot *testing.T) {
	fmt.Println("Snappy")
	tsetCompressor(SnappyMustCompress, SnappyMustUnCompress)
}
