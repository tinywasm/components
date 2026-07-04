//go:build !wasm

package themetoggle

import . "github.com/tinywasm/css"

func (t *ThemeToggle) RenderCSS() *Stylesheet {
	// The [data-theme="light"|"dark"] palette overrides live in tinywasm/css's
	// RenderCSS (always bundled), so the manual toggle works even if this
	// component's own CSS isn't collected. Here we only style the button.
	return NewStylesheet(
		Rule(clsTsBtn,
			Position(Str("fixed")),
			Top(Rem(1)),
			Right(Rem(1)),
			ZIndex(Str("9999")),
			Width(Rem(2.4)),
			Height(Rem(2.4)),
			Padding(Zero),
			BorderRadius(Pct(50)),
			Border(None),
			Cursor(Pointer),
			FontSize(Rem(1.1)),
			LineHeight(Str("1")),
			Display(Flex_),
			AlignItems(Center),
			JustifyContent(Center),
			Background(ColorPrimary),
			Color(ColorOnPrimary),
			Opacity(0.85),
			Transition(Str("opacity 0.2s, transform 0.2s")),
		),
		Rule(clsTsBtn.Hover(),
			Opacity(1),
			Transform(Str("scale(1.12)")),
			Background(ColorSecondary),
			Color(ColorOnSecondary),
		),
	)
}
