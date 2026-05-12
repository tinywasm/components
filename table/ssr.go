//go:build !wasm

package table

import . "github.com/tinywasm/css"

func SSRInstance() *Table { return &Table{} }

func (t *Table) RenderCSS() *Stylesheet {
	return New(
		Rule(ClsTable,
			Width(Pct(100)),
			RuleContent(Decl{"border-collapse", "collapse"}),
			RuleContent(Decl{"text-align", "left"}),
			RuleContent(Decl{"margin-bottom", Space2.Var()}),
			Color(ColorOnSurface),
		),
		Rule(Selector("."+string(ClsTable)+" th, ."+string(ClsTable)+" td"),
			Padding(Rem(0.75)),
			RuleContent(Decl{"border-bottom", "1px solid " + ColorMuted.Var()}),
		),
		Rule(Selector("."+string(ClsTable)+" th"),
			RuleContent(Decl{"font-weight", "600"}),
			BackgroundColor(ColorSurface),
		),
		Rule(Selector("."+string(ClsTable)+" tbody tr:hover"),
			BackgroundColor(ColorHover),
		),
	)
}

func (t *Table) IconSvg() map[string]string {
	return nil
}
