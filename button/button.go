package button

import (
	"github.com/tinywasm/components"
	dom "github.com/tinywasm/components/internal/dom"
)

type Button struct {
	components.BaseComponent
	Text    string
	Variant string // "primary", "secondary", "danger"
	OnClick func(dom.Event)
}

func (b *Button) Render() dom.Node {
	// Default to primary if no variant specified
	variant := b.Variant
	if variant == "" {
		variant = "primary"
	}

	btn := dom.Button().
		ID(b.ID()).
		Class("btn btn-" + variant).
		Text(b.Text)

	if b.OnClick != nil {
		btn.OnClick(b.OnClick)
	}

	return btn.ToNode()
}
