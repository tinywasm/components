//go:build !wasm

package themetoggle

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the themetoggle visual contract using style DSL.
func (t *ThemeToggle) RenderCSS() *css.Stylesheet {
	return style.Of(NameThemeToggle).
		Root(
			style.Fixed(),
			style.On(style.Accent),
			style.Round(style.RadiusFull),
			style.Pad(style.Space1),
		).
		Cue(widget.Hover, "",
			style.On(style.AccentHover),
		).
		Stylesheet()
}
