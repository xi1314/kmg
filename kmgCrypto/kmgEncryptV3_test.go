package kmgCrypto

import (
	"testing"
)

func TestEncryptV3(ot *testing.T) {
	EncryptTester(EncryptV3, DecryptV3, 20)
}
