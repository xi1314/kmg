package kmgDebug_test

import (
	"testing"

	"github.com/bronze1man/kmg/kmgDebug"
	"github.com/bronze1man/kmg/kmgTest"
)

func TestSprintln(ot *testing.T) {
	kmgTest.Equal(kmgDebug.Sprintln([]byte{0, 1}), "[kmgDebug.Println] [0 1]\n")
	kmgTest.Equal(kmgDebug.Sprintln([]byte{}), "[kmgDebug.Println] []\n")

	kmgTest.Equal(kmgDebug.Sprintln([]byte(nil)), "[kmgDebug.Println] nil\n")
}
