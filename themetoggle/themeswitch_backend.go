//go:build !wasm

package themetoggle

import . "github.com/tinywasm/dom"

func (t *ThemeToggle) Init(_ Ctx) {
	t.theme = NewString(string(defaultTheme))
}

func (t *ThemeToggle) onClick() {}
