package SubCommand

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgRand"
)

func NewPassowrd() {
	fmt.Println(kmgRand.MustCryptoRandToReadableAlphaNum(10))
}
