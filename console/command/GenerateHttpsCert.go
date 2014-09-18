package command

import (
	"flag"
	"fmt"
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgFile"
	"os"
	"path/filepath"
	"strings"
)

//https证书生成,会先生成一个根证书,然后生成几个客户端证书
type GenerateHttpsCert struct {
	outputPath string
	subject    string
}

func (command *GenerateHttpsCert) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{
		Name:  "GenerateHttpsCert",
		Short: "Generate Https Cert",
	}
}

func (command *GenerateHttpsCert) ConfigFlagSet(f *flag.FlagSet) {
	f.StringVar(&command.outputPath, "o", "certs", "output dir,it will create it if it is not created")
	f.StringVar(&command.subject, "subj", "/C=CN/ST=SiChuan/L=ChengDu/O=ZhuoZhuo/OU=IT Department/CN=www.new1.uestc.edu.cn", "the subj of the cert.")
}

func (command *GenerateHttpsCert) Execute(context *console.Context) (err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	workDir := filepath.Join(wd, "certs")
	kmgFile.MustMkdirAll(workDir)
	os.Chdir(workDir)
	kmgFile.MustWriteFile("index.txt", []byte(""))
	kmgFile.MustWriteFile("serial", []byte("01"))
	kmgFile.MustWriteFile("config.conf", []byte(`[ ca ]
default_ca = ca_default

[ ca_default ]
dir = .
certs = .
new_certs_dir = .
database = ./index.txt
serial = ./serial
RANDFILE = .rand
certificate = ca.crt
private_key = ca.key
default_days = 36500
default_crl_days = 30
default_md = md5
preserve = no
policy = generic_policy
[ policy_anything ]
countryName = optional
stateOrProvinceName = optional
localityName = optional
organizationName = optional
organizationalUnitName = optional
commonName = supplied
emailAddress = optional`))
	mustRunCmd("openssl req -passout pass:1234 -batch -new -x509 -newkey rsa:2048 -extensions v3_ca -keyout ca.key -out ca.crt -days 18250",
		"-subj", command.subject+" ca")
	mustRunCmd("openssl req -new -newkey rsa:2048 -nodes -keyout server.key -out csr.csr -days 18250",
		"-subj", command.subject)
	mustRunCmdString("openssl ca -config config.conf -batch -cert ca.crt -passin pass:1234 -keyfile ca.key -policy policy_anything -out server.crt -infiles csr.csr")
	mustRunCmd("openssl req -new -newkey rsa:2048 -nodes -keyout client.key -out csr.csr -days 18250",
		"-subj", command.subject+" client")
	mustRunCmdString("openssl ca -config config.conf -batch -cert ca.crt -passin pass:1234 -keyfile ca.key -policy policy_anything -out client.crt -infiles csr.csr")
	mustRunCmdString("openssl pkcs12 -export -passout pass:1234 -inkey client.key -in client.crt -out client.pfx")
	return nil
}

func mustRunCmdString(s string) {
	fmt.Println(s)
	err := kmgCmd.RunOsStdioCmdString(s)
	if err != nil {
		panic(err)
	}
}

func mustRunCmd(s string, args ...string) {
	fmt.Println(s)

	sParts := strings.Split(s, " ")
	args = append(sParts[1:], args...)
	err := kmgCmd.RunOsStdioCmd(sParts[0], args...)
	if err != nil {
		panic(err)
	}
}
