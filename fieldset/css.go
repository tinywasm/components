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
		Root(
			style.As(style.Panel),
			style.Round(style.RadiusMd),
			style.Pad(style.Space3),
			style.Stack(style.Space2),
		).
		Part(widget.PartLabel,
			style.As(style.Primary),
			style.Round(style.RadiusSm),
			style.FontSize(style.TextXs),
			style.FontWeight(style.WeightBold),
			style.Raise(style.Raised),
		).
		Part(widget.PartError,
			style.As(style.Danger),
			style.FontSize(style.TextXs),
		).
		Part(widget.PartInput,
			style.As(style.Panel),
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
