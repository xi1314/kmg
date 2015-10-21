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

type Massege struct {
	ProjectId string
	To        string
}

type Submail struct {
	Appid     string
	Signature string
}

var EmailConfig = Submail{}

type SubmailProject string

func SendMailForHtml(email Mail) (err error) {
	submailUrl := "https://api.submail.cn/mail/send.json"
	return MailBreack(email, url.Values{
		"appid":     {EmailConfig.Appid},
		"signature": {EmailConfig.Signature},
		"to":        {email.To},
		"html":      {email.Html},
		"subject":   {email.Title},
		"from":      {email.From},
		"from_name": {email.FromName},
	}, submailUrl)
}

func SendMailForModel(email Mail) (err error) {
	submailUrl := "https://api.submail.cn/mail/xsend.json"
	return MailBreack(email, url.Values{
		"appid":     {EmailConfig.Appid},
		"signature": {EmailConfig.Signature},
		"to":        {email.To},
		"project":   {kmgStrconv.InterfaceToString(email.Project)},
		"links":     {email.Links},
		"vars":      {email.Vars},
	}, submailUrl)
}

func MailBreack(email Mail, u url.Values, submailUrl string) (err error) {
	resp, e := http.PostForm(submailUrl, u)
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
	err := SendMailForModel(email)
	if err != nil {
		defer func() {
			err := recover()
			kmgLog.Log("error", err)
		}()
		panic(err)
	}
}

func MustSendMail(email Mail) {
	err := SendMailForHtml(email)
	if err != nil {
		defer func() {
			err := recover()
			kmgLog.Log("error", err)
		}()
		panic(err)
	}
}

func SendMessage(massage Massege) (err error) {
	massegeApiUrl := "https://api.submail.cn/message/xsend"
	resp, e := http.PostForm(massegeApiUrl, url.Values{
		//		"appid":     {EmailConfig.Appid},
		//		"signature": {EmailConfig.Signature},
		"appid":     {"10111"},
		"signature": {"142a3e0d66c4dda1e918487b1952b26c"},
		"to":        {massage.To},
		"project":   {massage.ProjectId},
	})
	if e != nil {
		return e
	}
	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return e
	}
	kmgLog.Log("Submail", string(body), massage)
	data := kmgJson.MustUnmarshalToMapDeleteBOM(body)
	if data["status"] == "error" {
		return errors.New(kmgStrconv.InterfaceToString(data["msg"]))
	}
	return nil
}

func MustXsendMessage(massage Massege) {
	err := SendMessage(massage)
	if err != nil {
		defer func() {
			err := recover()
			kmgLog.Log("error", err)
		}()
		panic(err)
	}
}
