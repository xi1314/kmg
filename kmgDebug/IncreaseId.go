package kmgDebug

import (
	"strconv"
	"sync/atomic"
)

var intId uint64

func NextIntIdString() string {
	var idInt = atomic.AddUint64(&intId, 1)
	return strconv.FormatUint(idInt, 10)
}
