//go:build !wasm

package navbar

import . "github.com/tinywasm/css"

func SSRInstance() *NavBar { return &NavBar{} }

func (n *NavBar) RenderCSS() *Stylesheet {
	return New(
		Rule(clsNav,
			BackgroundColor(ColorSurface),
			Color(ColorOnSurface),
			Padding(Space2),
		),
		Rule(clsNavList,
			RuleContent(Decl{Prop: "list-style", Val: "none"}),
			Padding(Zero),
			Margin(Zero),
			Display(Flex_),
			Gap(Space2),
		),
		Rule(Selector("."+string(clsNavItem)+" a"),
			Color(Str("inherit")),
			RuleContent(Decl{Prop: "text-decoration", Val: "none"}),
			Display(Flex_),
			AlignItems(Center),
			Gap(Space1),
		),
		Rule(Selector("."+string(clsNavItem)+" a:hover"),
			RuleContent(Decl{Prop: "text-decoration", Val: "underline"}),
			Color(ColorPrimary),
		),
		Rule(clsIcon,
			Width(Em(1)),
			Height(Em(1)),
			RuleContent(Decl{Prop: "fill", Val: "currentColor"}),
		),
	)
}

func (n *NavBar) IconSvg() map[string]string {
	return nil
}
