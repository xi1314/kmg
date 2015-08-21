package kmgCompress

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
)

// 正常使用应该不会报错.
func ZlibMustCompress(inb []byte) (outb []byte) {
	buf := &bytes.Buffer{}
	w := zlib.NewWriter(buf)
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

func ZlibMustUnCompress(inb []byte) (outb []byte) {
	buf := bytes.NewBuffer(inb)
	reader, err := zlib.NewReader(buf)
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

func ZlibUnCompress(inb []byte) (outb []byte, err error) {
	buf := bytes.NewBuffer(inb)
	reader, err := zlib.NewReader(buf)
	if err != nil {
		return nil, err
	}
	outb, err = ioutil.ReadAll(reader)
	if err != nil {
		reader.Close()
		return nil, err
	}
	err = reader.Close()
	if err != nil {
		panic(err)
	}
	return outb, nil
}
