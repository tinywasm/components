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
		).
		Part(PartClose,
			style.On(style.Panel),
		)
}

// RenderCSS returns custom positioning rules for modaldialog backdrop overlay
// since absolute overlay positioning rules are not supported by the style.Sheet DSL.
func (m *ModalDialog) RenderCSS() *css.Stylesheet {
	return css.NewStylesheet(
		css.Raw(`
.modaldialog__backdrop {
	position: absolute;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	background-color: color-mix(in srgb, var(--color-surface) 60%, transparent);
	z-index: 1;
}
		`),
	)
}
