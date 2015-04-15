package kmgTest

import (
	"testing"
)

func TestEqual(ot *testing.T) {
	Equal(true, true)
	Equal([]byte{1}, []byte{1})
	AssertPanic(func() {
		Equal(true, false)
	})
	Ok(true)
	AssertPanic(func() {
		Ok(false)
	})
	AssertPanic(func() {
		Equal(int64(1), int(1))
	})
}
func TestIsEqual(ot *testing.T) {
	if isEqual(map[string]interface{}{"a": 1}, map[string]interface{}{"a": 1}) == false {
		panic("fail")
	}
}

func TestV2(ot *testing.T) {
	msg := AssertPanic(func() {
		Ok(false)
	})
	Equal(msg, "ok fail")
	Ok(true)

	//equal byte
	Equal([]byte{1, 2}, []byte{1, 2})

	AssertPanic(func() {
		Equal([]byte{1, 2}, []byte{1, 3})
	})

	Equal(nil, (*testing.T)(nil))

	//assert panic
	flag := 1
	msg = AssertPanic(func() {
		AssertPanic(func() {
			flag = 2 //check this function has already run.
		})
		//panic should not pass to this line,so this test can verify that AssertPanic is working
		flag = 3
	})
	Equal(flag, 2)
	Equal(msg, "should panic")
}
