//go:build !wasm

package statgrid

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the stylesheet for statgrid.
func (s *StatGrid) RenderCSS() *css.Stylesheet {
	return style.For(s).
		Root(
			style.Grid(style.ColumnMedium, style.Space4),
		).
		Part(PartItem,
			style.Stack(style.Space1),
			style.As(style.Panel),
			style.Pad(style.Space4),
			style.Round(style.RadiusMd),
			style.CenterContent(),
		).
		Part(PartValue,
			style.FontSize(style.Text2xl),
			style.FontWeight(style.WeightBold),
		).
		Part(PartLabel,
			style.FontSize(style.TextSm),
		).
		Stylesheet()
}
