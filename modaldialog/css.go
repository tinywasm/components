//go:build !wasm

package modaldialog

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the modaldialog visual contract using the style DSL.
func (m *ModalDialog) RenderCSS() *css.Stylesheet {
	return style.For(m).
		Root(
			style.Fill(),
			style.KeepSize(),
			style.FillCentered(),
			style.As(style.Page),
		).
		Part(PartBackdrop,
			style.Backdrop(style.Parent),
			style.Veil(),
		).
		Part(PartPanel,
			style.As(style.Panel),
			style.Round(style.RadiusLg),
			style.Raise(style.Popover),
			style.Pad(style.Space3),
		).
		Part(PartHeader,
			style.Row(style.Space2),
		).
		Part(PartBody,
			style.Stack(style.Space1),
		).
		Part(PartClose,
			style.As(style.Panel),
		).
		Stylesheet()
}
