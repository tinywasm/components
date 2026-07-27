//go:build !wasm

package modaldialog

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget/style"
)

// Style defines the modaldialog visual contract using the style DSL.
func (m *ModalDialog) Style() *style.Sheet {
	return style.Of(NameModalDialog).
		Root(
			style.Fill(),
			style.Fixed(),
			style.Cover(),
			style.On(style.Page),
		).
		Part(PartBackdrop,
			style.On(style.Sunken),
		).
		Part(PartPanel,
			style.On(style.Panel),
			style.Round(style.RadiusLg),
			style.Raise(style.Overlay),
			style.Pad(style.Space3),
		).
		Part(PartHeader,
			style.Row(style.Space2),
		).
		Part(PartBody,
			style.Stack(style.Space1),
		)
}

// RenderCSS returns the compiled CSS stylesheet.
func (m *ModalDialog) RenderCSS() *css.Stylesheet {
	return m.Style().Stylesheet()
}
