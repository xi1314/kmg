package kmgRedis

import (
	"errors"
)

var ErrKeyExist = errors.New("key exist")
var ErrKeyNotExist = errors.New("key not exist")

var ErrListWrongType = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value,need a list type")
var ErrSortedSetWrongType = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value,need a sorted set type")
var ErrStringWrongType = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value,need a string type")

var ErrRenameSameName = errors.New("ERR source and destination objects are the same")
var ErrValueNotIntFormatOrOutOfRange = errors.New("ERR value is not an integer or out of range")
var ErrValueNotFloatFormat = errors.New("ERR value is not a valid float")
