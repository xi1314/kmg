package kmgCrypto_test

import (
	"testing"
	//"fmt"
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"

	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/kmgNet"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgTest"
)

func TestMustTlsTestListen(ot *testing.T) {
	listener := kmgCrypto.MustTlsTestListen(":0")
	defer listener.Close()
	listenAddr := kmgNet.MustGetLocalAddrFromListener(listener)
	waitAcceptChan := make(chan int)
	go func() {
		defer func() { waitAcceptChan <- 1 }()
		conn, err := listener.Accept()
		kmgTest.Equal(err, nil)
		defer conn.Close()
		//此处开始检测两条连接是否连在了一起
		result, err := ioutil.ReadAll(conn)
		kmgTest.Equal(err, nil)
		kmgTest.Ok(bytes.Equal(result, []byte(`hello world`)))
	}()
	conn, err := tls.Dial("tcp", listenAddr, &tls.Config{
		InsecureSkipVerify: true,
	})
	kmgTest.Equal(err, nil)
	defer conn.Close()
	_, err = conn.Write([]byte(`hello world`))
	kmgTest.Equal(err, nil)
	conn.Close()
	<-waitAcceptChan
}

func TestMustTlsTestListenHttps(ot *testing.T) {
	listener := kmgCrypto.MustTlsTestListen(":0")
	defer listener.Close()
	go http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("hello world"))
	}))
	c := kmgHttp.NewHttpsCertNotCheckClient()
	resp, err := c.Get("https://" + kmgNet.MustGetLocalAddrFromListener(listener))
	kmgTest.Equal(err, nil)
	kmgTest.Equal(resp.StatusCode, 200)
}
