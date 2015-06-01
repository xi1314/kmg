package kmgExchangeRate

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestMustGetExchangeRate(ot *testing.T) {
	rate := MustGetExchangeRate("JPY", "CNY")
	kmgTest.Ok(rate > 0)
}
