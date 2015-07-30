package kmgSubmail

import (
	"github.com/bronze1man/kmg/encoding/kmgJson"
	"github.com/bronze1man/kmg/errors"
	"github.com/bronze1man/kmg/kmgLog"
	"github.com/bronze1man/kmg/kmgStrconv"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Mail struct {
	Project  SubmailProject
	Links    string
	To       string
	Vars     string
	Html     string
	From     string
	FromName string
	Title    string
}

type Submail struct {
	Appid     string
	Signature string
}

var EmailConfig = Submail{}

type SubmailProject string

func SendMail(email Mail) string {
	submailUrl := "https://api.submail.cn/mail/send.json"
	resp, err := http.PostForm(submailUrl, url.Values{
		"appid":     {EmailConfig.Appid},
		"signature": {EmailConfig.Signature},
		"to":        {email.To},
		"html":      {email.Html},
		"subject":   {email.Title},
		"from":      {email.From},
		"from_name": {email.FromName},
	})
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	kmgLog.Log("Submail", string(body), email)
	data := kmgJson.MustUnmarshalToMapDeleteBOM(body)
	if data["status"] == "error" {
		panic(errors.New(kmgStrconv.InterfaceToString(data["msg"])))
	}
	return string(body)
}

func XSendMail(email Mail) (err error) {
	submailUrl := "https://api.submail.cn/mail/xsend.json"
	resp, e := http.PostForm(submailUrl, url.Values{
		"appid":     {EmailConfig.Appid},
		"signature": {EmailConfig.Signature},
		"to":        {email.To},
		"project":   {kmgStrconv.InterfaceToString(email.Project)},
		"links":     {email.Links},
		"vars":      {email.Vars},
	})
	if e != nil {
		return e
	}
	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return e
	}
	kmgLog.Log("Submail", string(body), email)
	data := kmgJson.MustUnmarshalToMapDeleteBOM(body)
	if data["status"] == "error" {
		return errors.New(kmgStrconv.InterfaceToString(data["msg"]))

	}
	return nil
}

func MustXSendMail(email Mail) {
	err := XSendMail(email)
	if err != nil {
		panic(err)
	}
}
