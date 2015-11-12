package SubCommand

import (
	"flag"
	"fmt"
	"github.com/bronze1man/kmg/kmgRand"
)

func NewPassword() {
	len := 0
	num := 0
	var typ string
	flag.StringVar(&typ, "type", "AlphaNum", "password type (AlphaNum,Num)")
	flag.IntVar(&len, "len", 10, "len of password")
	flag.IntVar(&num, "num", 1, "num of password")
	flag.Parse()
	var f func(length int) string
	switch typ {
	case "AlphaNum":
		f = kmgRand.MustCryptoRandToReadableAlphaNum
	case "Num":
		f = kmgRand.MustCryptoRandToNum
	default:
		fmt.Println("Unknow password type")
		flag.Usage()
		return
	}
	for i := 0; i < num; i++ {
		fmt.Println(f(len))
	}
}
