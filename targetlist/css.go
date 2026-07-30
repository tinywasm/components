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
		// One line per row: the label takes the free space and pushes the badge
		// and the menu affordance to the trailing edge. Without a flow the <li>
		// falls back to display:list-item and the menu, a block-level <details>,
		// drops onto a second full-width line.
		Part(PartRow,
			style.Row(style.Space2),
			style.Interactive(style.Panel),
			style.Pad(style.Space3),
			style.Round(style.RadiusMd),
		).
		Part(PartLabel,
			style.FontWeight(style.WeightBold),
			style.Grow(),
		).
		// PushEnd because the badge wraps onto its own line under the label:
		// nothing is left beside it to push it, so the free space goes in front.
		Part(PartBadge,
			style.As(style.Inset),
			style.Round(style.RadiusSm),
			style.FontSize(style.TextXs),
			style.PushEnd(),
			style.KeepSize(),
		).
		Part(PartMenu,
			style.As(style.Subtle),
			style.KeepSize(),
		).
		Part(PartButton,
			style.As(style.Subtle),
		).
		// IconBox is not optional: a bare <svg> with no width or height falls
		// back to the replaced-element default of 300x150 and drags the whole
		// row open with it.
		Part(PartIcon,
			style.As(style.Subtle),
			style.IconBox(style.IconMd),
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
