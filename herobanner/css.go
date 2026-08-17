//go:build !wasm

package herobanner

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the visual stylesheet for herobanner.
func (h *HeroBanner) RenderCSS() *css.Stylesheet {
	return style.For(h).
		Root(
			style.Backdrop(style.Parent),
			style.Veil(),
			style.Pad(style.Space8),
			style.CenterContent(),
			style.Animate(style.MotionSlow),
			style.HideOverflow(),
		).
		Part(PartMedia,
			style.Backdrop(style.Parent),
			style.Fill(),
			style.AutoRotate(),
		).
		// Foreground(): PartMedia is a Backdrop(Parent), which pins it at the
		// local stacking level. The caption is in normal flow at z-index auto
		// and therefore paints UNDERNEATH it — invisible as soon as the
		// photograph actually fills the banner.
		Part(PartContent,
			style.Stack(style.Space4),
			style.Pad(style.Space6),
			style.CenterContent(),
			style.As(style.Page),
			style.Foreground(),
		).
		Part(PartTitle,
			style.FontSize(style.Text2xl),
			style.FontWeight(style.WeightBold),
		).
		Part(PartSubtitle,
			style.FontSize(style.TextLg),
		).
		Part(PartActions,
			style.Row(style.Space3),
			style.CenterContent(),
		).
		Stylesheet()
}
