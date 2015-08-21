package kmgCompress

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"testing"
)

func GzipMustCompress(inb []byte) (outb []byte) {
	buf := &bytes.Buffer{}
	w := gzip.NewWriter(buf)
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

func GzipMustUnCompress(inb []byte) (outb []byte) {
	buf := bytes.NewBuffer(inb)
	reader, err := gzip.NewReader(buf)
	if err != nil {
		panic(err)
	}
	outb, err = ioutil.ReadAll(reader)
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

func TestGzip(ot *testing.T) {
	fmt.Println("gzip")
	tsetCompressor(GzipMustCompress, GzipMustUnCompress)
}
