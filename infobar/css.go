//go:build !wasm

package infobar

import (
	"webtyp.com/css"
	"webtyp.com/widget/style"
)

// RenderCSS defines the stylesheet for infobar.
func (ib *InfoBar) RenderCSS() *css.Stylesheet {
	return style.For(ib).
		Root(
			style.Row(style.Space4),
			style.As(style.Inset),
			style.Pad(style.Space2),
			style.FontSize(style.TextSm),
			style.HideOverflow(),
		).
		Part(PartItem,
			style.Row(style.Space2),
			style.CenterContent(),
			style.KeepSize(),
		).
		Part(PartIcon,
			style.IconBox(style.IconSm),
			style.KeepSize(),
		).
		Part(PartText,
			style.FontSize(style.TextSm),
		).
		Stylesheet()
}
