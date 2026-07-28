//go:build !wasm

package fieldset

import (
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// Style defines the fieldset widget visual contract using the style DSL.
func (f *Fieldset) Style() *style.Sheet {
	return style.Of(widget.NameField).
		Root(
			style.On(style.Panel),
			style.Round(style.RadiusMd),
			style.Pad(style.Space3),
			style.Stack(style.Space2),
		).
		Part(widget.PartLabel,
			style.On(style.Accent),
			style.Round(style.RadiusSm),
			style.Text(style.TextXs),
			style.FontWeight(style.WeightBold),
			style.Raise(style.Raised),
		).
		Part(widget.PartError,
			style.On(style.Danger),
			style.Text(style.TextXs),
		).
		Part(widget.PartInput,
			style.On(style.Panel),
		).
		Part(widget.PartRadioGroup,
			style.Row(style.Space3),
			style.Pad(style.Space1),
		).
		When(widget.Locked, "",
			style.On(style.Sunken),
		).
		When(widget.Invalid, "",
			style.On(style.Danger),
		)
}
