//go:build !wasm

package targetlist

import (
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// Style defines the targetlist visual contract using the style DSL.
func (t *TargetList) Style() *style.Sheet {
	return style.Of(NameTargetList).
		Root(
			style.Fill(),
			style.Stack(style.Space0),
		).
		Part(PartList,
			style.Stack(style.Space3),
			style.Scrolls(),
			style.Pad(style.Space1),
		).
		Part(PartBackdrop,
			style.On(style.Sunken),
			style.Fixed(),
		).
		Part(PartRow,
			style.Row(style.Space2),
			style.On(style.Panel),
			style.Pad(style.Space3),
			style.Round(style.RadiusMd),
		).
		Part(PartBadge,
			style.On(style.Sunken),
			style.Round(style.RadiusSm),
			style.Text(style.TextXs),
		).
		Part(PartMenu,
			style.On(style.Panel),
		).
		Part(PartOptions,
			style.Stack(style.Space0),
			style.On(style.Panel),
			style.Raise(style.Floating),
			style.Clip(),
		).
		When(widget.Selected, PartRow,
			style.On(style.Selected),
		).
		Cue(widget.Hover, PartRow,
			style.On(style.PanelHover),
		)
}
