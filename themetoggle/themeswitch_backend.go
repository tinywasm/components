//go:build !wasm

package themetoggle

import . "github.com/tinywasm/dom"

func (t *ThemeToggle) Init(_ Ctx) {
	t.theme = NewString("")
}

func (t *ThemeToggle) onClick() {}
