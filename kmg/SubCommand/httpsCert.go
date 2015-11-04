package SubCommand
import (
	"flag"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgFile"
	"os"
)

func httpsCertCsrCLI(){
	domain:=""
	outDir:=""
	flag.StringVar(&domain,"domain","","the domain need to generate csr(google.com)")
	flag.StringVar(&outDir,"outDir","","the output dir(default to ./doc/cert/{domain})")
	flag.Parse()
	if domain==""{
		flag.Usage()
		os.Exit(1)
	}
	if outDir==""{
		outDir = "doc/cert/"+domain
	}
	kmgFile.Mkdir(outDir)
	kmgCmd.CmdSlice([]string{"openssl", "req", "-out", "domain.csr", "-new", "-newkey", "rsa:4096", "-nodes",
		"-keyout" ,"domain.key" ,"-subj", "/C=US/ST=US/L=US/O=US/OU=US/CN="+domain+"/"}).SetDir(outDir).ProxyRun()
}