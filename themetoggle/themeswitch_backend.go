//go:build !wasm

package themetoggle

import . "github.com/tinywasm/dom"

func (t *ThemeToggle) Init(_ Ctx) {
	t.theme = NewString(string(ThemeDark))
}

func (t *ThemeToggle) onClick() {}
