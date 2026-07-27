//go:build !wasm

package themetoggle

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// Style defines the themetoggle visual contract using style DSL.
func (t *ThemeToggle) Style() *style.Sheet {
	return style.Of(NameThemeToggle).
		Root(
			style.Fixed(),
			style.On(style.Accent),
			style.Round(style.RadiusFull),
			style.Pad(style.Space1),
		).
		Cue(widget.Hover, "",
			style.On(style.AccentHover),
		)
}

// RenderCSS returns compiled CSS.
func (t *ThemeToggle) RenderCSS() *css.Stylesheet {
	return t.Style().Stylesheet()
}
