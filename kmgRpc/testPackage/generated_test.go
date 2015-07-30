package testPackage

import (
	"testing"
	"time"

	"github.com/bronze1man/kmg/kmgTest"
)

func TestGenerated(ot *testing.T) {
	closer := ListenAndServe_Demo(":34895", &Demo{})
	defer closer()
	client := NewClient_Demo("http://127.0.0.1:34895/f")
	info, err := client.PostScoreInt("LbId", 1)
	kmgTest.Equal(err, nil)
	kmgTest.Equal(info, "LbId")

	info, err = client.PostScoreInt("LbId", 2)
	kmgTest.Equal(err.Error(), "Score!=1")
	kmgTest.Equal(info, "")
}

func BenchmarkGenerated(ot *testing.B) {
	closer := ListenAndServe_Demo(":34896", &Demo{})
	defer func() {
		closer()
		time.Sleep(10 * time.Millisecond)
	}()
	client := NewClient_Demo("http://127.0.0.1:34896/f")
	ot.ResetTimer()
	for i := 0; i < ot.N; i++ {
		info, err := client.PostScoreInt("LbId", 1)
		kmgTest.Equal(err, nil)
		kmgTest.Equal(info, "LbId")
	}
	ot.StopTimer()
	// 结果 770 qps BenchmarkGenerated	    1000	   1298511 ns/op
}
