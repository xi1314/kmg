package kmgHex

import (
	"encoding/hex"
	"strings"
)

func EncodeToUpperString(b []byte) string {
	return strings.ToUpper(hex.EncodeToString(b))
}
