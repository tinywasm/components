//go:build !wasm

package fieldset

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the fieldset widget visual contract using the style DSL.
func (f *Fieldset) RenderCSS() *css.Stylesheet {
	return style.For(f).
		// The field is not a card. It is a label chip sitting on an input box:
		// wrapping the pair in a second bordered panel nests a box inside a box
		// and reads as clutter once several fields stack up.
		// Pad, not a gap: the <form> that holds the fields carries no class, so
		// the only place to put air between consecutive fields is inside each
		// one.
		Root(
			style.Anchor(),
			style.Stack(style.SpaceNone),
			style.Pad(style.Space2),
			style.KeepSize(),
		).
		// A legend, not a caption: OnEdge centres the chip ON the input's top
		// border rather than stacking it above, and Space2 is the field's own
		// padding — the distance from the field's border to the input's.
		// No vertical Pad: the badge treatment a list row carries is the height
		// this chip must match (ChipBox alone is 20px — adding Space1 made the
		// legend 26px and taller than every badge beside it). PadInline is
		// exempt: it never touches the height, and text flush against the
		// chip's filled edge reads as a bug, not as density.
		Part(widget.PartLabel,
			style.OnEdge(style.EdgeTop, style.SideStart, style.Space2, style.Space4),
			style.As(style.Primary),
			style.Round(style.RadiusSm),
			style.FontSize(style.TextXs),
			style.FontWeight(style.WeightBold),
			style.Raise(style.Raised),
			style.ChipBox(),
			style.StartContent(),
			style.PadInline(style.Space2),
		).
		// Red text in the field's top-right corner, growing leftward — not a
		// filled bar across the whole field. Docked against the Root's Anchor.
		Part(widget.PartError,
			style.Docked(style.Parent, style.EdgeTop, style.SideEnd, style.Space2),
			style.Glyph(style.Danger),
			style.FontSize(style.TextXs),
		).
		// Space4, not Space2: the legend rides the input's top border and hangs
		// half its height inside the box, so the value needs room to clear it.
		// At the chip's 20px height the hang is 10px and Space4's 16px still
		// clears it.
		Part(widget.PartInput,
			style.As(style.Panel),
			style.Round(style.RadiusMd),
			style.Pad(style.Space4),
			style.ControlBox(),
		).
		Part(widget.PartRadioGroup,
			style.Row(style.Space3),
			style.Pad(style.Space1),
		).
		When(widget.Locked, "",
			style.As(style.Inset),
		).
		// No filled state on the field itself: the red message in the corner is
		// the signal. Painting the whole box danger buried the value.
		When(widget.Invalid, widget.PartInput,
			style.Glyph(style.Danger),
		).
		Stylesheet()
}
