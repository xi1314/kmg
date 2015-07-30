package kmgExchangeRate

import (
	"testing"

	"github.com/bronze1man/kmg/kmgTest"
)

func TestMustGetExchangeRate(ot *testing.T) {
	rate := MustGetExchangeRate("JPY", "CNY")
	kmgTest.Ok(rate > 0)
}
