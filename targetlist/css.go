//go:build !wasm

package targetlist

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the targetlist visual contract using the style DSL.
func (t *TargetList) RenderCSS() *css.Stylesheet {
	return style.For(t).
		Root(
			style.Fill(),
			style.Stack(style.SpaceNone),
		).
		Part(PartList,
			style.Stack(style.Space3),
			style.Scroll(),
			style.Pad(style.Space1),
		).
		Part(PartBackdrop,
			style.Backdrop(style.Viewport),
			style.RevealedBy(widget.Open),
		).
		Part(PartRow,
			style.Interactive(style.Panel),
			style.Pad(style.Space3),
			style.Round(style.RadiusMd),
		).
		Part(PartLabel,
			style.As(style.Panel),
		).
		Part(PartBadge,
			style.As(style.Inset),
			style.Round(style.RadiusSm),
			style.FontSize(style.TextXs),
		).
		Part(PartMenu,
			style.As(style.Panel),
		).
		Part(PartButton,
			style.As(style.Panel),
		).
		Part(PartIcon,
			style.As(style.Panel),
		).
		Part(PartOptions,
			style.Stack(style.SpaceNone),
			style.As(style.Panel),
			style.Raise(style.Floating),
			style.HideOverflow(),
		).
		Part(PartItem,
			style.As(style.Panel),
		).
		When(widget.Selected, PartRow,
			style.As(style.Highlight),
		).
		Stylesheet()
}
