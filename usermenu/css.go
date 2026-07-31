//go:build !wasm

package usermenu

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the usermenu visual contract using the style DSL.
func (m *UserMenu) RenderCSS() *css.Stylesheet {
	return style.For(m).
		// Anchor, not a flow: the panel hangs off this element, and a corner
		// pin needs a positioned ancestor to hang from.
		Root(
			style.Anchor(),
			style.KeepSize(),
		).
		// The resting state: avatar and name, nothing else. Whatever a shell
		// puts here has to have a bounded width, and a name does.
		Part(PartTrigger,
			style.Row(style.Space2),
			style.Interactive(style.Subtle),
			style.Pad(style.Space1),
			style.Round(style.RadiusFull),
			style.KeepSize(),
		).
		Part(PartAvatar,
			style.IconBox(style.IconLg),
			style.Round(style.RadiusFull),
			style.HideOverflow(),
		).
		Part(PartName,
			style.FontSize(style.TextBase),
			style.FontWeight(style.WeightBold),
			style.KeepSize(),
		).
		// Flyout, not flow: an inline panel would push the header's row open
		// and shove the content under it down every time someone looked at who
		// they were logged in as.
		Part(PartPanel,
			style.Stack(style.Space2),
			style.As(style.Panel),
			style.Round(style.RadiusMd),
			style.Raise(style.Floating),
			style.Pad(style.Space2),
			style.Width(style.Content),
			style.Flyout(style.SideEnd),
		).
		Part(PartRoles,
			style.Row(style.Space1),
		).
		// The same chip box the rest of the system uses, so a role beside a
		// field legend or a list badge reads as the same kind of object.
		Part(PartRole,
			style.As(style.Inset),
			style.Round(style.RadiusSm),
			style.FontSize(style.TextXs),
			style.ChipBox(),
			style.CenterContent(),
		).
		Part(PartActions,
			style.Row(style.Space1),
		).
		// Transparent and only present while the menu is open: it exists to
		// catch the click that closes it, not to dim anything.
		Part(PartBackdrop,
			style.Backdrop(style.Viewport),
			style.RevealedBy(widget.Open),
		).
		Stylesheet()
}
