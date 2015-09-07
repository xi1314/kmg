package kmgGoParser

import (
	"github.com/bronze1man/kmg/kmgGoSource/kmgGoReader"
	//	"fmt"
)

// 此处暂时仅跳过该部分
func (gofile *File) readGoType(r *kmgGoReader.Reader) {
	r.MustReadMatch([]byte("type"))
	r.ReadAllSpace()
	b := r.ReadByte()
	if b == '(' {
		for {
			r.ReadAllSpace()
			b = r.ReadByte()
			if b == ')' {
				return
			}
			r.UnreadByte()
			name := readIdentifier(r)
			r.ReadAllSpace()
			typ := gofile.readType(r)
			gofile.NamedTypeList = append(gofile.NamedTypeList, &NamedType{
				PackagePath: gofile.PackagePath,
				Name:        string(name),
				underType:   typ,
			})
		}
		return
	} else {
		r.UnreadByte()
		name := readIdentifier(r)
		r.ReadAllSpace()
		typ := gofile.readType(r)
		gofile.NamedTypeList = append(gofile.NamedTypeList, &NamedType{
			PackagePath: gofile.PackagePath,
			Name:        string(name),
			underType:   typ,
		})
		return
	}
}

// 正确跳过该部分.
func (gofile *File) readGoVar(r *kmgGoReader.Reader) {
	r.MustReadMatch([]byte("var"))
	r.ReadAllSpace()
	b := r.ReadByte()
	if b == '(' {
		readMatchSmallParantheses(r)
		return
	}
	r.UnreadByte()
	readIdentifier(r)
	r.ReadAllSpace()
	b = r.ReadByte()
	if b == '=' {
		r.ReadAllSpace()
	}
	for {
		if b == '"' || b == '`' {
			r.UnreadByte()
			mustReadGoString(r)
		}
		if b == '\'' {
			r.UnreadByte()
			mustReadGoChar(r)
		}
		if b == '\n' {
			return
		}
		if b == '{' { //里面可以字面写东西.
			readMatchBigParantheses(r)
			//TODO 正确跳过该部分
			/*
							简单解决下列情况
							var UnreadRuneErrorTests = []struct {
					name string
					f    func(*Reader)
				}{
					{"Read", func(r *Reader) { r.Read([]byte{0}) }},
					{"ReadByte", func(r *Reader) { r.ReadByte() }},
					{"UnreadRune", func(r *Reader) { r.UnreadRune() }},
					{"Seek", func(r *Reader) { r.Seek(0, 1) }},
					{"WriteTo", func(r *Reader) { r.WriteTo(&Buffer{}) }},
				}

			*/
		}
		if b == '(' { //里面可以调用函数
			readMatchSmallParantheses(r)
		}
		if r.IsEof() {
			return
		}
		b = r.ReadByte()
	}
}

// 正确跳过该部分.
func (gofile *File) readGoConst(r *kmgGoReader.Reader) {
	r.MustReadMatch([]byte("const"))
	r.ReadAllSpace()
	b := r.ReadByte()
	if b == '(' {
		readMatchSmallParantheses(r)
		return
	}
	for {
		if b == '"' || b == '`' {
			r.UnreadByte()
			mustReadGoString(r)
		}
		if b == '\'' {
			r.UnreadByte()
			mustReadGoChar(r)
		}
		if b == '\n' {
			return
		}
		b = r.ReadByte()
	}
}
