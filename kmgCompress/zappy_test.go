package kmgCompress

import (
	"fmt"
	"github.com/cznic/zappy"
	"testing"
)

func ZappyMustCompress(inb []byte) (outb []byte) {
	outb, err := zappy.Encode(nil, inb)
	if err != nil {
		panic(err)
	}
	return outb
}

func ZappyMustUnCompress(inb []byte) (outb []byte) {
	outb, err := zappy.Decode(nil, inb)
	if err != nil {
		panic(err)
	}
	return outb
}

func TestZappy(ot *testing.T) {
	fmt.Println("zappy")
	tsetCompressor(ZappyMustCompress, ZappyMustUnCompress)
}
