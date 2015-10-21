package kmgSubmail

import (
	"github.com/bronze1man/kmg/encoding/kmgJson"
	"github.com/bronze1man/kmg/errors"
	"github.com/bronze1man/kmg/kmgDebug"
	"github.com/bronze1man/kmg/kmgLog"
	"github.com/bronze1man/kmg/kmgStrconv"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Message struct {
	Project MessageProject
	To      string //手机号
	Vars    string // 也可用json传入多个变量名哟

}

type SubMessage struct {
	Appid     string
	Signature string
}

var MessageConfig = SubMessage{}

type MessageProject string

func XSendMessage(message Message) (err error) {
	subMessageUrl := "https://api.submail.cn/message/xsend.json"
	kmgDebug.Println(message.Vars)
	response, err := http.PostForm(subMessageUrl, url.Values{
		"appid":     {MessageConfig.Appid},
		"signature": {MessageConfig.Signature},
		"to":        {message.To},
		"project":   {kmgStrconv.InterfaceToString(message.Project)},
		"vars":      {message.Vars},
	})
	defer response.Body.Close()
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	kmgLog.Log("SubMessage", string(body), message)
	data := kmgJson.MustUnmarshalToMapDeleteBOM(body)
	if data["status"] == "error" {
		return errors.New(kmgStrconv.InterfaceToString(data["msg"]))
	}
	return nil

}

func MustXSendMessage(message Message) {
	err := XSendMessage(message)
	if err != nil {
		panic(err)
	}
}
