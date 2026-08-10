//go:build !wasm

package selectsearch

import (
	"github.com/tinywasm/css"
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
	return style.For(c).
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
		).
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
		Part(PartDropdown,
			style.Stack(style.Space1),
			style.As(style.Panel),
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
		Part(PartHeader,
			style.Row(style.SpaceNone),
			style.Round(style.RadiusMd),
			style.HideOverflow(),
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
		// corners.
		Part(PartIcon,
			style.As(style.Primary),
			style.MediaBox(style.AspectSquare),
			style.ControlBox(),
		).
		// A bare <svg> with no box falls back to 300x150 — same gotcha
		// searchbar/css.go's PartGlyph documents; IconBox pins it.
		Part(PartGlyph,
			style.IconBox(style.IconSm),
		).
		Part(PartSearch,
			style.Pad(style.Space2),
			style.As(style.Inset),
		).
		Part(PartOptions,
			style.Stack(style.SpaceNone),
			style.Scroll(),
		).
		Part(PartOption,
			style.Row(style.Space2),
			style.Pad(style.Space2),
			style.Interactive(style.Panel),
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
