package kmgPing

import (
	"github.com/tatsushid/go-fastping"
	"net"
	"time"
	"github.com/bronze1man/kmg/kmgLog"
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

// 丢包返回1,不丢包返回0
func (e Echo) GetLostRateFloat()float64{
	if e.Status==EchoStatusSuccess{
		return 0
	}else{
		return 1
	}
}

var MaxRtt time.Duration = time.Duration(1e9)

func Ping(address string) Echo {
	p := fastping.NewPinger()
	p.MaxRTT = MaxRtt
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
	if echo.Rtt==0 && echo.Status==EchoStatusSuccess{
		kmgLog.Log("error","[kmgPing.Ping] echo.Rtt==0 && echo.Status==EchoStatusSuccess",address)
	}
	return echo
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
