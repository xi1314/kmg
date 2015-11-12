package kmgBootstrap

import (
	"github.com/bronze1man/kmg/kmgRand"
	"github.com/bronze1man/kmg/kmgView"
)

type DropDown struct {
	id         string
	Title      kmgView.HtmlRenderer
	OptionList []NavBarNode
}

func (d DropDown) HtmlRender() string {
	if d.id == "" {
		d.id = kmgRand.MustCryptoRandToReadableAlphaNum(20)
	}
	return tplDropDown(d)
}

func NewMoreButton(optionList []NavBarNode) DropDown {
	d := DropDown{
		Title: Button{
			Size: ButtonSizeExtraSmall,
			Content: Icon{
				IconName: "cog",
			},
		},
		OptionList: optionList,
	}
	return d
}

func NewCaret() kmgView.HtmlRenderer {
	return kmgView.Html(`<span class="caret"></span>`)
}
