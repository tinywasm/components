//go:build !wasm

package card

import . "github.com/tinywasm/css"

func SSRInstance() *Card { return &Card{} }

func (c *Card) RenderCSS() *Stylesheet {
	return New(
		Rule(ClsCard,
			Border(Px(1), Str("solid"), ColorMuted),
			BackgroundColor(ColorBackground),
			BorderRadius(RadiusMd),
			BoxShadow(Str("0 1px 3px 0 rgba(0, 0, 0, 0.3)")),
			Display(Flex_),
			FlexDirection(Str("column")),
		),
		Rule(ClsCardHeader,
			Padding(Space2),
			RuleContent(Decl{"border-bottom", "1px solid " + ColorMuted.Var()}),
			FontWeight(FontWeightBold),
			Color(ColorOnSurface),
		),
		Rule(ClsCardBody,
			Padding(Space2),
			Flex(Str("1")),
			Color(ColorOnSurface),
		),
		Rule(ClsCardFooter,
			Padding(Space2),
			RuleContent(Decl{"border-top", "1px solid " + ColorMuted.Var()}),
			BackgroundColor(ColorSurface),
			RuleContent(Decl{"border-bottom-left-radius", RadiusMd.Var()}),
			RuleContent(Decl{"border-bottom-right-radius", RadiusMd.Var()}),
			Color(ColorOnSurface),
		),
	)
}

func (c *Card) IconSvg() map[string]string {
	return nil
}
