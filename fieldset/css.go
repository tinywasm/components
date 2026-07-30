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
			style.Stack(style.SpaceNone),
			style.Pad(style.Space2),
			style.KeepSize(),
		).
		// Width(Content) is what makes this a chip. Inside a Stack the default
		// cross-axis alignment is stretch, so without it the label spans the
		// whole field and turns into a solid bar above the input.
		Part(widget.PartLabel,
			style.As(style.Primary),
			style.Round(style.RadiusSm),
			style.FontSize(style.TextXs),
			style.FontWeight(style.WeightBold),
			style.Raise(style.Raised),
			style.Pad(style.Space1),
			style.Width(style.Content),
		).
		Part(widget.PartError,
			style.As(style.Danger),
			style.FontSize(style.TextXs),
		).
		Part(widget.PartInput,
			style.As(style.Panel),
			style.Round(style.RadiusMd),
			style.Pad(style.Space2),
		).
		Part(widget.PartRadioGroup,
			style.Row(style.Space3),
			style.Pad(style.Space1),
		).
		When(widget.Locked, "",
			style.As(style.Inset),
		).
		When(widget.Invalid, "",
			style.As(style.Danger),
		).
		Stylesheet()
}
