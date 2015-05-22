package testPackage

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestGenerated(ot *testing.T) {
	go ListenAndServe_Demo(":34895", &Demo{})
	client := NewClient_Demo("http://127.0.0.1:34895/f")
	info, err := client.PostScoreInt("LbId", 1)
	kmgTest.Equal(err, nil)
	kmgTest.Equal(info, "LbId")

	info, err = client.PostScoreInt("LbId", 2)
	kmgTest.Equal(err.Error(), "Score!=1")
	kmgTest.Equal(info, "")
}
