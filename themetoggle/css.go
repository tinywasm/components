//go:build !wasm

package themetoggle

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the themetoggle visual contract using style DSL.
func (t *ThemeToggle) RenderCSS() *css.Stylesheet {
	return style.For(t).
		Root(
			style.KeepSize(),
			style.Interactive(style.Primary),
			style.Round(style.RadiusFull),
			style.Pad(style.Space1),
		).
		Stylesheet()
}
