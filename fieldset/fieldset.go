// Package fieldset is the global form-field skin for the tinywasm ecosystem.
//
// It renders nothing itself. It exists so the SSR pipeline collects its
// RenderCSS(), which styles the structure tinywasm/form emits for every field —
// a `.tw-field` wrapper containing a `<label>` and the widget's input — as the
// reference "labeled box" look (a chip legend over a white box).
//
// A consumer opts EVERY form in the app into this look exactly once, by depending
// on this package from its composition root. The CSS is emitted globally, so all
// forms in all panels render identically with zero per-field or per-model config.
// Swapping the whole app's form look = swapping this one dependency.
//
// It deliberately styles form elements (`.tw-field`), which the generic component
// guideline discourages — that is precisely this package's job: it is the one
// place the ecosystem centralizes form appearance.
package fieldset

import (
	. "github.com/tinywasm/dom"
)

// Fieldset is the skin's receiver type. The SSR collector instantiates it to call
// RenderCSS(); it carries no state and paints no markup of its own.
type Fieldset struct {
	Element
}
