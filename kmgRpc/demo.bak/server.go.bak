package demo

var lastPostLbId string
var lastPostScore int

type Demo struct {
}

//简单函数
func (s *Demo) PostScoreInt(LbId string, Score int) error {
	lastPostLbId = LbId
	lastPostScore = Score
	return nil
}

func (s *Demo) GetMaxScoreInt(LbId string) (int, error) {
	return 0, nil
}

func (s *Demo) AutoRegister(req *AutoRegisterRequest) (Id string, SK string, err error) {
	return "", "", nil
}

func (s *Demo) GetFrontUserInfo(Id string, Sk string) (info *FrontUserInfo, err error) {
	return nil, nil
}

//允许服务器端不返回错误,但是客户端总是会有错误(网络错误)
func (s *Demo) SimpleAdd(a int, b int) int {
	return a + b
}
