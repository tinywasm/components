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
		Part(PartToggle,
			style.As(style.Panel),
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
		Part(PartHeader,
			style.Row(style.Space2),
			style.Pad(style.Space2),
		).
		Part(PartIcon,
			style.As(style.Panel),
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
		Part(PartLabel,
			style.As(style.Panel),
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
