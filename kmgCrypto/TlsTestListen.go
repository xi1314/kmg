package kmgCrypto

import (
	"crypto/tls"
	"net"
	"sync"
	//"github.com/bronze1man/kmg/kmgNet"
)

var testServerPrivateCert = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAo4pofDVEvWur9udLPWY/0J73rdm+O1vIFRLIh9ZOfAX6T15Q
HFU9OhJr2NO+TAj5aGgs7ANkUYpp69z218Gz4JM6XmNvzy0OtFW+G+wRDp55upfj
CCJpoUgM1DW0o2HYGv+0Z3s6Q2bgcleP7lvh3PM+Szed2PxysU7AC6f+M1upeyCT
KJfsR6mC10mCCTDFPxWWONITnrtl6gRpUMaQAIxskCaqBT9MeI1Q7KBhWPiQpCEi
I7RSdIya/gRC9deg7XTJXwYYM4cPgDw9UZXizS+YD/vHhQ02zwyMDz7/d1umwNFU
kaDUPpWD2NQcMxJFZxphP+V9oni6FQcdx3p9WQIDAQABAoIBAE3vp9uJhzi+FzWT
BEg+cir5lw9CtWWUV8WzLd2OQ9jXAHLSf1sBRCxomvy56/ZQUItxGfdfpm81h0Tg
CMLJjO95OlfBC4ev8VH/ipnD5l0RKPHDcG5v9+bkbziyX2N9PX+gXFj0YM++NzhG
glEVeI3TXdx6pL7Dj+TvopXreTj/srOkvy2p3dX5+38w1R9hIAXteMLPE/p3uM0/
2EV6Gd5H1f1y59l4//jkjQf/Ou3cGnzpW6oDlmFLJl9czZe9oL8HK2NFjGek4KkW
dx507A+rLX+UuxWO+Z+l7UcbMuvjCM1AWYwA287npLqL/RC0jDKf1wdJIRU8DCFv
MJnsyAECgYEA0RyGf7RkFSj5lgewhg7ljjeJpDdRz2vT2ki2T9Qa3OCNiRybMylV
2mSecHJIjL9EgwYCA4JCF+3W8PHyvS9LFnDzbmbUiMrJXGkfNbc3KmvLy+dCIulp
VRLiq9TUzCTKoPRRvr1gviNUSAdjKaYB1giUoOBiIIDeQsqNHihR/ecCgYEAyDYE
JrZ9ME0uNUKk4iPQ/YdvoNK/jhZ7D0LwajryQy/7T/+JlbVY+Ymqoa/AMy5/hpfL
giTct5nYOW+27LBsmkBf9WAvmHs46QIz2QFab/MJEQ3z6o5xZYyxHTxxl/xajth9
KzxtZj1TK8zYXEh19/dXgiaYBDoFNXm6YBJTwr8CgYEAg8CBTb0OwfZLKyg1JIIG
SJDdfEYOma3KkWH23F07f6dMBfOsJZQJr8xtt1OKOoPYWuVSJ3vOwNzt7GnFE0XU
/ZK1Df5kMrvyGvNw4ptJesToZtSSawS9hQidIL68RNN5h+foCVGwvpvr4mYlKHyb
84r8elBmAyyu3U5Zk4K8BkkCgYARAb2iiDfkJXo6Xfnhl8dF8f5CfAR3jmNPrZD8
hRtVJ7tCVWObivcO42nSKDq8XkPI7BYGbRkuo2vhnSK9wlLHW5aLImuImVcBPAWp
dlr3TX7EqxnAH+9z/9p/sEW58l7C6ZLgXFayq5zoCJOMaz9SG/mb/alGGqOcokV5
qbThwwKBgBDdUwd6RWYOpmmSsqB2GaxrO2+EjrlAXKSqNHj2aycOkJiCYyZi6eQZ
b3ePtnsPj6cxuP8Mh7jyVYn16WTFc/u4P4/AsQsjI0ky8Iwd98r5TDkqV4CEV2BZ
oDzCG1nht+zDmNE6surpue276FUSmaPOZ3L1PTqAA8H7NKQwiM33
-----END RSA PRIVATE KEY-----
`)

var testServerPublicCert = []byte(`-----BEGIN CERTIFICATE-----
MIID/DCCAuSgAwIBAgIJAJZ3p6LwP/ATMA0GCSqGSIb3DQEBBQUAMFsxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIEwpTb21lLVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQxFDASBgNVBAMTC2V4YW1wbGUuY29tMCAXDTE0MTEyMjA2
MjkyMFoYDzMwMTQwMzI1MDYyOTIwWjBbMQswCQYDVQQGEwJBVTETMBEGA1UECBMK
U29tZS1TdGF0ZTEhMB8GA1UEChMYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMRQw
EgYDVQQDEwtleGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC
ggEBAKOKaHw1RL1rq/bnSz1mP9Ce963ZvjtbyBUSyIfWTnwF+k9eUBxVPToSa9jT
vkwI+WhoLOwDZFGKaevc9tfBs+CTOl5jb88tDrRVvhvsEQ6eebqX4wgiaaFIDNQ1
tKNh2Br/tGd7OkNm4HJXj+5b4dzzPks3ndj8crFOwAun/jNbqXsgkyiX7EepgtdJ
ggkwxT8VljjSE567ZeoEaVDGkACMbJAmqgU/THiNUOygYVj4kKQhIiO0UnSMmv4E
QvXXoO10yV8GGDOHD4A8PVGV4s0vmA/7x4UNNs8MjA8+/3dbpsDRVJGg1D6Vg9jU
HDMSRWcaYT/lfaJ4uhUHHcd6fVkCAwEAAaOBwDCBvTAdBgNVHQ4EFgQUNoJ/9/eY
7GRpV3FG9s3vBp3V3vAwgY0GA1UdIwSBhTCBgoAUNoJ/9/eY7GRpV3FG9s3vBp3V
3vChX6RdMFsxCzAJBgNVBAYTAkFVMRMwEQYDVQQIEwpTb21lLVN0YXRlMSEwHwYD
VQQKExhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxFDASBgNVBAMTC2V4YW1wbGUu
Y29tggkAlnenovA/8BMwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQUFAAOCAQEA
RK6IyocmCl7lhRg3E9MbljIxO3Em1bzH5rlALQMnCnsiSKOdThZuzTBU9jmo80fo
1ox7vlQdIz/3QIImy1kGZ/tkzgII1qZYIY0S7oYIEYFp8UBfjYyan0LsjeGcXovl
qg7azw3z0fCAg3TZ/IugyxTmKk+j/1qGdiDMRChAmqnYtuC+vajGvUgP7JpeWeFR
00iJuIeAp5mjVwgKXdZWkZlVsQp91ijOlYKZGuAjRazWLw7wA7oYWRNrVmjrbg3b
xQ8SBVkAq3+ucw342yKjDtFop713AXNg5bhNakFof3A9SYD/nBENjhlW0C06gsR2
T88PyTE1M0XtM7XUffqqYw==
-----END CERTIFICATE-----
`)

//开启一个tls服务器监听,证书在上面是自签名的. 目标: 使各种客户端都可以兼容使用
//这个服务器不是安全的,1.本项目开源,2.所以上面的证书也是到处都可以取到的.
// 有下列两个用处
// 1.想用ssl,不管安全性
// 2.只是测试代码正确性
func MustTlsTestListen(addr string) (listener net.Listener) {
	listener, err := tls.Listen("tcp", addr, GetTlsTestServerConfig())
	if err != nil {
		panic(err)
	}
	return listener
}

var tlsTestServerConfigOnce sync.Once
var tlsTestServerConfig *tls.Config

func GetTlsTestServerConfig() *tls.Config {
	tlsTestServerConfigOnce.Do(func() {
		cert, err := tls.X509KeyPair(testServerPublicCert, testServerPrivateCert)
		if err != nil {
			panic(err)
		}
		tlsTestServerConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
	})
	return tlsTestServerConfig
}

// 返回一个不验证对方证书有效性,并且发送127.0.0.1域名给对方的拨号器
// 这个拨号器是不安全的,仅供测试或不需要安全的应用
/*
func NewTestTlsDialer(addr string) kmgNet.RwcDialer{
	return kmgNet.RwcDialerFunc(func(){
		tls.Dial("tcp", addr, &tls.Config{
				ServerName:         "127.0.0.1",
				InsecureSkipVerify: true,
			})
	})
}
	*/