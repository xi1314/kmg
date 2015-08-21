package kmgGoParser

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgGoSource/kmgGoReader"
)

// 此处仅跳过函数
func (gofile *File) readGoFunc(r *kmgGoReader.Reader) *FuncOrMethodDeclaration {
	funcDecl := &FuncOrMethodDeclaration{}
	r.MustReadMatch([]byte("func"))
	// 读函数头,至少要有括号
	r.ReadAllSpace()
	b := r.ReadByte()
	if b == '(' {
		r.UnreadByte()
		receiver := gofile.readParameters(r)
		if len(receiver) != 1 {
			panic(fmt.Errorf("%s receiver must have one parameter", r.GetFileLineInfo()))
		}
		funcDecl.ReceiverType = receiver[0].Type
		r.ReadAllSpace()
		// 暂不处理方法
	} else {
		r.UnreadByte()
	}
	id := readIdentifier(r)
	funcDecl.Name = string(id)
	if funcDecl.Name == "" {
		panic(fmt.Errorf("%s need function name", r.GetFileLineInfo()))
	}
	r.ReadAllSpace()
	b = r.ReadByte()
	if b != '(' {
		panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
	}
	r.UnreadByte()
	funcDecl.InParameter = gofile.readParameters(r)
	r.ReadAllSpaceWithoutLineBreak()
	b = r.ReadByte()
	if b == '\n' { //没有body
		return funcDecl
	} else if b != '{' {
		r.UnreadByte()
		funcDecl.OutParameter = gofile.readParameters(r)
		r.ReadAllSpaceWithoutLineBreak()
		b = r.ReadByte()
	}
	if b == '\n' { //没有body
		return funcDecl
	}
	if b != '{' {
		panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
	}

	//跳过函数体
	readMatchBigParantheses(r)
	return funcDecl
}

func (gofile *File) readParameters(r *kmgGoReader.Reader) (output []FuncParameter) {

	b := r.ReadByte()
	if b != '(' {
		// 处理 int 这种类型
		r.UnreadByte()
		return []FuncParameter{
			{
				Type: gofile.readType(r),
			},
		}
	}
	parameterPartList := []*astParameterPart{}
	lastPart := &astParameterPart{}
	for {
		r.ReadAllSpace()
		b := r.ReadByte()
		if b == ')' || b == ',' {
			if lastPart.partList[0].originByte != nil {
				parameterPartList = append(parameterPartList, lastPart)
				lastPart = &astParameterPart{}
			}
			if b == ')' {
				break
			}
			if b == ',' {
				continue
			}
		}

		r.UnreadByte()
		if r.IsMatchAfter([]byte("...")) {
			r.MustReadMatch([]byte("..."))
			lastPart.isVariadic = true
		}
		startPos := r.Pos()
		typ := gofile.readType(r)
		buf := r.BufToCurrent(startPos)
		//fmt.Println(string(buf))
		hasSet := false
		for i := range lastPart.partList {
			if lastPart.partList[i].originByte == nil {
				lastPart.partList[i].originByte = buf
				lastPart.partList[i].typ = typ
				hasSet = true
				break
			}
		}
		if !hasSet {
			panic(r.GetFileLineInfo() + " unexcept func parameterList.")
		}
	}

	output = make([]FuncParameter, len(parameterPartList))
	onlyHavePart1Num := 0
	for i := range parameterPartList {
		if parameterPartList[i].partList[1].originByte == nil {
			onlyHavePart1Num++
		}
	}
	//重新分析出函数参数来. 处理(int,int) 这种类型.
	if onlyHavePart1Num == len(parameterPartList) {
		for i := range parameterPartList {
			output[i].Type = parameterPartList[i].partList[0].typ
			output[i].IsVariadic = parameterPartList[i].isVariadic
		}
		return output
	}
	// 处理 (x,y int) (x int,y int) 这种类型.
	for i, parameterPart := range parameterPartList {
		output[i].Name = string(parameterPart.partList[0].originByte)
		if parameterPart.partList[1].typ != nil {
			output[i].Type = parameterPart.partList[1].typ
		}
		output[i].IsVariadic = parameterPart.isVariadic
	}
	// 补全 (x,y int) 里面 x的类型.
	for i := range parameterPartList {
		if output[i].Type == nil {
			for j := i + 1; j < len(parameterPartList); j++ {
				if output[j].Type != nil {
					output[i].Type = output[j].Type
				}
			}
		}
	}
	return output
}

type astParameterPart struct {
	partList [2]struct {
		originByte []byte
		typ        Type
	}
	isVariadic bool //是否有3个点
}
