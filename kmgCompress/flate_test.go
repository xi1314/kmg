package kmgCompress

import (
	"testing"

	"fmt"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgNet"
	"github.com/bronze1man/kmg/kmgRand"
	"github.com/bronze1man/kmg/kmgTest"
	"time"
)

func TestFlate(ot *testing.T) {
	fmt.Println("flate")
	tsetCompressor(FlateMustCompress, FlateMustUnCompress)
}

func tsetCompressor(compressor func(inb []byte) (outb []byte), decompressor func(inb []byte) (outb []byte)) {
	origin := []byte("123")
	ob := compressor([]byte("123"))
	output := decompressor(ob)
	kmgTest.Equal(origin, output)

	//for i:=1;i<20;i++ {
	//	fmt.Println(kmgHex.UpperEncodeBytesToString(compressor(bytes.Repeat([]byte{'1'},i))))
	//}
	//for i:=1;i<100;i++ {
	//	fmt.Println(kmgHex.UpperEncodeBytesToString(compressor(kmgRand.MustCryptoRandBytes(i))))
	//}
	for _, i := range []int{1, 10, 100, 1000, 1e4, 1e5} {
		ob = compressor(kmgRand.MustCryptoRandBytes(i))
		fmt.Println(i, len(ob), len(ob)-i)
	}
	for _, path := range []string{
		"/Users/bronze1man/tmp/vanke.sql",
		"/bin/kmg",
	} {
		c := kmgFile.MustReadFile(path)
		t := time.Now()
		ob = compressor(c)
		dur := time.Since(t)
		fmt.Println(path, kmgNet.SizeString(int64(len(c))), kmgNet.SizeString(int64(len(ob))),
			float64(len(ob))/float64(len(c)), dur)
	}

	fmt.Println()
}
