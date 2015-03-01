package kmgRand

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"sync"
)

//这个东西的首要目标就是快,内容看上去很随机而已 有 200M/s 左右
var FastRandReader io.Reader = &fastRandReader{}

const randBlockSize = 256

type fastRandReader struct {
	stream  cipher.Stream
	buf     [randBlockSize]byte
	lock    sync.Mutex
	hasInit bool
}

func (r *fastRandReader) Read(dst []byte) (n int, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.init()
	remainSize := len(dst)
	for {
		if remainSize >= randBlockSize {
			r.stream.XORKeyStream(dst[remainSize-randBlockSize:remainSize], r.buf[:])
			remainSize -= randBlockSize
			continue
		}
		r.stream.XORKeyStream(dst[0:remainSize], r.buf[:remainSize])
		break
	}
	return len(dst), nil
}

func (r *fastRandReader) init() {
	if r.stream != nil {
		return
	}
	_, err := io.ReadFull(rand.Reader, r.buf[:])
	if err != nil {
		panic(err)
	}
	block, err := aes.NewCipher(r.buf[:32])
	if err != nil {
		panic(err)
	}
	r.stream = cipher.NewCTR(block, r.buf[32:48])
}
