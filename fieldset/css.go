//go:build !wasm

package fieldset

import (
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// Style defines the fieldset widget visual contract using the style DSL.
func (f *Fieldset) Style() *style.Sheet {
	return style.Of(NameFieldset).
		Root(
			style.On(style.Panel),
			style.Round(style.RadiusMd),
			style.Pad(style.Space3),
			style.Stack(style.Space2),
		).
		Part(PartLabel,
			style.On(style.Accent),
			style.Round(style.RadiusSm),
			style.Text(style.TextXs),
			style.FontWeight(style.WeightBold),
			style.Raise(style.Raised),
		).
		Part(PartError,
			style.On(style.Danger),
			style.Text(style.TextXs),
		).
		Part(PartInput,
			style.On(style.Panel),
		).
		When(widget.Locked, "",
			style.On(style.Sunken),
		).
		When(widget.Invalid, "",
			style.On(style.Danger),
		)
}
