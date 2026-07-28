//go:build !wasm

package targetlist

import (
	"github.com/tinywasm/css"
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
		Part(PartLabel,
			style.On(style.Panel),
		).
		Part(PartBadge,
			style.On(style.Sunken),
			style.Round(style.RadiusSm),
			style.Text(style.TextXs),
		).
		Part(PartMenu,
			style.On(style.Panel),
		).
		Part(PartButton,
			style.On(style.Panel),
		).
		Part(PartIcon,
			style.On(style.Panel),
		).
		Part(PartOptions,
			style.Stack(style.Space0),
			style.On(style.Panel),
			style.Raise(style.Floating),
			style.Clip(),
		).
		Part(PartItem,
			style.On(style.Panel),
		).
		When(widget.Selected, PartRow,
			style.On(style.Selected),
		).
		Cue(widget.Hover, PartRow,
			style.On(style.PanelHover),
		)
}

// RenderCSS returns custom positioning rules for targetlist backdrop overlay
// since those overlay/fixed rules are not currently supported by style.Sheet DSL.
func (t *TargetList) RenderCSS() *css.Stylesheet {
	return css.NewStylesheet(
		css.Raw(`
.targetlist__backdrop {
	display: none;
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	z-index: 4;
}
.targetlist:has(.targetlist__menu[open]) .targetlist__backdrop {
	display: block;
}
		`),
	)
}
