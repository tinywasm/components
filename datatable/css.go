//go:build !wasm

package datatable

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the datatable visual contract using the style DSL.
func (t *DataTable) RenderCSS() *css.Stylesheet {
	return style.Of(NameDataTable).
		Root(
			style.Width(style.Full),
			style.On(style.Panel),
		).
		Part(PartHeader,
			style.On(style.Sunken),
			style.FontWeight(style.WeightBold),
			style.Pad(style.Space2),
		).
		Part(PartRow,
			style.Pad(style.Space2),
		).
		Cue(widget.Hover, PartRow,
			style.On(style.PanelHover),
		).
		Stylesheet()
}
