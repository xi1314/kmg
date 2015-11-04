package kmgJpush

import (
	"fmt"
	"github.com/bronze1man/kmg/errors"
	"github.com/bronze1man/kmg/kmgLog"
	"github.com/isdamir/jpush"
)

type Client struct {
	c            *jpush.PushClient
	IsIosProduct bool
	name         string
	Platform     SystemPlatform
}

type NewClientRequest struct {
	Name         string //这个客户端的名字
	AppKey       string
	Secret       string
	Platform     SystemPlatform //iOS、Android平台
	IsIosProduct bool           //如果是否false表示向测试设备推送,如果是true表示向正式设备推送,后台的那个开发与正式似乎没有作用.
}

type SystemPlatform string

const (
	Android     SystemPlatform = "Android"
	IOS         SystemPlatform = "iOS"
	AllPlatform SystemPlatform = "AllPlatform"
)

func NewClient(req NewClientRequest) *Client {
	return &Client{
		c:            jpush.NewPushClient(req.Secret, req.AppKey),
		IsIosProduct: req.IsIosProduct,
		name:         req.Name,
		Platform:     req.Platform,
	}
}

var NotFoundUser error = errors.New("[kmgJpush] user not exist")

func (c *Client) PushToOne(alias string, content string) (err error) {
	nb := jpush.NewNoticeBuilder()
	nb.SetPlatform(jpush.AllPlatform())
	au := &jpush.Audience{}
	au.SetAlias([]string{alias})
	nb.SetAudience(au)

	//Android配置
	notice := jpush.NewNoticeAndroid()
	notice.Alert = content
	nb.SetNotice(notice)

	//iOS配置
	iosNotice := jpush.NewNoticeIos()
	iosNotice.Sound = "default"
	iosNotice.Badge = "1"
	iosNotice.Alert = content
	nb.SetNotice(iosNotice)

	op := jpush.NewOptions()
	op.SetApns_production(c.IsIosProduct)
	nb.SetOptions(op)
	ret, err := c.c.Send(nb)
	if err != nil {
		return err
	}
	if ret.Error.Code == 0 {
		kmgLog.Log("jpush", "PushToOne success", c.name, alias, content, c.IsIosProduct)
		return nil
	}
	if ret.Error.Code == 1011 {
		kmgLog.Log("jpush", "PushToOne NotFoundUser", c.name, alias, content)
		return NotFoundUser
	}
	return fmt.Errorf("code:%d err: %s", ret.Error.Code, ret.Error.Message)
}

func (c *Client) PushToAll(content string) (err error) {
	nb := jpush.NewNoticeBuilder()
	nb.SetPlatform(jpush.AllPlatform())
	nb.SetAudience(jpush.AllAudience())

	// Android配置
	notice := jpush.NewNoticeAndroid()
	notice.Alert = content
	nb.SetNotice(notice)

	// iOS配置
	iosNotice := jpush.NewNoticeIos()
	iosNotice.Sound = "default"
	iosNotice.Badge = "1"
	iosNotice.Alert = content
	nb.SetNotice(iosNotice)

	op := jpush.NewOptions()
	op.SetApns_production(c.IsIosProduct)
	//	op.SetBigPushDuration(60) //过快的进行全局推送,会导致系统其他地方压力太大而挂掉.先设置成60分钟.
	nb.SetOptions(op)
	ret, err := c.c.Send(nb)
	if err != nil {
		return err
	}
	if ret.Error.Code == 0 {
		kmgLog.Log("jpush", "PushToAll success", c.name, content)
		return nil
	}
	if ret.Error.Code == 1011 {
		kmgLog.Log("jpush", "PushToAll NotFoundUser", c.name, content)
		return NotFoundUser
	}
	return fmt.Errorf("code:%d err: %s", ret.Error.Code, ret.Error.Message)
}
