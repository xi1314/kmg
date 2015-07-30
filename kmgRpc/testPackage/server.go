package testPackage

import (
	"errors"
)

type Demo struct {
}

func (s *Demo) PostScoreInt(LbId string, Score int) (Info string, err error) {
	if Score == 1 {
		return LbId, nil
	} else {
		return "", errors.New("Score!=1")
	}
}

type DemoRequest struct {
}

func (s *Demo) DemoFunc2(Req1 DemoRequest, Req2 *DemoRequest) (err error) {
	return nil
}

//返回参数没有名字
func (s *Demo) DemoFunc3(Req1 DemoRequest, Req2 *DemoRequest) error {
	return nil
}

//没有返回参数
func (s *Demo) DemoFunc4(Req1 DemoRequest, Req2 *DemoRequest) {
	return
}

//有返回参数,但是不包含error
func (s *Demo) DemoFunc5(Req1 DemoRequest, Req2 *DemoRequest) (Info string) {
	return ""
}

//不管私有方法
func (s *Demo) demoFunc6(Req1 DemoRequest, Req2 *DemoRequest) (Info string) {
	return ""
}

//返回值变成一个参数
func (s *Demo) DemoFunc7(Req1 DemoRequest, Req2 *DemoRequest) (Response string) {
	return ""
}

// 允许有 小写的参数名
func (s *Demo) DemoFunc8(req1 DemoRequest, req2 *DemoRequest, req3 int) (info string, err error) {
	if req3 == 1 {
		return "info1", nil
	}
	return "info", nil
}
