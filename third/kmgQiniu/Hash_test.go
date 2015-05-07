package kmgQiniu

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestHash(ot *testing.T) {
	etag, err := ComputeHashFromFile("testFile/1.txt")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(etag, "FqmZPjZHBoFquj4lcXhQwmyc0Nid")

	etag = ComputeHashFromBytes([]byte("abc"))
	kmgTest.Equal(etag, "FqmZPjZHBoFquj4lcXhQwmyc0Nid")
}
