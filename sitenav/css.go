//go:build !wasm

package sitenav

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the visual stylesheet for sitenav.
func (sn *SiteNav) RenderCSS() *css.Stylesheet {
	return style.For(sn).
		Root(
			style.Row(style.Space4),
			style.As(style.Page),
			style.Pad(style.Space3),
			style.CenterContent(),
		).
		Part(PartBrand,
			style.KeepSize(),
			style.CenterContent(),
		).
		// The logo ships as whatever the site declares (WideLogoSrc/
		// CompactLogoSrc) — often an SVG whose viewBox has nothing to do with
		// a nav bar. KeepSize() alone only stops flex from stretching or
		// shrinking it; it does not cap its size, so a bare <img> renders at
		// its own intrinsic dimensions — a full-page illustration traced at
		// 918x478 becomes the entire header. LogoBox() is exactly this: cap
		// height to the row's control-height, width auto to keep the aspect
		// ratio (IconBox forces width=height, wrong for a non-square mark).
		Part(PartLogo,
			style.LogoBox(),
		).
		Part(PartToggle,
			style.IconBox(style.IconSm),
			style.As(style.Panel),
			style.KeepSize(),
		).
		// The hamburger only exists where the links do not fit. Leaving it on
		// every viewport put a stray button next to a nav bar that was already
		// showing all six links.
		On(css.Tablet, PartToggle, style.Hide()).
		On(css.Desktop, PartToggle, style.Hide()).
		// The mirror image, and the reason RevealedBy sits inside On(): on a
		// phone the links collapse behind the toggle and come back with the
		// Open state the button writes; on anything wider they simply stay in
		// the row. Declared at the top level it would hide them on every
		// viewport. Before this, nothing was hidden at all — the six labels
		// stayed in the row and overflowed a narrow screen horizontally, and
		// the toggle was rendered, wired and useless.
		// Stack, not Row: revealed on a phone the links have the full width of
		// the bar and one line each. Kept as a row they ran off the side of
		// the screen — the same overflow the collapse exists to prevent, only
		// now behind a tap.
		On(css.Mobile, PartMenu, style.Stack(style.Space2), style.RevealedBy(widget.Open)).
		On(css.Mobile, PartNav, style.Stack(style.Space2)).
		Part(PartMenu,
			style.Row(style.Space4),
			style.Grow(),
			style.CenterContent(),
		).
		Part(PartNav,
			style.Row(style.Space3),
			style.Grow(),
			style.CenterContent(),
		).
		Part(PartLink,
			style.Pad(style.Space1),
			style.FontSize(style.TextSm),
			style.FontWeight(style.WeightMedium),
		).
		Part(PartActions,
			style.Row(style.Space2),
			style.KeepSize(),
		).
		Stylesheet()
}
