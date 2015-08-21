package kmgCompress

import (
	"bytes"
	"compress/lzw"
	"fmt"
	"io/ioutil"
	"testing"
)

func LzwMustCompress(inb []byte) (outb []byte) {
	buf := &bytes.Buffer{}
	w := lzw.NewWriter(buf, lzw.LSB, 8)
	_, err := w.Write(inb)
	if err != nil {
		w.Close()
		panic(err)
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func LzwMustUnCompress(inb []byte) (outb []byte) {
	buf := bytes.NewBuffer(inb)
	reader := lzw.NewReader(buf, lzw.LSB, 8)
	outb, err := ioutil.ReadAll(reader)
	if err != nil {
		reader.Close()
		panic(err)
	}
	err = reader.Close()
	if err != nil {
		panic(err)
	}
	return outb
}

func TestLzw(ot *testing.T) {
	fmt.Println("lzw")
	tsetCompressor(LzwMustCompress, LzwMustUnCompress)
}
