package kmgConsole

import (
	"io"
	"os"
)

//把stderr和stdout重定向到文件里面
func MustStdoutErrAppendToFile(path string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0777))
	if err != nil {
		panic(err)
	}
	os.Stderr.Close()
	os.Stdout.Close()
	os.Stderr = f
	os.Stdout = f
}

//把stderr和stdout使用tee的方法,重定向到文件,并且也输出到原始的stdout,(原始的stderr被忽略)
//测试可用
//大概相当于 xxx 2>&1 | tee path.log
func MustStdoutErrTeeToFile(path string) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0777))
	if err != nil {
		panic(err)
	}
	oldStdOut := os.Stdout
	oldStdErr := os.Stderr
	oldStdErr.Close()
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w
	os.Stderr = w
	writer := io.MultiWriter(oldStdOut, f)
	go func() {
		defer f.Close()
		defer oldStdOut.Close()
		io.Copy(writer, r)
	}()
	return
}
