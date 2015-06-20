package kmgSubmail

import (
	"github.com/bronze1man/kmg/kmgLog"
	"github.com/bronze1man/kmg/kmgStrconv"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Mail struct {
	Project SubmailProject
	Links   string
	To      string
	Vars    string
	Html    string
	From    string
	Title   string
}

type DefaultMail struct {
	Appid     string
	Signature string
}

var EmailConfig = DefaultMail{}

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
	})
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}

func XSendMail(email Mail) string {
	submailUrl := "https://api.submail.cn/mail/xsend.json"
	resp, err := http.PostForm(submailUrl, url.Values{
		"appid":     {EmailConfig.Appid},
		"signature": {EmailConfig.Signature},
		"to":        {email.To},
		"project":   {kmgStrconv.InterfaceToString(email.Project)},
		"links":     {email.Links},
		"vars":      {email.Vars},
	})
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	kmgLog.Log("response", string(body), email)
	return string(body)
}
