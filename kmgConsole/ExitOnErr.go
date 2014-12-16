package kmgConsole

import (
	"fmt"
	"os"
)

func ExitOnErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func ExitOnStderr(err string) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
