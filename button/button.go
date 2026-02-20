package button

import (
	"github.com/tinywasm/dom"
)

type Button struct {
	*dom.Element
	Text    string
	Variant string // "primary", "secondary", "danger"
	OnClick func(dom.Event)
}

func (b *Button) Render() *dom.Element {
	// Default to primary if no variant specified
	variant := b.Variant
	if variant == "" {
		variant = "primary"
	}

	if b.Element == nil {
		b.Element = &dom.Element{}
	}

	btn := dom.Button().
		Class("btn btn-" + variant).
		Text(b.Text)

	if b.OnClick != nil {
		btn.On("click", b.OnClick)
	}

	return btn
}
