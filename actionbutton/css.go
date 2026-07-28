//go:build !wasm

package actionbutton

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the actionbutton visual contract using the style DSL.
func (b *ActionButton) RenderCSS() *css.Stylesheet {
	return style.For(b).
		Root(
			style.Pad(style.Space2),
			style.Round(style.RadiusSm),
			style.As(style.Page),
		).
		Part(PartPrimary,
			style.Interactive(style.Primary),
		).
		Part(PartSecondary,
			style.Interactive(style.Secondary),
		).
		Part(PartDanger,
			style.Interactive(style.Danger),
		).
		Stylesheet()
}
