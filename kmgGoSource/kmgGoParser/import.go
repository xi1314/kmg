package kmgGoParser

// 读取一个import语法里面的数据,此处从import关键词开始
func (gofile *File) readImport(r *reader) {
	r.MustReadMatch([]byte("import"))
	r.ReadAllSpace()
	if r.IsEof() {
		panic(r.GetFileLineInfo() + " unexcept EOF ")
	}
	b := r.ReadByte()
	if b == '(' {
		for {
			r.ReadAllSpace()
			b := r.ReadByte()
			if b == ')' {
				return
			} else {
				r.UnreadByte()
				gofile.readImportSpec(r)
			}
		}
	} else {
		r.UnreadByte()
		gofile.readImportSpec(r)
	}
}

// 此处暂时只保留路径,其他数据抛弃.
func (gofile *File) readImportSpec(r *reader) {
	r.ReadAllSpace()
	b := r.ReadByte()
	//fmt.Println(b,string(rune(b)))
	if b == '"' || b == '`' {
		r.UnreadByte()
		gofile.AddImport(string(mustReadGoString(r)), "")
	} else if b == '.' {
		r.ReadAllSpace()
		gofile.AddImport(string(mustReadGoString(r)), ".")
	} else {
		r.UnreadByte()
		alias := readIdentifier(r)
		r.ReadAllSpace()
		gofile.AddImport(string(mustReadGoString(r)), string(alias))
	}
}
