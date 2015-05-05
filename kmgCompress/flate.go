package kmgCompress

import (
	"bytes"
	"compress/flate"
	"io/ioutil"
)

//flate压缩,panic+buffer实现
func FlateMustCompress(inb []byte) (outb []byte) {
	buf := &bytes.Buffer{}
	w, err := flate.NewWriter(buf, -1)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(inb)
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

func FlateMustUnCompress(inb []byte) (outb []byte) {
	buf := bytes.NewBuffer(inb)
	reader := flate.NewReader(buf)
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

func FlateUnCompress(inb []byte) (outb []byte, err error) {
	buf := bytes.NewBuffer(inb)
	reader := flate.NewReader(buf)
	outb, err = ioutil.ReadAll(reader)
	if err != nil {
		reader.Close()
		return nil, err
	}
	err = reader.Close()
	if err != nil {
		return
	}
	return outb, nil
}
