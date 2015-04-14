package kmgConsole

import "os"

//把stderr和stdout重定向到文件里面
func MustStdoutErrAppendToFile(path string) {
	os.Stderr.Close()
	os.Stdout.Close()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0777))
	if err != nil {
		panic(err)
	}
	os.Stderr = f
	os.Stdout = f
}
