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
			style.Anchor(),
			style.Row(style.Space2),
			style.ControlBox(),
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
			style.ChipBox(),
			style.CenterContent(),
			style.KeepSize(),
			style.OnEdge(style.EdgeBottom, style.SideEnd, style.SpaceNone, style.Space3),
		).
		// Both the menu and the badge leave the flow, so the label is the only
		// thing sizing the row: every row ends up the same height regardless of
		// how long its title or its badge is. No Anchor() here — Docked already
		// makes this a containing block, and the two fight over `position`.
		Part(PartMenu,
			style.As(style.Subtle),
			style.KeepSize(),
			style.Docked(style.Parent, style.EdgeTop, style.SideEnd, style.Space2),
		).
		Part(PartButton,
			style.Interactive(style.Subtle),
			style.Round(style.RadiusSm),
		).
		// IconBox is not optional: a bare <svg> with no width or height falls
		// back to the replaced-element default of 300x150 and drags the whole
		// row open with it.
		Part(PartIcon,
			style.As(style.Subtle),
			style.IconBox(style.IconMd),
		).
		// Fixed to the viewport, not anchored to the row. Six ancestors between
		// the row and the screen carry overflow — the list's own <ul>, crudview's
		// list and aside, the platform panel and stage, the shell root — and an
		// absolutely positioned panel is clipped by every one of them. Only
		// `fixed` escapes. The backdrop above is what closes it.
		Part(PartOptions,
			style.Stack(style.SpaceNone),
			style.As(style.Panel),
			style.Raise(style.Floating),
			// SideStart, not SideEnd: the trailing corner is where a host's
			// floating action button sits, and the two would land on top of
			// each other.
			style.Docked(style.Viewport, style.EdgeBottom, style.SideStart, style.Space4),
		).
		Part(PartItem,
			style.Interactive(style.Panel),
			style.Pad(style.Space2),
			style.Width(style.Full),
		).
		When(widget.Selected, PartRow,
			style.As(style.Highlight),
		).
		Stylesheet()
}
