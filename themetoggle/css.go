//go:build !wasm

package themetoggle

import . "github.com/tinywasm/css"

func (t *ThemeToggle) RenderCSS() *Stylesheet {
	return NewStylesheet(
		Rule(Selector("[data-theme=\"light\"]"),
			Bind(ColorBackground, ColorBackgroundLight),
			Bind(ColorSurface, ColorSurfaceLight),
			Bind(ColorOnSurface, ColorOnSurfaceLight),
			Bind(ColorMuted, ColorMutedLight),
			Bind(ColorHover, ColorHoverLight),
		),
		Rule(Selector("[data-theme=\"dark\"]"),
			Bind(ColorBackground, ColorBackgroundDark),
			Bind(ColorSurface, ColorSurfaceDark),
			Bind(ColorOnSurface, ColorOnSurfaceDark),
			Bind(ColorMuted, ColorMutedDark),
			Bind(ColorHover, ColorHoverDark),
		),
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
