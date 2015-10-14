package main

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgFile"
	"path/filepath"
	"time"
)

var dockerPath = filepath.Join(kmgConfig.DefaultEnv().ProjectPath, "src/github.com/bronze1man/kmg/kmg/SubCommand/serviceCmd/test")

func main() {
	kmgCmd.MustRun(`kmg GoCrossCompile -platform linux_amd64 github.com/bronze1man/kmg/kmg`)
	kmgCmd.MustRun(`kmg GoCrossCompile -platform linux_amd64 github.com/bronze1man/kmg/kmg/SubCommand/serviceCmd/testBin`)
	kmgFile.MustCopyFile(filepath.Join(kmgConfig.DefaultEnv().ProjectPath, "bin/kmg_linux_amd64"), filepath.Join(dockerPath, "kmg"))
	kmgFile.MustCopyFile(filepath.Join(kmgConfig.DefaultEnv().ProjectPath, "bin/testBin_linux_amd64"), filepath.Join(dockerPath, "testBin"))
	kmgFile.MustWriteFile(filepath.Join(dockerPath, "Dockerfile"), []byte(`FROM ubuntu
WORKDIR /
COPY kmg /bin/
COPY testBin /bin/
RUN chmod +x /bin/kmg
RUN chmod +x /bin/testBin
CMD testBin
`))
	//CMD kmg service setandrestart t testBin && kmg service stop t && kmg service start t && kmg service restart t
	kmgCmd.MustRunAndReturnOutput("docker build -t kmgtest " + dockerPath)
	for i := 0; i < 20; i++ {
		time.Sleep(time.Second)
		kmgCmd.MustRunAndReturnOutput("docker run kmgtest")
	}
}
