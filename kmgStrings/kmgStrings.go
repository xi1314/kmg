package kmgStrings

import (
	"bytes"
	"sort"
	"strings"
	"unicode"
)

func IsInSlice(slice []string, s string) bool {
	for _, thisS := range slice {
		if thisS == s {
			return true
		}
	}
	return false
}

func SliceNoRepeatMerge(s1 []string, s2 []string) []string {
	for _, s := range s2 {
		if !IsInSlice(s1, s) {
			s1 = append(s1, s)
		}
	}
	return s1
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

//是否s里面全部都是英文字母(只有那26个,大小写均可) 空字符串也是true
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

//是否s里面全部都是数值(只有那10个,0-9) 空字符串也是true
func IsAllNum(s string)bool{
	for _,rune:=range s{
		if !(rune>='0' && rune<='9'){
			return false
		}
	}
	return true
}

//仅对 UTF-8 有效
func FirstLetterToUpper(s string) string {
	b := []byte(s)
	if len(b) == 0 {
		return s
	}
	if b[0] > unicode.MaxASCII {
		return s
	}
	firstLetter := bytes.ToUpper(b[0:1])
	b = append(firstLetter, b[1:]...)
	return string(b)
}

func MapStringBoolToSortedSlice(m map[string]bool) []string {
	output := make([]string, len(m))
	i := 0
	for s := range m {
		output[i] = s
		i++
	}
	sort.Strings(output)
	return output
}

// LastTwoPartSplit("github.com/bronze1man/kmg/kmgGoSource/testPackage.Demo",".") -> "github.com/bronze1man/kmg/kmgGoSource/testPackage","Demo",false
// LastTwoPartSplit("Demo",".") -> "","",true
func LastTwoPartSplit(originS string, splitS string) (p1 string, p2 string, ok bool) {
	part := strings.Split(originS, splitS)
	if len(part) < 2 {
		return "", "", false
	}
	return strings.Join(part[:len(part)-1], splitS), part[len(part)-1], true
}

func LineDataToSlice(lineData string) []string {
	part := strings.Split(lineData, "\n")
	out := []string{}
	for _, s := range part {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		out = append(out, s)
	}
	return out
}
func SubStr(s string, from int, to int) string {
	rs := []rune(s)
	rl := len(rs)
	if to == 0 {
		to = rl
	}
	if to < 0 {
		to = rl + to
	}
	if to > rl {
		to = rl
	}
	if to < 0 {
		to = 0
	}
	return string(rs[from:to])
}
func IsStartWith(s string,start string)bool{
	return  SubStr(s,0,1) == start
}
