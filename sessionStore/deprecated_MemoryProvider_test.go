package sessionStore

func (t *Tester) TestManagerProvider() {
	p := NewMemoryProvider()

	err := p.Set("g881r0H-B4fIGF8ktUWeUg==", []byte("2"))
	t.Equal(err, nil)

	v, ok, err := p.Get("g881r0H-B4fIGF8ktUWeUg==")
	t.Equal(err, nil)
	t.Equal(ok, true)
	t.Equal(v, []byte("2"))

	err = p.Delete("g881r0H-B4fIGF8ktUWeUg==")
	t.Equal(err, nil)

	v, ok, err = p.Get("g881r0H-B4fIGF8ktUWeUg==")
	t.Equal(err, nil)
	t.Equal(ok, false)

	err = p.Set("g881r0H-B4fIGF8ktUWeUg==", []byte("3"))
	t.Equal(err, nil)
	v, ok, err = p.Get("g881r0H-B4fIGF8ktUWeUg==")
	t.Equal(err, nil)
	t.Equal(ok, true)
	t.Equal(v, []byte("3"))
}
