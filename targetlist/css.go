//go:build !wasm

package targetlist

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the targetlist visual contract using the style DSL.
func (t *TargetList) RenderCSS() *css.Stylesheet {
	return t.sheet().Stylesheet()
}

// sheet builds the style Sheet. Split from RenderCSS so Validate() can check
// the declared part tree — the containment Without the StyleSheet round-trip.
func (t *TargetList) sheet() *style.Sheet {
	return style.For(t).
		Root(
			style.Fill(),
			style.Stack(style.SpaceNone),
		).
		// The list's own gutter is INLINE only. Scroll() already reserves the
		// block edges through var(--floating-top|bottom, 0px) — the strip a
		// FloatingChrome host declares over this list — and a `padding:`
		// shorthand here would silently clobber that reservation (the widgets
		// layer beats the seam's primitives layer), putting the last row back
		// under the host's floating button. No block Pad means the block
		// gutter is the seam's business now: flush by default, the host's
		// strip when one is declared.
		Part(PartList,
			style.Stack(style.Space3),
			style.Scroll(),
			style.PadInline(style.Space1),
		).
		// Flush on mobile, not Space1: this is INSIDE the crudview card
		// (between its Inset border and the row content), unlike rightpanel's
		// aside gutter which stays put on every breakpoint specifically so
		// that outer border keeps a buffer against the screen edge. Reclaiming
		// space here instead of there is what makes the ⋮ menu's clearance in
		// the mobile sliver not cost the card its own visible frame — see
		// rightpanel/css.go's partAside comment for the two sides of this
		// trade.
		//
		// Inline only, like the base rule: the block edges belong to the seam,
		// so the flush reclaims the sliver without clobbering the bottom
		// reservation of a FloatingChrome host.
		On(css.Mobile, PartList,
			style.PadInline(style.SpaceNone),
		).
		// The row wraps, and KeepSize is what lets it GROW. Row(Space2) already
		// sets flex-wrap, which is what puts the open options on their own line
		// under the trigger. But this <li> is a flex item of the list's own
		// column, and a flex item defaults to flex-shrink: 1 — inside a
		// Scroll() column whose height is constrained, that shrinks every row
		// back to its min-height and the options render OUTSIDE the row's box.
		// Measured: the row stayed at 50px while its own options painted from
		// 228 to 269.6, past a bottom edge of 225.2. With KeepSize the row goes
		// to 107.2px and contains them.
		Part(PartRow,
			style.Anchor(),
			style.Row(style.Space2),
			style.KeepSize(),
			style.ControlBox(),
			style.Interactive(style.Panel),
			style.Pad(style.Space3),
			style.Round(style.RadiusMd),
		).
		// No PadInline: the label is the row's second in-flow child, after the
		// menu, so the row's own Row(Space2) gap is what clears the leading-edge
		// ⋮ (IconMd, ~24px). The PadInline(Space8) this part used to carry
		// existed only to dodge the icon when PartMenu was Docked on top of it —
		// with the trigger back in the flow there is nothing to dodge.
		Part(PartLabel,
			style.FontWeight(style.WeightBold),
			style.Grow(),
		).
		// PushEnd because the badge wraps onto its own line under the label:
		// nothing is left beside it to push it, so the free space goes in front.
		//
		// The straddle is half a --chip-height of negative bottom margin — the
		// token both this badge and a fieldset legend share — never a
		// transform: the badge's box now EXISTS for the list's scrollHeight,
		// so a host that floats an action button over this list reserves the
		// strip with FloatingChrome and Scroll() pads by
		// var(--floating-bottom, 0px). Do not swap the margin back for a
		// transform: that is what made the badge's real position invisible to
		// every scroll-size calculation.
		Part(PartBadge,
			style.As(style.Inset),
			style.Round(style.RadiusSm),
			style.FontSize(style.TextXs),
			style.ChipBox(),
			style.CenterContent(),
			style.KeepSize(),
			style.OnEdge(style.EdgeBottom, style.SideEnd, style.SpaceNone, style.Space3),
		).
		// The trigger in flow, the row's FIRST child: as a flex item it sits at
		// the row's leading edge by order alone — no Docked, no Anchor, nothing
		// positioned anywhere on this row's menu chain.
		//
		// Leading edge, not trailing: on the mobile master-detail strip
		// (rightpanel's MasterDetail(Most)) selecting a row navigates to the
		// form panel and leaves only a 10% sliver of the list showing — the
		// row's LEADING edge, by construction (Size.Most: "leaves a sliver of
		// what sits behind it"). A trailing-edge menu is the one part of the
		// row that sliver can never show, which strands the only control that
		// unlocks the now-read-only form on the panel the user just left. Do
		// not move this back to the trailing edge without re-solving that.
		//
		// In flow the icon starts at the row's own Pad(Space3) = 12px instead
		// of the old Docked Space1 = 4px inset. Measured in the browser: the
		// ~37.5px sliver at a 375px viewport still fits the 24px icon with
		// 1.5px to spare (see layout/docs/PLAN.md for the budget).
		Part(PartButton,
			style.As(style.Subtle),
			style.KeepSize(),
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
		// An accordion INSIDE the row, not an overlay hanging off it. This is
		// the second time this part's positioning caused a shipped bug, and
		// both had the same root: an out-of-flow panel cannot coexist with the
		// list's own Scroll().
		//
		//   1. As a Flyout it resolved `inset-block-start: 100%` against the
		//      nearest POSITIONED ancestor. With the trigger Docked that was
		//      the 24px <details>, so the panel opened over its own row's
		//      label (measured: 21.2px inside the row, 8.2px of text covered).
		//   2. Anchored correctly to the row instead, it then hit the OTHER
		//      wall: .targetlist__list is a Scroll() region, and an absolutely
		//      positioned descendant is clipped by it. On the last row only
		//      10px of an 84.8px panel survived.
		//
		// The mobile Docked(Viewport, …) that used to be here was the escape
		// hatch from (2) — it left the clipper by leaving the row entirely, and
		// landed 502px away in a screen corner, on top of two unrelated rows.
		// That detachment is what then needed a Veil()'d Backdrop to explain
		// itself, and the veil blurred the very row the user was acting on.
		// A patch on a patch on a patch.
		//
		// In flow there is no clipper to escape: the content grows the row,
		// the row grows the scroller, and the browser can scroll to it. Same
		// last row, measured after: 41.6px of 41.6px visible. The target stays
		// on screen, wearing its own Selected amber, with its actions inside
		// it — so nothing has to be dimmed to explain what acts on what.
		//
		// Width(Full) + KeepSize is what puts it on its own line: Row(Space2)
		// on PartRow already wraps, and a flex item with basis 100% that
		// cannot shrink has to take a line of its own. RevealedBy(Open) is the
		// open/close — a per-row signal, since the options are no longer a
		// child of a native <details> that could hide them itself.
		Part(PartOptions,
			style.Row(style.Space2),
			style.Width(style.Full),
			style.KeepSize(),
			style.RevealedBy(widget.Open),
		).
		// Square: the items are flush rows inside the panel, not buttons floating
		// Content-width buttons side by side, not full-width flush rows: the
		// options are a line INSIDE the row now, not a card of their own, so
		// two stretched bars would read as a second list rather than as this
		// row's two actions. Rounded for the same reason — RadiusNone only
		// made sense when a panel's HideOverflow was clipping the outer
		// corners, and there is no panel left.
		Part(PartItem,
			style.Interactive(style.Panel),
			style.Round(style.RadiusSm),
			style.Pad(style.Space2),
			style.Width(style.Content),
		).
		// Danger, not Panel: mirrors crudview's own delconfirm-btn/
		// delconfirm-btn-danger split — same base shape, Eliminar alone tinted
		// for a destructive action. Desktop/Tablet keep this identical to
		// PartItem apart from color; see the mobile override below for why.
		Part(PartItemDanger,
			style.Interactive(style.Danger),
			style.Round(style.RadiusSm),
			style.Pad(style.Space2),
			style.Width(style.Content),
		).
		// Mobile: both options become floating icon buttons matching
		// crudview's own action button exactly (As, Round, Pad, Raise(Floating),
		// CenterContent — same recipe, see crudview/css.go's mobile "action"
		// override) instead of a text dropdown row. Primary for Editar, the
		// chassis' default "this is a control" color; Danger for Eliminar,
		// matching the destructive-action color already used elsewhere in this
		// same component set.
		On(css.Mobile, PartItem,
			style.As(style.Primary),
			style.Round(style.RadiusMd),
			style.Pad(style.Space2),
			style.Width(style.Content),
			style.Raise(style.Floating),
			style.CenterContent(),
		).
		On(css.Mobile, PartItemDanger,
			style.As(style.Danger),
			style.Round(style.RadiusMd),
			style.Pad(style.Space2),
			style.Width(style.Content),
			style.Raise(style.Floating),
			style.CenterContent(),
		).
		// OnlyOn: the icon exists ONLY on mobile, where it is the entire visible
		// content of the button — desktop shows the text label instead (see
		// PartItemLabel below) and never renders this icon at all. IconLg
		// matches crudview's own mobile action-new/action-cancel icon sizing.
		// CenterContent is not decoration here: OnlyOn hides the part by
		// default (display: none) and only a flow primitive on the device rule
		// gives it a real display to reveal into — IconBox alone sizes a box
		// that would stay display:none forever.
		OnlyOn(css.Mobile, PartItemIcon,
			style.CenterContent(),
			style.IconBox(style.IconLg),
		).
		// The label is the visible content everywhere except mobile, where the
		// icon above takes over and the text would just double up inside an
		// already-labelled icon button. No base Part() call: plain text needs
		// no styling of its own, and an empty Part emits nothing at all, which
		// the sheet rejects as pointless — the On(Mobile, Hide()) below is the
		// only rule this class ever needs.
		On(css.Mobile, PartItemLabel,
			style.Hide(),
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
		)
}
