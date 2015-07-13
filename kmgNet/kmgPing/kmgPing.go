package kmgPing

import (
	"github.com/tatsushid/go-fastping"
	"net"
	"time"
)

type EchoStatus string

const (
	EchoStatusSuccess EchoStatus = "Success"
	EchoStatusFail    EchoStatus = "Fail"
)

type Echo struct {
	Rtt     time.Duration
	Status  EchoStatus
	Address string
}

func Ping(address string) Echo {
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", address)
	handleErr(err)
	p.AddIPAddr(ra)
	echo := Echo{
		Address: address,
	}
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		echo.Rtt = rtt
		echo.Status = EchoStatusSuccess
	}
	p.OnIdle = func() {
		if echo.Status == EchoStatusSuccess {
			return
		}
		echo.Status = EchoStatusFail
		echo.Rtt = time.Duration(1e9)
	}
	err = p.Run()
	handleErr(err)
	return echo
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
