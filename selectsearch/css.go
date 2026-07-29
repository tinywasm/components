//go:build !wasm

package selectsearch

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the selectsearch visual contract using style DSL.
func (c *SelectSearch) RenderCSS() *css.Stylesheet {
	return style.For(c).
		Root(
			style.Stack(style.Space1),
			style.As(style.Panel),
			style.Round(style.RadiusMd),
		).
		Part(PartToggle,
			style.As(style.Panel),
		).
		Part(PartDropdown,
			style.Stack(style.Space1),
			style.As(style.Panel),
			style.Raise(style.Floating),
			style.HideOverflow(),
		).
		Part(PartHeader,
			style.Row(style.Space2),
			style.Pad(style.Space2),
		).
		Part(PartIcon,
			style.As(style.Panel),
		).
		Part(PartSearch,
			style.Pad(style.Space2),
			style.As(style.Inset),
		).
		Part(PartOptions,
			style.Stack(style.SpaceNone),
			style.Scroll(),
		).
		Part(PartOption,
			style.Row(style.Space2),
			style.Pad(style.Space2),
			style.Interactive(style.Panel),
		).
		Part(PartLabel,
			style.As(style.Panel),
		).
		Part(PartDesc,
			style.As(style.Inset),
			style.FontSize(style.TextXs),
		).
		Stylesheet()
}
