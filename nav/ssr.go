//go:build !wasm

package nav

import . "github.com/tinywasm/css"

func SSRInstance() *Nav { return &Nav{} }

func (n *Nav) RenderCSS() *Stylesheet {
	return New(
		Rule(ClsNav,
			BackgroundColor(ColorSurface),
			Color(ColorOnSurface),
			Padding(Space2),
		),
		Rule(ClsNavList,
			RuleContent(Decl{"list-style", "none"}),
			Padding(Zero),
			Margin(Zero),
			Display(Flex_),
			Gap(Space2),
		),
		Rule(Selector("."+string(ClsNavItem)+" a"),
			Color(Str("inherit")),
			RuleContent(Decl{"text-decoration", "none"}),
			Display(Flex_),
			AlignItems(Center),
			Gap(Space1),
		),
		Rule(Selector("."+string(ClsNavItem)+" a:hover"),
			RuleContent(Decl{"text-decoration", "underline"}),
			Color(ColorPrimary),
		),
		Rule(ClsIcon,
			Width(Em(1)),
			Height(Em(1)),
			RuleContent(Decl{"fill", "currentColor"}),
		),
	)
}

func (n *Nav) IconSvg() map[string]string {
	return nil
}
