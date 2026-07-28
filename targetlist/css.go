//go:build !wasm

package targetlist

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// Style defines the targetlist visual contract using the style DSL.
func (t *TargetList) Style() *style.Sheet {
	return style.Of(NameTargetList).
		Root(
			style.Fill(),
			style.Stack(style.Space0),
		).
		Part(PartList,
			style.Stack(style.Space3),
			style.Scrolls(),
			style.Pad(style.Space1),
		).
		Part(PartBackdrop,
			style.On(style.Sunken),
			style.Fixed(),
		).
		Part(PartRow,
			style.Row(style.Space2),
			style.On(style.Panel),
			style.Pad(style.Space3),
			style.Round(style.RadiusMd),
		).
		Part(PartLabel,
			style.On(style.Panel),
		).
		Part(PartBadge,
			style.On(style.Sunken),
			style.Round(style.RadiusSm),
			style.Text(style.TextXs),
		).
		Part(PartMenu,
			style.On(style.Panel),
		).
		Part(PartButton,
			style.On(style.Panel),
		).
		Part(PartIcon,
			style.On(style.Panel),
		).
		Part(PartOptions,
			style.Stack(style.Space0),
			style.On(style.Panel),
			style.Raise(style.Floating),
			style.Clip(),
		).
		Part(PartItem,
			style.On(style.Panel),
		).
		When(widget.Selected, PartRow,
			style.On(style.Selected),
		).
		Cue(widget.Hover, PartRow,
			style.On(style.PanelHover),
		)
}

// RenderCSS carries the ⋮ menu's backdrop: a full-viewport click-catcher that is
// hidden until a row's <details> menu is open, so clicking anywhere outside closes
// it (native <details> has no such behaviour on its own).
//
// KNOWN UPSTREAM GAP — widget/style has no vocabulary for an overlay: position
// outside the flow, full inset, stacking order, and "visible only while a sibling
// is open". style.Fixed() is the *does-not-reflow* exception (flex-grow/shrink: 0),
// NOT positioning — using it here silently turned the backdrop into a visible block
// in the flow and broke click-outside-to-close. Report the gap to tinywasm/widget;
// delete this method once the vocabulary exists.
//
// Selectors are DERIVED from the widget anatomy, never written by hand, so markup
// and stylesheet cannot drift apart.
func (t *TargetList) RenderCSS() *css.Stylesheet {
	var (
		root     = "." + clsListWrap.String()
		menu     = "." + clsMenu.String()
		backdrop = "." + clsMenuBackdrop.String()
	)
	return css.NewStylesheet(
		css.Raw(
			backdrop + "{display:none;position:fixed;top:0;left:0;right:0;bottom:0;z-index:4;}" +
				root + ":has(" + menu + "[open]) " + backdrop + "{display:block;}",
		),
	)
}
