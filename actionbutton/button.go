package actionbutton

import (
	. "github.com/tinywasm/dom"
	. "github.com/tinywasm/html"
	"github.com/tinywasm/widget"
)

// NameActionButton is the widget name for actionbutton.
const NameActionButton = widget.Name("actionbutton")

const (
	PartPrimary     = widget.Part("primary")
	PartSecondary   = widget.Part("secondary")
	PartDanger      = widget.Part("danger")
	PartUrlUp       = widget.Part("url-up")
	PartUrlDown     = widget.Part("url-down")
	PartUrl         = widget.Part("url")
	PartUrlDisable  = widget.Part("url-disable")
	PartSelected    = widget.Part("selected")
	PartLogin       = widget.Part("login")
	PartUrlPulse    = widget.Part("url-pulse")
	PartContebuton  = widget.Part("contebuton")
	PartCenteredBtn = widget.Part("centered-btn")
)

var (
	clsBtn             = NameActionButton.Root()
	clsBtnPrimary      = NameActionButton.Class(PartPrimary)
	clsBtnSecondary    = NameActionButton.Class(PartSecondary)
	clsBtnDanger       = NameActionButton.Class(PartDanger)
	clsBtnUrlUp        = NameActionButton.Class(PartUrlUp)
	clsBtnUrlDown      = NameActionButton.Class(PartUrlDown)
	clsBtnUrl          = NameActionButton.Class(PartUrl)
	clsBtnUrlDisable   = NameActionButton.Class(PartUrlDisable)
	clsBtnSelected     = NameActionButton.Class(PartSelected)
	clsBtnLogin        = NameActionButton.Class(PartLogin)
	clsBtnUrlPulse     = NameActionButton.Class(PartUrlPulse)
	clsContebuton      = NameActionButton.Class(PartContebuton)
	clsContCenteredBtn = NameActionButton.Class(PartCenteredBtn)
)

type ActionButton struct {
	Element
	Text    string
	Variant string // "primary", "secondary", "danger"
	OnClick func(Event)
}

func (b *ActionButton) WidgetName() widget.Name { return NameActionButton }
func (b *ActionButton) WidgetKind() widget.Kind { return widget.Region }

func (b *ActionButton) Render() *Element {
	var variantCls widget.Class
	switch b.Variant {
	case "secondary":
		variantCls = clsBtnSecondary
	case "danger":
		variantCls = clsBtnDanger
	default:
		variantCls = clsBtnPrimary
	}

	btn := Button().
		Attr("class", string(clsBtn)+" "+string(variantCls)).
		Text(b.Text)

	if b.OnClick != nil {
		btn.On("click", b.OnClick)
	}

	return btn
}
