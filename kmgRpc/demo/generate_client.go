package demo

import (
	"errors"
	"github.com/bronze1man/kmg/encoding/kmgGob"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/qiniu/bytes"
	"net/http"
)

//还可以生成一个函数版本(简化客户端的客户的使用)
type Client struct {
	addr string
}

//简单函数
func (c *Client) PostScoreInt(LbId string, Score int) error {
	req := RpcRequest{
		ApiName: "PostScoreInt",
		InData: [][]byte{
			kmgGob.MustMarshal(LbId),
			kmgGob.MustMarshal(Score),
		},
	}
	inByte := kmgGob.MustMarshal(req)
	httpResp, err := http.Post("http://"+c.addr+"/", "object", bytes.NewReader(inByte))
	if err != nil {
		return err
	}
	outBytes, err := kmgHttp.ResponseReadAllBody(httpResp)
	if err != nil {
		return err
	}
	var resp RpcResponse
	kmgGob.MustUnmarshal(outBytes, &resp)
	if resp.Error != "" {
		return errors.New(resp.Error)
	}
	return nil
}

func (s *Client) GetMaxScoreInt(LbId string) (int, error) {
	return 0, nil
}

func (s *Client) AutoRegister(req *AutoRegisterRequest) (Id string, SK string, err error) {
	return "", "", nil
}

func (s *Client) GetFrontUserInfo(Id string, Sk string) (info *FrontUserInfo, err error) {
	return nil, nil
}

//允许服务器端不返回错误,但是客户端总是会有错误(网络错误)
func (s *Client) SimpleAdd(a int, b int) (int, error) {
	return a + b, nil
}

func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
	}
}
