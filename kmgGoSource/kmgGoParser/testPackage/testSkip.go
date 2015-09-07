package testPackage

import (
	. "bytes"
)

var UnreadRuneErrorTests = []struct {
	name string
	f    func(*Reader)
}{
	{"Read", func(r *Reader) { r.Read([]byte{0}) }},
	{"ReadByte", func(r *Reader) { r.ReadByte() }},
	{"UnreadRune", func(r *Reader) { r.UnreadRune() }},
	{"Seek", func(r *Reader) { r.Seek(0, 1) }},
	{"WriteTo", func(r *Reader) { r.WriteTo(&Buffer{}) }},
}

// Do the 64-bit functions panic?  If so, don't bother testing.
var test64err = func() (err interface{}) {
	defer func() {
		err = recover()
	}()
	//var x int64
	return nil
}()

var test1 = "string"

type serverType func()

// dial constructs a new test server and returns a *ClientConn.
func dial() {

}
