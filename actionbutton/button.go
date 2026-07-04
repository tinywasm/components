package actionbutton

import (
	. "github.com/tinywasm/css"
	. "github.com/tinywasm/dom"
	. "github.com/tinywasm/html"
)

var (
	clsBtn             Class = "btn"
	clsBtnPrimary      Class = "btn btn-primary"
	clsBtnSecondary    Class = "btn btn-secondary"
	clsBtnDanger       Class = "btn btn-danger"
	clsBtnUrlUp        Class = "btn-url-up"
	clsBtnUrlDown      Class = "btn-url-down"
	clsBtnUrl          Class = "btn-url"
	clsBtnUrlDisable   Class = "btn-url-disable"
	clsBtnSelected     Class = "btn-selected"
	clsBtnLogin        Class = "btn-login"
	clsBtnUrlPulse     Class = "btn-url-pulse"
	clsContebuton      Class = "contebuton"
	clsContCenteredBtn Class = "cont-centered-btn"
)

type ActionButton struct {
	Element
	Text    string
	Variant string // "primary", "secondary", "danger"
	OnClick func(Event)
}

func (b *ActionButton) Render() *Element {
	var variantCls Class
	switch b.Variant {
	case "secondary":
		variantCls = clsBtnSecondary
	case "danger":
		variantCls = clsBtnDanger
	default:
		variantCls = clsBtnPrimary
	}

	btn := Button().Set(variantCls.AsAttr()).
		Text(b.Text)

	if b.OnClick != nil {
		btn.On("click", b.OnClick)
	}

	return btn
}
