package testPackage

import (
	"testing"
	"time"

	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/kmgTest"
	"github.com/bronze1man/kmg/kmgTime"
)

func TestGenerated(ot *testing.T) {
	psk := kmgCrypto.Get32PskFromString("abc")
	closer := ListenAndServe_Demo(":34895", &Demo{}, psk)
	defer closer()
	client := NewClient_Demo("http://127.0.0.1:34895/f", psk)
	info, err := client.PostScoreInt("LbId", 1)
	kmgTest.Equal(err, nil)
	kmgTest.Equal(info, "LbId")

	info, err = client.PostScoreInt("LbId", 2)
	kmgTest.Equal(info, "")
	kmgTest.Ok(err != nil, err)
	kmgTest.Equal(err.Error(), "Score!=1")

	info, err = client.DemoFunc8(DemoRequest{}, &DemoRequest{}, 1)
	kmgTest.Equal(info, "info1")

	info, err = client.DemoFunc8(DemoRequest{}, &DemoRequest{}, 2)
	kmgTest.Equal(info, "info")

	inT := kmgTime.MustParseAutoInDefault("2001-01-01 01:01:01")
	outT, err := client.DemoTime(inT)
	kmgTest.Equal(err, nil)
	kmgTest.Ok(outT.Equal(kmgTime.MustParseAutoInDefault("2001-01-01T02:01:01.001000001+08:00")), outT)

	ip, err := client.DemoClientIp()
	kmgTest.Equal(err, nil)
	kmgTest.Equal(ip, "127.0.0.1")
}

func BenchmarkGenerated(ot *testing.B) {
	psk := kmgCrypto.Get32PskFromString("abcd")
	closer := ListenAndServe_Demo(":34896", &Demo{}, psk)
	defer func() {
		closer()
		time.Sleep(10 * time.Millisecond)
	}()
	client := NewClient_Demo("http://127.0.0.1:34896/f", psk)
	ot.ResetTimer()
	for i := 0; i < ot.N; i++ {
		info, err := client.PostScoreInt("LbId", 1)
		kmgTest.Equal(err, nil)
		kmgTest.Equal(info, "LbId")
	}
	ot.StopTimer()
	// 结果 770 qps BenchmarkGenerated	    1000	   1298511 ns/op
}
