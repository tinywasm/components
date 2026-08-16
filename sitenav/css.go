//go:build !wasm

package sitenav

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the visual stylesheet for sitenav.
func (sn *SiteNav) RenderCSS() *css.Stylesheet {
	return style.For(sn).
		Root(
			style.Row(style.Space4),
			style.As(style.Page),
			style.Pad(style.Space3),
			style.CenterContent(),
		).
		Part(PartBrand,
			style.KeepSize(),
			style.CenterContent(),
		).
		Part(PartLogo,
			style.KeepSize(),
		).
		Part(PartToggle,
			style.IconBox(style.IconSm),
			style.As(style.Panel),
			style.KeepSize(),
		).
		Part(PartMenu,
			style.Row(style.Space4),
			style.Grow(),
			style.CenterContent(),
		).
		Part(PartNav,
			style.Row(style.Space3),
			style.Grow(),
			style.CenterContent(),
		).
		Part(PartLink,
			style.Pad(style.Space1),
			style.FontSize(style.TextSm),
			style.FontWeight(style.WeightMedium),
		).
		Part(PartActions,
			style.Row(style.Space2),
			style.KeepSize(),
		).
		Stylesheet()
}
