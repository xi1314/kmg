package kmgLog

import (
	"fmt"
	"os"
)

func StdoutLogWriter(r LogRow){
	b, err := r.Marshal()
	if err != nil {
		fmt.Fprintln(os.Stderr, "[StdoutLogger] Marshal fail", err)
		return
	}
	_, err = fmt.Println(string(b))
	if err != nil {
		fmt.Fprintln(os.Stderr, "[StdoutLogger] printf fail", err)
		return
	}
	return
}
