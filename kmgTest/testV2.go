package kmgTest

import (
	"fmt"
)

func Ok(expectTrue bool) {
	if !expectTrue {
		panic("ok fail")
	}
}

func Equal(get interface{}, expect interface{}) {
	if isEqual(expect, get) {
		return
	}
	msg := fmt.Sprintf("\tget1: %s\n\texpect2: %s", valueDetail(get), valueDetail(expect))
	panic(msg)
}

func valueDetail(value interface{}) string {
	stringer, ok := value.(tStringer)
	if ok {
		return fmt.Sprintf("%s (%T) %#v", stringer.String(), value, value)
	} else {
		return fmt.Sprintf("%#v (%T)", value, value)
	}
}

type tStringer interface {
	String() string
}
