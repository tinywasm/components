//go:build !wasm

package datatable

import . "github.com/tinywasm/css"

func (t *DataTable) RenderCSS() *Stylesheet {
	return NewStylesheet(
		Rule(clsTable,
			Width(Pct(100)),
			RuleContent(Decl{Prop: "border-collapse", Val: "collapse"}),
			RuleContent(Decl{Prop: "text-align", Val: "left"}),
			RuleContent(Decl{Prop: "margin-bottom", Val: Space2.Var()}),
			Color(ColorOnSurface),
		),
		Rule(Selector("."+string(clsTable)+" th, ."+string(clsTable)+" td"),
			Padding(Rem(0.75)),
			RuleContent(Decl{Prop: "border-bottom", Val: "1px solid " + ColorMuted.Var()}),
		),
		Rule(Selector("."+string(clsTable)+" th"),
			RuleContent(Decl{Prop: "font-weight", Val: "600"}),
			BackgroundColor(ColorSurface),
		),
		Rule(Selector("."+string(clsTable)+" tbody tr:hover"),
			BackgroundColor(ColorHover),
		),
	)
}
