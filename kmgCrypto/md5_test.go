package kmgCrypto

import (
	"testing"

	"github.com/bronze1man/kmg/kmgTest"
)

/*
func TestBytes(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	expect_bytes, _ := hex.DecodeString("d41d8cd98f00b204e9800998ecf8427e")
	t.Equal(Bytes([]byte("")), expect_bytes)
	t.Equal(Bytes([]byte("")), expect_bytes)
	expect_bytes, _ = hex.DecodeString("9e107d9d372bb6826bd81d3542a419d6")
	t.Equal(Bytes([]byte("The quick brown fox jumps over the lazy dog")),
		expect_bytes)

}
*/
func TestMd5Hex(ot *testing.T) {
	kmgTest.Equal(Md5Hex([]byte("")), "d41d8cd98f00b204e9800998ecf8427e")
	kmgTest.Equal(Md5Hex([]byte("The quick brown fox jumps over the lazy dog")),
		"9e107d9d372bb6826bd81d3542a419d6")
}

/*
func TestHexFromString(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	t.Equal(HexFromString(strings.Repeat("1", 10000)), "b223cca8b360eae4e49568512e2de29f")
}
*/
