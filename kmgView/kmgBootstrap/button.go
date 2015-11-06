package kmgBootstrap
import "github.com/bronze1man/kmg/kmgView"

type ButtonSize string

const (
	ButtonSizeLarge      ButtonSize = "btn-lg"
	ButtonSizeDefault    ButtonSize = ""
	ButtonSizeSmall      ButtonSize = "btn-sm"
	ButtonSizeExtraSmall ButtonSize = "btn-xs"
)

type ButtonColor string

const (
	ButtonColorDefault ButtonColor = "btn-default"
	ButtonColorPrimary ButtonColor = "btn-primary"
	ButtonColorSuccess ButtonColor = "btn-success"
	ButtonColorInfo    ButtonColor = "btn-info"
	ButtonColorWarning ButtonColor = "btn-warning"
	ButtonColorDanger  ButtonColor = "btn-danger"
	ButtonColorLink    ButtonColor = "btn-link"
)

type ButtonType string

const (
	ButtonTypeA      ButtonType = "a"
	ButtonTypeButton ButtonType = "button"
)

type Button struct {
	Type          ButtonType
	Url           string
	Color         ButtonColor
	Size          ButtonSize
	Content       kmgView.HtmlRenderer
	AttributeNode kmgView.HtmlRenderer
	ClassName     string
	Id            string
	Name          string
	Value         string
}

func (b Button) HtmlRender() string {
	if b.Color == "" {
		b.Color = ButtonColorDefault
	}
	if b.Type == "" {
		b.Type = ButtonTypeButton
	}
	if b.Content == nil {
		panic("kmgBootstrap.Button.Content must not be empty")
	}
	return tplButton(b)
}

