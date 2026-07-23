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
		// The field box. Padding uses tinywasm/css's public spacing scale
		// (Space1/Space2/Space3), not ad hoc em values, so this box's
		// vertical rhythm stays numerically consistent with other
		// components (e.g. targetlist's row padding) that share the scale.
		Rule(Selector(".tw-field"),
			Position(Relative),
			Display(Flex_),
			FlexDirection(Column),
			Background(ColorBackground),
			Border(Str("1px solid "+ColorMuted.Var())),
			BorderRadius(Em(0.4)),
			// 4 separate longhands, not one Padding(a,b,c) shorthand call —
			// joinValues()'s multi-token output (multiple var()s) loses its
			// spaces somewhere downstream in the CSS pipeline, rendering as
			// invalid `padding:var(...)var(...)var(...)` with no separators
			// (verified via the live stylesheet). A single-token longhand
			// per side sidesteps that entirely.
			PaddingTop(Space2),
			PaddingRight(Space3),
			PaddingBottom(Space1),
			PaddingLeft(Space3),
			MarginBottom(Rem(0.7)),
			Transition(Str("box-shadow .2s ease, border-color .2s ease")),
		),

		// The label chip (legend), centered exactly ON the box's top border —
		// half above (outside the box), half below (inside it). `position:
		// absolute`'s containing-block origin is the border's inner edge, so
		// `top: 0` already sits ON the border with no padding math needed;
		// `translateY(-50%)` then shifts it up by exactly half of its OWN
		// rendered height, so the straddle stays exact regardless of the
		// chip's font-size/padding. (Previously a flex child with a fixed
		// `-1em` margin, which undershot the true half-height offset.)
		Rule(Selector(".tw-field label"),
			Position(Absolute),
			Top(Zero),
			Left(Space3),
			Margin(Zero),
			// Adjacent RawRules concatenate with no separator — this one
			// MUST carry its own trailing ';' or it silently glues onto the
			// next Decl (Background below) into one invalid `transform`
			// value that the browser drops entirely (verified via
			// getComputedStyle().transform === "none" without the ';').
			RawRule("transform: translateY(-50%);"),
			Background(ColorPrimary),
			Color(ColorOnPrimary),
			FontSize(Rem(0.72)),
			FontWeight(Str("600")),
			Padding(Em(0.1), Em(0.6)),
			BorderRadius(Em(0.35)),
			Cursor(Pointer),
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

		// Focus (ColorPrimary — the field is active/being edited) + hover
		// (ColorHover — the dedicated hover token, consistent with every other
		// hover indicator in the app, not a repurposed accent color).
		Rule(Selector(".tw-field:focus-within"),
			BorderColor(ColorPrimary),
			RawRule("box-shadow: inset 0 0 0 1px "+ColorPrimary.Var()),
		),
		Rule(Selector(".tw-field:hover"),
			BorderColor(ColorHover),
		),

		// The error message sits inside the box's top-right corner — same top/
		// right offset as the box's own padding, so its margins read as even —
		// transparent (no chip/badge look) and absolutely positioned so it
		// NEVER affects the box's height, even across multiple lines of text.
		// Hidden by default; form/render_input.go toggles tw-field-error--visible
		// only when there's an actual validation error.
		Rule(Selector(".tw-field .tw-field-error"),
			Position(Absolute),
			Top(Space2),
			Right(Space3),
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
