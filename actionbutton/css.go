//go:build !wasm

package actionbutton

import (
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// Style defines the actionbutton visual contract using the style DSL.
func (b *ActionButton) Style() *style.Sheet {
	return style.Of(NameActionButton).
		Root(
			style.Pad(style.Space2),
			style.Round(style.RadiusSm),
			style.On(style.Page),
		).
		Part(PartPrimary,
			style.On(style.Accent),
		).
		Part(PartSecondary,
			style.On(style.Secondary),
		).
		Part(PartDanger,
			style.On(style.Danger),
		).
		Cue(widget.Hover, PartPrimary,
			style.On(style.AccentHover),
		).
		Cue(widget.Hover, PartSecondary,
			style.On(style.SecondaryHover),
		).
		Cue(widget.Hover, PartDanger,
			style.On(style.DangerHover),
		)
}
