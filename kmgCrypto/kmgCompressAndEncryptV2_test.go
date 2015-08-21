package kmgCrypto

import (
	"testing"
)

func TestCompressAndEncryptBytesV2(ot *testing.T) {
	EncryptTester(CompressAndEncryptBytesEncodeV2, CompressAndEncryptBytesDecodeV2, 21)
}
