//go:build !wasm

package contentcard

import (
	"github.com/tinywasm/widget/style"
)

// Style defines the visual sheet for contentcard.
func (c *ContentCard) Style() *style.Sheet {
	return style.Of(NameContentCard).
		Root(
			style.Stack(style.Space0),
			style.On(style.Panel),
			style.Round(style.RadiusMd),
			style.Raise(style.Raised),
		).
		Part(PartHeader,
			style.Pad(style.Space2),
			style.FontWeight(style.WeightBold),
		).
		Part(PartBody,
			style.Pad(style.Space2),
			style.Fill(),
		).
		Part(PartFooter,
			style.Pad(style.Space2),
			style.On(style.Sunken),
		)
}
