package kmgCrypto

import (
	"testing"
	//"fmt"
	"bytes"
	"crypto/tls"
	"github.com/bronze1man/kmg/kmgNet"
	"github.com/bronze1man/kmg/kmgTest"
	"github.com/bronze1man/kmg/net/kmgHttp"
	"io/ioutil"
	"net/http"
)

func TestMustTlsTestListen(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	listener := MustTlsTestListen(":0")
	defer listener.Close()
	listenAddr := kmgNet.MustGetLocalAddrFromListener(listener)
	waitAcceptChan := make(chan int)
	go func() {
		defer func() { waitAcceptChan <- 1 }()
		conn, err := listener.Accept()
		t.Equal(err, nil)
		defer conn.Close()
		//此处开始检测两条连接是否连在了一起
		result, err := ioutil.ReadAll(conn)
		t.Equal(err, nil)
		t.Ok(bytes.Equal(result, []byte(`hello world`)))
	}()
	conn, err := tls.Dial("tcp", listenAddr, &tls.Config{
		InsecureSkipVerify: true,
	})
	t.Equal(err, nil)
	defer conn.Close()
	_, err = conn.Write([]byte(`hello world`))
	t.Equal(err, nil)
	conn.Close()
	<-waitAcceptChan
}

func TestMustTlsTestListenHttps(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	listener := MustTlsTestListen(":0")
	defer listener.Close()
	go http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("hello world"))
	}))
	c := kmgHttp.NewHttpsCertNotCheckClient()
	resp, err := c.Get("https://" + kmgNet.MustGetLocalAddrFromListener(listener))
	t.Equal(err, nil)
	t.Equal(resp.StatusCode, 200)
}
