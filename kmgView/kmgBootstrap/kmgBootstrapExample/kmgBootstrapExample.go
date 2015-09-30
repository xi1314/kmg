package main
import (
	"github.com/bronze1man/kmg/kmgControllerRunner"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgView/kmgBootstrap"
	"github.com/bronze1man/kmg/kmgView"
)

func main(){
	kmgControllerRunner.RegisterController(Example{})
	kmgControllerRunner.EnterPointApiName = "main.Example.Demo1"
	kmgControllerRunner.StartServerCommand()
}

type Example struct{}

func (e Example) Demo1(ctx *kmgHttp.Context){
	ctx.WriteString(kmgBootstrap.NewWrap("kmg bootstrap demo",
		kmgBootstrap.Panel{
			Title: "Panel",
			Body: kmgBootstrap.Panel{
				Title: "Panel",
				Body: kmgView.String("Body"),
			},
		},
		kmgBootstrap.Panel{
			Title: "Table",
			Body: kmgBootstrap.Table{
				Caption: kmgView.String("Caption"),
				TitleList: []kmgView.HtmlRenderer{
					kmgView.String("title1"),
					kmgView.String("title2"),
				},
				DataList: [][]kmgView.HtmlRenderer{
					{
						kmgView.String("r1c1"),
						kmgView.String("r1c2"),
					},
					{
						kmgView.String("r2c1"),
						kmgView.String("r2c2"),
					},
				},
			},
		},
		kmgBootstrap.Panel{
			Title: "Icon",
			Body: kmgBootstrap.Icon{
				IconName:      "exchange",
				IconColor:     kmgBootstrap.TextDanger,
				AttributeNode: kmgBootstrap.Popover{Title: "这里应该有一个图标"},
			},
		},
	).HtmlRender())
}