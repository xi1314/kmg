package command

import (
	"bytes"
	"flag"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgTextTemplate"
	"text/template"
)

func init() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "rpc",
		Runner: Generate,
	})
}

type Request struct {
	PackagePath string
	TypeName    string
}

func Generate() {
	req := Request{}
	flag.StringVar(&req.PackagePath, "PackagePath", "", "the path of the package to generate rpc files")
	flag.StringVar(&req.TypeName, "TypeName", "", "type name")
	//phase 1 get reflect data of the type
	kmgTextTemplate.MustRender(`package main
import (
    inPkg "{{.Package}}"
    "reflect"
)

func main(){

}
    `, struct {
	}{})
	//phase 2 generate client and server code
	buf := &bytes.Buffer{}
	template.New("kmgRpc").Parse(``)
}
