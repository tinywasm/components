//go:build !wasm

package contentcard

import . "github.com/tinywasm/css"

func (c *ContentCard) RenderCSS() *Stylesheet {
	return NewStylesheet(
		Rule(clsCard,
			Border(Px(1), Str("solid"), ColorMuted),
			BackgroundColor(ColorBackground),
			BorderRadius(RadiusMd),
			BoxShadow(Str("0 1px 3px 0 rgba(0, 0, 0, 0.3)")),
			Display(Flex_),
			FlexDirection(Str("column")),
		),
		Rule(clsCardHeader,
			Padding(Space2),
			RuleContent(Decl{Prop: "border-bottom", Val: "1px solid " + ColorMuted.Var()}),
			FontWeight(FontWeightBold),
			Color(ColorOnSurface),
		),
		Rule(clsCardBody,
			Padding(Space2),
			Flex(Str("1")),
			Color(ColorOnSurface),
		),
		Rule(clsCardFooter,
			Padding(Space2),
			RuleContent(Decl{Prop: "border-top", Val: "1px solid " + ColorMuted.Var()}),
			BackgroundColor(ColorSurface),
			RuleContent(Decl{Prop: "border-bottom-left-radius", Val: RadiusMd.Var()}),
			RuleContent(Decl{Prop: "border-bottom-right-radius", Val: RadiusMd.Var()}),
			Color(ColorOnSurface),
		),
	)
}
