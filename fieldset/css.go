//go:build !wasm

package fieldset

import (
	. "github.com/tinywasm/css"
)

// RenderCSS is the form-field skin. It targets the structure tinywasm/form emits
// (the `.tw-field` wrapper, its `<label>`, and the input widget) so every form in
// the app renders as a labeled box: a chip legend over a white, bordered field.
// Ported from Dev/Pkg/UI/Components/fieldset. Consumes theme tokens only — the
// app's RootCSS owns their values.
func (f *Fieldset) RenderCSS() *Stylesheet {
	return NewStylesheet(
		// The field box.
		Rule(Selector(".tw-field"),
			Position(Relative),
			Display(Flex_),
			FlexDirection(Column),
			Background(ColorBackground),
			Border(Str("1px solid "+ColorMuted.Var())),
			BorderRadius(Em(0.4)),
			Padding(Em(0.5), Em(0.6), Em(0.25)),
			MarginBottom(Rem(0.7)),
			Transition(Str("box-shadow .2s ease, border-color .2s ease")),
		),

		// The label chip (legend), lifted to sit on the box's top border.
		Rule(Selector(".tw-field label"),
			AlignSelf(Str("flex-start")),
			Background(ColorPrimary),
			Color(ColorOnPrimary),
			FontSize(Rem(0.72)),
			FontWeight(Str("600")),
			Padding(Em(0.1), Em(0.6)),
			BorderRadius(Em(0.35)),
			Cursor(Pointer),
			RawRule("margin: -1em 0 .25em"),
		),

		// The widget input, borderless inside the box.
		Rule(Selector(".tw-field input, .tw-field textarea, .tw-field select"),
			Border(None),
			Background(Str("transparent")),
			Width(Pct(100)),
			Padding(Em(0.15), Zero),
			FontSize(Rem(1)),
			Color(ColorOnSurface),
			Outline(None),
		),

		// Focus + hover feedback on the whole box.
		Rule(Selector(".tw-field:focus-within"),
			BorderColor(ColorPrimary),
			RawRule("box-shadow: inset 0 0 0 1px "+ColorPrimary.Var()),
		),
		Rule(Selector(".tw-field:hover"),
			BorderColor(ColorPrimary),
		),

		// The error message sits inside the box's top-right corner — same top/
		// right offset as the box's own padding, so its margins read as even —
		// transparent (no chip/badge look) and absolutely positioned so it
		// NEVER affects the box's height, even across multiple lines of text.
		// Hidden by default; form/render_input.go toggles tw-field-error--visible
		// only when there's an actual validation error.
		Rule(Selector(".tw-field .tw-field-error"),
			Position(Absolute),
			Top(Em(0.5)),
			Right(Em(0.6)),
			Background(Str("transparent")),
			Color(ColorError),
			FontSize(Rem(0.68)),
			FontWeight(Str("600")),
			Opacity(0),
			Transition(Str("opacity .15s ease")),
		),
		Rule(Selector(".tw-field .tw-field-error--visible"),
			Opacity(1),
		),

		// Locked/read-only state (Form.SetLocked, or a field's own static
		// disabled flag): a "frosted glass" look — the value stays fully
		// legible/selectable, just a hair darker than the box's normal
		// background (ColorSurface, not the much darker ColorMuted — this is
		// not a disabled BUTTON, the text must stay easy to read/copy) and a
		// not-allowed cursor. No opacity dimming: that fades the text too,
		// which is exactly what must stay readable until the row's ⋮ menu →
		// Editar unlocks it.
		Rule(Selector(".tw-field:has(:disabled)"),
			Background(ColorSurface),
		),
		Rule(Selector(".tw-field input:disabled, .tw-field textarea:disabled, .tw-field select:disabled"),
			Cursor(Str("not-allowed")),
		),
	)
}
