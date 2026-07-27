//go:build !wasm

package selectsearch

import (
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// Style defines the selectsearch visual contract using style DSL.
func (c *SelectSearch) Style() *style.Sheet {
	return style.Of(NameSelectSearch).
		Root(
			style.Stack(style.Space1),
			style.On(style.Panel),
			style.Round(style.RadiusMd),
		).
		Part(PartDropdown,
			style.Stack(style.Space1),
			style.On(style.Panel),
			style.Raise(style.Floating),
			style.Clip(),
		).
		Part(PartHeader,
			style.Row(style.Space2),
			style.Pad(style.Space2),
		).
		Part(PartSearch,
			style.Pad(style.Space2),
			style.On(style.Sunken),
		).
		Part(PartOptions,
			style.Stack(style.Space0),
			style.Scrolls(),
		).
		Part(PartOption,
			style.Row(style.Space2),
			style.Pad(style.Space2),
		).
		Cue(widget.Hover, PartOption,
			style.On(style.PanelHover),
		)
}
