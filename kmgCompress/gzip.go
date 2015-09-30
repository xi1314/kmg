package kmgCompress

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

// 此处不会panic
func GzipMustCompress(inb []byte) (outb []byte) {
	// 只会由于无法write,而报错,但是bytes.Buffer不会报错.
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
