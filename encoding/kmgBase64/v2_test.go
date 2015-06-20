package kmgBase64

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestV2(ot *testing.T) {
	for _, testCase := range []struct {
		in  string
		out string
	}{
		{"", ""},
		{"A", "QQ"},
		{"\x00", "AA"},
		{"\x00\x00", "AAA"},
		{"\x00\x00\x00", "AAAA"},
		{"\xff\xff\xff", "____"},
		{"\x14\xfb\x9c\x03\xd9\x7e", "FPucA9l-"},
		{"\x14\xfb\x9c\x03\xd9", "FPucA9k"},
		{"\x14\xfb\x9c\x03", "FPucAw"},
		// RFC 4648 examples
		{"", ""},
		{"f", "Zg"},
		{"fo", "Zm8"},
		{"foo", "Zm9v"},
		{"foob", "Zm9vYg"},
		{"fooba", "Zm9vYmE"},
		{"foobar", "Zm9vYmFy"},

		// Wikipedia examples
		{"sure.", "c3VyZS4"},
		{"sure", "c3VyZQ"},
		{"sur", "c3Vy"},
		{"su", "c3U"},
		{"leasure.", "bGVhc3VyZS4"},
		{"easure.", "ZWFzdXJlLg"},
		{"asure.", "YXN1cmUu"},
		{"sure.", "c3VyZS4"},
		{"Twas brillig, and the slithy toves", "VHdhcyBicmlsbGlnLCBhbmQgdGhlIHNsaXRoeSB0b3Zlcw"},
	} {
		kmgTest.Equal(EncodeByteToStringV2([]byte(testCase.in)), testCase.out)
		thisIn, err := DecodeStringToByteV2(testCase.out)
		kmgTest.Equal(err, nil)
		kmgTest.Equal(thisIn, []byte(testCase.in))
	}
}

// 性能还好,暂时不用管性能 52.21 MB/s 官方 107.21 MB/s
func BenchmarkEncodeToStringV2(b *testing.B) {
	data := make([]byte, 8192)
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		EncodeByteToStringV2(data)
	}
}

// 性能还好,暂时不用管性能 51.37 MB/s 官方 54.87 MB/s
func BenchmarkDecodeStringV2(b *testing.B) {
	data := EncodeByteToStringV2(make([]byte, 8192))
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		DecodeStringToByteV2(data)
	}
}
