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
		// PadInline(Space8), not the row's own Pad(Space3): the label is the
		// row's only in-flow child, so its text starts flush at the row's own
		// padding edge -- exactly where the leading-edge ⋮ (Docked, IconMd
		// ~24px + Space1 inset) now also sits, since Docked keeps the whole
		// element INSIDE the box rather than straddling the border the way
		// PartBadge's OnEdge does below. Space8 (32px) clears the icon's own
		// footprint with room to spare. The same clearance lands on the
		// trailing edge too -- harmless, the badge lives in a separate
		// vertical zone below the text line.
		Part(PartLabel,
			style.FontWeight(style.WeightBold),
			style.Grow(),
			style.PadInline(style.Space8),
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
		// Both the menu and the badge leave the flow, so the label is the only
		// thing sizing the row: every row ends up the same height regardless of
		// how long its title or its badge is. No Anchor() here — Docked already
		// makes this a containing block, and the two fight over `position`.
		//
		// Leading edge, not trailing: on the mobile master-detail strip
		// (rightpanel's MasterDetail(Most)) selecting a row navigates to the
		// form panel and leaves only a 10% sliver of the list showing — the
		// row's LEADING edge, by construction (Size.Most: "leaves a sliver of
		// what sits behind it"). A trailing-edge menu is the one part of the
		// row that sliver can never show, which strands the only control that
		// unlocks the now-read-only form on the panel the user just left. Do
		// not move this back to the trailing edge without re-solving that.
		// Space1, not Space2: the sliver is ~37.5px at a 375px viewport and
		// every pixel of inset is budgeted (see layout/docs/PLAN.md).
		Part(PartMenu,
			style.As(style.Subtle),
			style.KeepSize(),
			style.Docked(style.Parent, style.EdgeTop, style.SideStart, style.Space1),
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
		// The menu hangs off the row it belongs to on desktop: the same
		// gesture should put it in the same place. SideStart, matching
		// PartMenu's own leading-edge trigger above — a SideEnd flyout hanging
		// off a SideStart trigger would open backwards, away from the finger
		// that opened it.
		//
		// As(Panel)/Raise/HideOverflow are Tablet+Desktop only, not base: on
		// mobile the two options stop being flush rows inside one shared card
		// and become their own independent floating icon buttons (see PartItem/
		// PartItemDanger below) — a bordered, shadowed panel wrapped around
		// two things that already float on their own would double the chrome.
		// Mirrors the same base/device split rightpanel already uses for its
		// aside and article cards.
		Part(PartOptions,
			style.Stack(style.SpaceNone),
			style.Flyout(style.SideStart),
		).
		On(css.Tablet, PartOptions,
			style.As(style.Panel),
			style.Raise(style.Floating),
			style.HideOverflow(),
		).
		On(css.Desktop, PartOptions,
			style.As(style.Panel),
			style.Raise(style.Floating),
			style.HideOverflow(),
		).
		// On mobile the row-anchored Flyout above is measured from the
		// trigger's own position — which swings between ~10px and ~370px into
		// a 375px viewport depending on which side of the master-detail strip
		// is scrolled into view (the row's own leading edge, wherever that
		// currently sits on screen). A panel anchored there can overflow either
		// edge. Docked to the viewport corner instead, it always has the full
		// screen to render in regardless of where the trigger is.
		// PartBackdrop (Backdrop(Viewport) + RevealedBy(Open)) already
		// dismisses it on an outside tap; this does not need a second one.
		// Space2 gap: the two options are now independent floating squares,
		// not flush rows in a card, so the stack needs real air between them.
		On(css.Mobile, PartOptions,
			style.Stack(style.Space2),
			style.Docked(style.Viewport, style.EdgeBottom, style.SideStart, style.Space4),
		).
		// Square: the items are flush rows inside the panel, not buttons floating
		// in it. An explicit Round overrides the radius As(Panel) would default
		// to; the panel's own HideOverflow is what rounds the outer corners.
		// (Desktop/Tablet only in effect — see the mobile override below.)
		Part(PartItem,
			style.Interactive(style.Panel),
			style.Round(style.RadiusNone),
			style.Pad(style.Space2),
			style.Width(style.Full),
		).
		// Danger, not Panel: mirrors crudview's own delconfirm-btn/
		// delconfirm-btn-danger split — same base shape, Eliminar alone tinted
		// for a destructive action. Desktop/Tablet keep this identical to
		// PartItem apart from color; see the mobile override below for why.
		Part(PartItemDanger,
			style.Interactive(style.Panel),
			style.Round(style.RadiusNone),
			style.Pad(style.Space2),
			style.Width(style.Full),
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
		).
		Stylesheet()
}
