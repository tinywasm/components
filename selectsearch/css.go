//go:build !wasm

package selectsearch

import (
	"github.com/tinywasm/components/listgap"
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the selectsearch visual contract using the style DSL.
func (c *SelectSearch) RenderCSS() *css.Stylesheet {
	return c.sheet().Stylesheet()
}

// sheet is split out from RenderCSS so tests can call Validate() on the
// *style.Sheet directly, without the *css.Stylesheet conversion — the same
// split targetlist/css.go uses (see targetlist_test.go's TestSheetValidates).
func (c *SelectSearch) sheet() *style.Sheet {
	s := style.For(c).
		// Anchor: what lets PartDropdown's Flyout (added below, tablet/desktop
		// only) hang from THIS box instead of some positioned ancestor further
		// up the tree. Mirrors usermenu's Root(Anchor()) exactly — same
		// primitive, same reason ("a corner pin needs a positioned ancestor to
		// hang from").
		Root(
			style.Anchor(),
			style.Stack(style.Space1),
			style.As(style.Panel),
			style.Round(style.RadiusMd),
		)
	// The same list container targetlist and targetdate assemble — one lego
	// piece, so the dropdown's rows breathe with the same rhythm as every
	// other list in the app instead of inventing a third spacing.
	listgap.Apply(s, PartOptions)
	s.On(css.Mobile, PartOptions, listgap.MobileOpts()...)
	return s.
		// The checkbox drives PartDropdown's open state via the label's `for`
		// attribute (see selectsearch.go) — it is the CSS-only toggle
		// mechanism, never meant to be seen. A label still activates a
		// display:none checkbox natively, so Hide() loses nothing.
		Part(PartToggle,
			style.Hide(),
		).
		// Mobile-first base: the dropdown is a plain IN-FLOW block here — no
		// Raise, no Flyout. Same choice usermenu's PartPanel makes for a phone:
		// there is no room to spare for an anchored floating overlay on a
		// narrow viewport, so opening the control expands it in place instead
		// (an accordion, not a popover). Raise(Floating)+Flyout are added ONLY
		// from Tablet up, in the On(...) blocks below — do not add them here.
		// Pad(Space1): the dropdown is a rounded, clipping panel
		// (As(Panel) brings RadiusMd by default; HideOverflow clips to it), so
		// a child flush against its top edge loses its corners — and with them
		// the inset focus ring, which is what made the search box look sawn
		// off. The padding is what keeps every child inside the rounded area.
		Part(PartDropdown,
			style.Stack(style.Space1),
			style.As(style.Panel),
			style.Pad(style.Space1),
			style.HideOverflow(),
		).
		// SpaceNone, not Space2: the gap between the icon and the text comes
		// from PartHeaderText's own left padding below, same as searchbar's
		// root uses a zero gap and lets PartInput's own Pad(Space2) make the
		// space — a Row gap here would ADD to that, pushing the icon away
		// from the header's own edge instead of leaving it flush.
		// Round+HideOverflow (not on Root, where the floating PartDropdown
		// lives — clipping there would cut the dropdown off) is what makes
		// PartIcon's own square corners read as the header's rounded ones:
		// the icon carries no radius of its own, the header clips it to
		// match.
		// ControlBox+KeepSize on BOTH the bar and its cap, exactly as
		// searchbar/css.go declares them: Row() carries flex-wrap: wrap, so a
		// narrow viewport wrapped the text onto its own line under the square
		// and the header stopped being a bar. KeepSize is what forbids the
		// wrap; ControlBox is what keeps the two halves the same height.
		Part(PartHeader,
			style.Row(style.SpaceNone),
			style.Round(style.RadiusMd),
			style.HideOverflow(),
			style.ControlBox(),
			style.KeepSize(),
		).
		// The padded text — see PartHeader above for why the padding lives
		// here and not on the header itself.
		Part(PartHeaderText,
			style.Grow(),
			style.Pad(style.Space2),
		).
		// The filled square cap around the arrow, FIRST child of the header
		// (see selectsearch.go) — searchbar's own PartIcon recipe
		// (MediaBox(AspectSquare) sized off the same --control-height
		// ControlBox answers to, so the cap is a true square), As(Primary)
		// so it reads as part of the same chrome the rest of the chassis
		// uses instead of a bare mark floating in the header. No Round of
		// its own — PartHeader's clip above is what rounds its visible
		// corners. KeepSize, as on the header: a square cannot wrap into a
		// sliver when the row next to it runs out of space.
		Part(PartIcon,
			style.As(style.Primary),
			style.MediaBox(style.AspectSquare),
			style.ControlBox(),
			style.KeepSize(),
		).
		// A bare <svg> with no box falls back to 300x150 — same gotcha
		// searchbar/css.go's PartGlyph documents; IconBox pins it.
		// The GLYPH turns, not PartIcon: rotating the cap would spin the whole
		// filled square. TurnNone is the resting rule Animate() transitions
		// from — without a base value there is no start state to move off.
		Part(PartGlyph,
			style.IconBox(style.IconSm),
			style.Rotate(style.TurnNone),
			style.Animate(style.MotionBase),
		).
		// A radius of its own so the box reads as a control, not as a slab
		// filling the panel's width. ControlBox answers to the same
		// --control-height token the header uses, so the search field and the
		// cap stay the same height.
		Part(PartSearch,
			style.Pad(style.Space2),
			style.As(style.Inset),
			style.Round(style.RadiusSm),
			style.ControlBox(),
		).
		// The row recipe targetlist/css.go declares for PartRow: same surface,
		// same box, same radius. A dropdown option and a list row are the same
		// object to the user; they must not be two different shapes.
		Part(PartOption,
			style.Row(style.Space2),
			style.KeepSize(),
			style.ControlBox(),
			style.Interactive(style.Page),
			style.Pad(style.Space3),
			style.Round(style.RadiusMd),
		).
		// The Label/Sublabel column: Grow() takes the row's free width so
		// PartDesc's badge lands at the trailing edge instead of hugging
		// straight after the name — same reasoning as targetlist's own
		// PartLabel.
		Part(PartText,
			style.Stack(style.SpaceNone),
			style.Grow(),
		).
		Part(PartLabel,
			style.FontWeight(style.WeightBold),
		).
		// Muted and small — a second line under the name, not competing with
		// it for attention. Glyph, not As(Inset): a tinted color reads as
		// secondary text, a filled box reads as its own control — a plain
		// two-line name+id has no control to be.
		Part(PartSublabel,
			style.Glyph(style.Subtle),
			style.FontSize(style.TextXs),
		).
		Part(PartDesc,
			style.As(style.Inset),
			style.FontSize(style.TextXs),
		).
		// The chevron IS the open state. WhenWithin from the root, because
		// data-open lives on the root (see selectsearch.go's BindState) and the
		// glyph is two levels down inside the header.
		WhenWithin(widget.Open, "", PartGlyph,
			style.Rotate(style.TurnHalf),
		).
		// Amber is the chassis' one "where I am" statement — the rail's current
		// nav item and a selected list row wear it too. AccentWash on hover and
		// focus reads as "on the way to selected"; a grey mix reads as chrome
		// with no relation to what clicking does. Focus repeats Hover on
		// purpose: a keyboard user gets no :hover and must not be stranded on a
		// different colour.
		When(widget.Selected, PartOption, style.As(style.Accent)).
		Cue(widget.Hover, PartOption, style.As(style.AccentWash)).
		Cue(widget.Focus, PartOption, style.As(style.AccentWash)).
		Cue(widget.Press, PartOption, style.As(style.Accent)).
		// From a tablet up, the dropdown becomes a floating panel hanging off
		// the Root anchor — the identical escalation usermenu's PartPanel uses,
		// gated on the identical two breakpoints, for the identical reason.
		On(css.Tablet, PartDropdown,
			style.Raise(style.Floating),
			style.Flyout(style.SideStart),
		).
		On(css.Desktop, PartDropdown,
			style.Raise(style.Floating),
			style.Flyout(style.SideStart),
		)
}
