//go:build !wasm

package countbadge

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the countbadge visual contract using the style DSL.
//
// The bubble is OnEdge(EdgeTop, SideEnd): centred on the host's top-end
// corner, half outside and half inside — the notification-bubble placement.
// OnEdge is absolute, so the host's box never depends on the count; that is
// the whole point of this component (an in-flow counter stretches its
// button the moment it appears).
//
// As(AccentInverse): amber on a Primary button, the same language the
// toggle button speaks while open. No ChipBox: a bubble shrink-wraps its
// own figure — "7" is round, "128" is a pill — instead of lining up with a
// column that does not exist here.
func (b *CountBadge) RenderCSS() *css.Stylesheet {
	return style.For(b).
		Part(PartBadge,
			style.OnEdge(style.EdgeTop, style.SideEnd, style.SpaceNone, style.SpaceNone),
			style.As(style.AccentInverse),
			style.Round(style.RadiusFull),
			style.FontSize(style.TextXs),
			style.FontWeight(style.WeightBold),
			style.PadInline(style.Space1),
			style.CenterContent(),
			style.RevealedBy(widget.Open),
		).
		Stylesheet()
}
