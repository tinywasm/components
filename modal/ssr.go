//go:build !wasm

package modal

import . "github.com/tinywasm/css"

func SSRInstance() *Modal { return &Modal{} }

func (m *Modal) RenderCSS() *Stylesheet {
	return New(
		Rule(ClsModal,
			Position(Str("fixed")),
			Top(Zero),
			Left(Zero),
			Width(Pct(100)),
			Height(Pct(100)),
			ZIndex(ZModal),
			Display(Flex_),
			JustifyContent(Center),
			AlignItems(Center),
		),
		Rule(Selector("."+string(ClsModal)+"."+string(ClsModalHidden)),
			Display(None),
		),
		Rule(ClsModalBackdrop,
			Position(Str("absolute")),
			Top(Zero),
			Left(Zero),
			Width(Pct(100)),
			Height(Pct(100)),
			RuleContent(Decl{"background-color", "color-mix(in srgb, " + ColorSurface.Var() + ", transparent 40%)"}),
			ZIndex(Str("1")),
		),
		Rule(ClsModalContent,
			BackgroundColor(ColorBackground),
			Color(ColorOnSurface),
			BorderRadius(RadiusMd),
			BoxShadow(Str("0 4px 6px -1px rgba(0, 0, 0, 0.3)")),
			ZIndex(Str("2")),
			RuleContent(Decl{"min-width", "300px"}),
			RuleContent(Decl{"max-width", "90%"}),
		),
		Rule(ClsModalHeader,
			Display(Flex_),
			JustifyContent(Str("space-between")),
			AlignItems(Center),
			Padding(Space2),
			RuleContent(Decl{"border-bottom", "1px solid " + ColorMuted.Var()}),
		),
		Rule(Selector("."+string(ClsModalHeader)+" h2"),
			Margin(Zero),
			FontSize(Rem(1.25)),
		),
		Rule(ClsModalClose,
			Background(None),
			Border(None),
			FontSize(Rem(1.5)),
			Cursor(Pointer),
			Padding(Zero),
			LineHeight(Str("1")),
			Color(ColorOnSurface),
		),
		Rule(ClsModalBody,
			Padding(Space2),
		),
	)
}

func (m *Modal) IconSvg() map[string]string {
	return nil
}
