package demo

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestBoth(ot *testing.T) {
	go ListenAndServe_Demo(":34895", &Demo{})
	client := NewClient("127.0.0.1:34895")
	err := client.PostScoreInt("LbId", 1)
	kmgTest.Equal(err, nil)
	kmgTest.Equal(lastPostLbId, "LdId")
}
