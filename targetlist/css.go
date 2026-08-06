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
		// The menu hangs off the row it belongs to, on every screen: the same
		// gesture should put it in the same place. Measured, it is not clipped
		// by the list's scrollport at either width.
		Part(PartOptions,
			style.Stack(style.SpaceNone),
			style.As(style.Panel),
			style.Raise(style.Floating),
			style.HideOverflow(),
			style.Flyout(style.SideEnd),
		).
		// Square: the items are flush rows inside the panel, not buttons floating
		// in it. An explicit Round overrides the radius As(Panel) would default
		// to; the panel's own HideOverflow is what rounds the outer corners.
		Part(PartItem,
			style.Interactive(style.Panel),
			style.Round(style.RadiusNone),
			style.Pad(style.Space2),
			style.Width(style.Full),
		).
		// Accent, not Highlight: the selection tint is a 15% blue wash that
		// reads as "disabled" on a light panel. The amber fill is the one
		// "where I am" statement the whole chassis shares — the rail's current
		// nav item wears the same surface.
		When(widget.Selected, PartRow,
			style.As(style.Accent),
		).
		// AccentWash, not Interactive(Panel)'s own grey mix: a hover that
		// leans toward the same amber the selected state commits to reads as
		// "on the way to selected" -- a grey hover reads as unrelated chrome
		// with no connection to what clicking it does.
		Cue(widget.Hover, PartRow,
			style.As(style.AccentWash),
		).
		// Same treatment as Hover above, on purpose: a keyboard user tabbing
		// through the list gets no :hover at all, so leaving Focus on
		// Interactive(Panel)'s default grey would give mouse/touch users the
		// amber preview and strand keyboard users on the old unrelated color
		// -- the exact inconsistency pairing Hover with Focus exists to rule
		// out.
		Cue(widget.Focus, PartRow,
			style.As(style.AccentWash),
		).
		// Same Accent as When(Selected) above, not Interactive(Panel)'s own
		// grey press mix: a tap is :active for the whole time the finger is
		// down, which outlasts the click handler that sets data-selected —
		// same specificity, same layer, :active declared later, so it wins
		// for that whole window. A grey Press painted during exactly that
		// window is what read as "select, flash grey, THEN turn amber" on a
		// phone. Matching Press to Selected's own color makes the two
		// indistinguishable, so there is nothing left to flash.
		Cue(widget.Press, PartRow,
			style.As(style.Accent),
		).
		Stylesheet()
}
