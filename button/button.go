package button

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/dom"
)

var (
	ClsBtn             css.Class = "btn"
	ClsBtnPrimary      css.Class = "btn-primary"
	ClsBtnSecondary    css.Class = "btn-secondary"
	ClsBtnDanger       css.Class = "btn-danger"
	ClsBtnUrlUp        css.Class = "btn-url-up"
	ClsBtnUrlDown      css.Class = "btn-url-down"
	ClsBtnUrl          css.Class = "btn-url"
	ClsBtnUrlDisable   css.Class = "btn-url-disable"
	ClsBtnSelected     css.Class = "btn-selected"
	ClsBtnLogin        css.Class = "btn-login"
	ClsBtnUrlPulse     css.Class = "btn-url-pulse"
	ClsContebuton      css.Class = "contebuton"
	ClsContCenteredBtn css.Class = "cont-centered-btn"
)

type Button struct {
	*dom.Element
	Text    string
	Variant string // "primary", "secondary", "danger"
	OnClick func(dom.Event)
}

func (b *Button) Render() *dom.Element {
	if b.Element == nil {
		b.Element = &dom.Element{}
	}

	var variantCls css.Class
	switch b.Variant {
	case "secondary":
		variantCls = ClsBtnSecondary
	case "danger":
		variantCls = ClsBtnDanger
	default:
		variantCls = ClsBtnPrimary
	}

	btn := dom.Button(dom.Classes(ClsBtn, variantCls)).
		Text(b.Text)

	if b.OnClick != nil {
		btn.On("click", b.OnClick)
	}

	return btn
}
