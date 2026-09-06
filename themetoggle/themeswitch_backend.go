//go:build !wasm

package themetoggle

import . "webtyp.com/dom"

func (t *ThemeToggle) Init(_ Ctx) {
	t.theme = NewString(string(defaultTheme))
	t.supported = true // no browser to probe at SSR time — see the field's doc comment
}

func (t *ThemeToggle) onClick() {}
