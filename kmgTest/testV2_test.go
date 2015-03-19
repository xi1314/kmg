package kmgTest

import "testing"

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
	flag := false
	msg = AssertPanic(func() {
		AssertPanic(func() {
		})
		//panic should not pass to this line,so this test can verify that AssertPanic is working
		flag = true
	})
	Equal(flag, false)
	Equal(msg, "should panic")
}
