package kmgGoParser

import (
	"bytes"
	"fmt"
	"unicode"
)

// 此处仅跳过函数
func (gofile *File) readGoFunc(r *reader) *FuncOrMethodDeclaration {
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

func (gofile *File) readParameters(r *reader) (output []FuncParameter) {

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
		startPos := r.pos
		typ := gofile.readType(r)
		buf := r.buf[startPos:r.pos]
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

/*
第一个字符可能为
	letter -> identifier(单独的类型名,带package的类型的package部分)
	"struct" struct类型开头
	"func" func类型开头
	"interface" interface类型开头
	"*" 指针开头
	"[" (数组,slice) 开头
	"map[" map开头
	"chan " chan开头
	"chan<- " chan<- 开头
	"<-chan" chan<- 开头
*/
func (gofile *File) readType(r *reader) Type {
	id := readIdentifier(r)
	if len(id) == 0 {
		if r.IsMatchAfter([]byte("<-chan")) {
			r.MustReadMatch([]byte("<-chan"))
			r.ReadAllSpace()
			return ChanType{
				Dir:  RecvDir,
				Elem: gofile.readType(r),
			}
		}
		b := r.ReadByte()
		if b == '*' {
			return PointerType{
				Elem: gofile.readType(r),
			}
		} else if b == '[' {
			content := readMatchMiddleParantheses(r)
			if len(content) == 1 {
				return SliceType{
					Elem: gofile.readType(r),
				}
			} else {
				// 仅跳过
				return ArrayType{
					Elem: gofile.readType(r),
				}
			}
		} else if b == '(' {
			typ := gofile.readType(r)
			r.MustReadMatch([]byte(")"))
			return typ
		} else {
			panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
		}
	}
	if bytes.Equal(id, []byte("struct")) {
		// 仅跳过
		r.ReadAllSpace()
		b := r.ReadByte()
		if b != '{' {
			panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
		}
		readMatchBigParantheses(r)
		return StructType{}
	} else if bytes.Equal(id, []byte("interface")) {
		// 仅跳过
		r.ReadAllSpace()
		b := r.ReadByte()
		if b != '{' {
			panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
		}
		readMatchBigParantheses(r)
		return InterfaceType{}
	} else if bytes.Equal(id, []byte("map")) {
		b := r.ReadByte()
		m := MapType{}
		if b != '[' {
			panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
		}
		m.Key = gofile.readType(r)
		r.MustReadMatch([]byte("]"))
		m.Value = gofile.readType(r)
		return m
	} else if bytes.Equal(id, []byte("func")) {
		// 仅跳过
		r.ReadAllSpace()
		b := r.ReadByte()
		if b != '(' {
			panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
		}
		readMatchSmallParantheses(r) //跳过输入参数
		r.ReadAllSpaceWithoutLineBreak()
		run := r.ReadRune() //跳过输出参数
		if run == '(' {
			readMatchSmallParantheses(r)
		} else if run == '\n' { //换行符可以标识这个函数定义结束了.
			return FuncType{}
		} else if unicode.IsLetter(run) || run == '[' || run == '*' || run == '<' {
			r.UnreadRune() //输出参数只有一个类型
			gofile.readType(r)
		} else {
			r.UnreadRune() //读到了其他东西,退回.
		}
		return FuncType{}
	} else if bytes.Equal(id, []byte("chan")) {
		if r.IsMatchAfter([]byte("<-")) {
			r.MustReadMatch([]byte("<-"))
			r.ReadAllSpace()
			return ChanType{
				Dir:  SendDir,
				Elem: gofile.readType(r),
			}
		} else {
			r.ReadAllSpace()
			return ChanType{
				Dir:  BothDir,
				Elem: gofile.readType(r),
			}
		}
	} else {
		b := r.ReadByte()
		if b == '.' {
			pkgPath := string(id)
			pkgPath, err := gofile.LookupFullPackagePath(pkgPath)
			if err != nil {
				fmt.Println(r.GetFileLineInfo(), err.Error()) //TODO 以目前的复杂度暂时无解.需要把所有相关的package都看一遍才能正确.
			}
			id2 := readIdentifier(r)
			return &NamedType{
				PackagePath: pkgPath,
				Name:        string(id2),
			}
		} else {
			r.UnreadByte()
			name := string(id)
			if builtinTypeMap[BuiltinType(name)] != Invalid {
				return BuiltinType(name)
			} else {
				return &NamedType{
					PackagePath: gofile.PackagePath,
					Name:        string(id),
				}
			}
		}
	}
	/*
		}else if r.IsMatchAfter([]byte("struct")) { //TODO 解决struct后面必须有一个空格的问题.
			r.ReadAllSpace()
			b := r.ReadByte()
			if b!='{' {
				panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
			}
			readMatchBigParantheses(r)
			return StructType{}
		}else if r.IsMatchAfter([]byte("interface")) {
			r.ReadAllSpace()
			b := r.ReadByte()
			if b!='{' {
				panic(fmt.Errorf("%s unexcept %s", r.GetFileLineInfo(), string(rune(b))))
			}
			readMatchBigParantheses(r)
			return InterfaceType{}
		}else if r.IsMatchAfter([]byte("map[")){

		}else if r.IsMatchAfter([]byte("func")){
			//TODO finish it.
		}
	*/
}
