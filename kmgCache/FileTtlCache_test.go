package kmgCache

import (
	"fmt"
	"testing"
	"time"

	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgTest"
)

func TestFileTtlCache(t *testing.T) {
	kmgFile.MustDelete(getFileTtlCachePath("test_file_ttl_cache"))
	_, err := FileTtlCache("test_file_ttl_cache", func() (b []byte, ttl time.Duration, err error) {
		return []byte("1"), time.Millisecond, fmt.Errorf("error")
	})
	kmgTest.Equal(err.Error(), "error")
	b, err := FileTtlCache("test_file_ttl_cache", func() (b []byte, ttl time.Duration, err error) {
		return []byte("1"), time.Millisecond, nil
	})
	kmgTest.Equal(b, []byte("1"))
	kmgTest.Equal(err, nil)

	b, err = FileTtlCache("test_file_ttl_cache", func() (b []byte, ttl time.Duration, err error) {
		return []byte("2"), time.Millisecond, nil
	})
	kmgTest.Equal(b, []byte("1"))
	kmgTest.Equal(err, nil)

	time.Sleep(time.Millisecond)
	b, err = FileTtlCache("test_file_ttl_cache", func() (b []byte, ttl time.Duration, err error) {
		return []byte("2"), time.Millisecond, nil
	})
	kmgTest.Equal(b, []byte("2"))
	kmgTest.Equal(err, nil)
}
