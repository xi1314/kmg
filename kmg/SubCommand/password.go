package SubCommand
import (
	"github.com/bronze1man/kmg/kmgRand"
	"fmt"
)

func NewPassowrd(){
	fmt.Println(kmgRand.MustCryptoRandToReadableAlphaNum(10))
}