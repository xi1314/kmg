package kmgStrings

func IsInSlice(slice []string, s string) bool {
	for _, thisS := range slice {
		if thisS == s {
			return true
		}
	}
	return false
}

//里面的字符串有3种状态,
// 1. false 表示需要检查,但是还没有遇到
// 2. true	表示需要检查,并且遇到了
// 3. 查不到 表示不需要检查,是否遇到过无关紧要.
type SliceExistChecker map[string]bool

//向Checker中插入一个字符串,返回该字符串是否在最开始要检查的列表里面.
func (c SliceExistChecker) Add(s string) bool {
	ret, ok := c[s]
	if !ok {
		return false
	}
	if ret == false {
		c[s] = true
	}
	return true
}
func (c SliceExistChecker) Check() (NotExist string) {
	for s, ret := range c {
		if ret == false {
			return s
		}
	}
	return ""
}

//新建一个数据是否存在的检查者
func NewSliceExistChecker(slice ...string) SliceExistChecker {
	out := SliceExistChecker{}
	for _, s := range slice {
		out[s] = false
	}
	return out
}

//是否s里面全部都是英文字母(只有那26个,大小写均可)
func IsAllAphphabet(s string) bool {
	for _, rune := range s {
		if !((rune >= 65 && //A
			rune <= 90) || //Z
			(rune >= 97 && //a
				rune <= 122)) { //z
			return false
		}
	}
	return true
}
