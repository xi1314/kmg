package kmgDebug

import (
	"github.com/bronze1man/kmgTest"
	"testing"
)

func TestSprintln(ot *testing.T) {
	kmgTest.Equal(Sprintln([]byte{0, 1}), "[kmgDebug.Println] [0 1]\n")
	kmgTest.Equal(Sprintln([]byte{}), "[kmgDebug.Println] []\n")

	kmgTest.Equal(Sprintln([]byte(nil)), "[kmgDebug.Println] nil\n")
}
