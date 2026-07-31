//go:build !wasm

package modaldialog

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the modaldialog visual contract using the style DSL.
func (m *ModalDialog) RenderCSS() *css.Stylesheet {
	return style.For(m).
		// Backdrop(Viewport), not Fill(): in the flow the dialog was a flex item
		// of whatever mounted it and got dealt a column of that layout — a
		// sliver between the form and the list. A modal answers to the screen,
		// not to the box it happens to be mounted in.
		Root(
			style.Backdrop(style.Viewport),
			style.Veil(),
			style.CenterContent(),
			style.Pad(style.Space4),
		).
		// The wash and its blur live on the root, which is the element that
		// covers the screen. This one is only the click-catcher: transparent,
		// filling the root, and behind the panel.
		Part(PartBackdrop,
			style.Backdrop(style.Parent),
		).
		// Anchor is what puts the panel above the click-catcher. Both carry
		// z-index auto, and a positioned element paints after an unpositioned
		// one no matter the DOM order — so the catcher covered the panel until
		// the panel was positioned too.
		Part(PartPanel,
			style.Anchor(),
			style.As(style.Panel),
			style.Round(style.RadiusLg),
			style.Raise(style.Popover),
			style.Pad(style.Space3),
			style.Stack(style.Space2),
			style.Width(style.Content),
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
