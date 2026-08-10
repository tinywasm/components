//go:build !wasm

package targethour

import (
	"github.com/tinywasm/components/listgap"
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the targethour visual contract using the style DSL.
func (t *TargetHour) RenderCSS() *css.Stylesheet {
	return t.sheet().Stylesheet()
}

// sheet builds the style Sheet — same split-for-Validate() reason as
// targetlist's own sheet().
func (t *TargetHour) sheet() *style.Sheet {
	return style.For(t).
		Root(
			style.Fill(),
			style.Stack(style.SpaceNone),
		).
		Part(PartList,
			style.Stack(listgap.Gap),
			style.Scroll(),
			style.PadInline(listgap.Gap),
		).
		On(css.Mobile, PartList,
			style.Stack(listgap.GapMobile),
			style.PadInline(listgap.GapMobile),
		).
		// Same KeepSize reasoning as targetlist.TargetList's PartRow: the row is
		// a flex item inside a Scroll() column and must not shrink back under
		// its own open ⋮ options — see that file for the measured numbers.
		//
		// No Pad here on purpose, unlike targetlist's own PartRow: this row's
		// FIRST direct child is PartLead, a fixed square (ControlBox's
		// --control-height), and the reference this row is matching (a
		// legacy hour/day badge) has that square flush against the row's
		// leading, top, and bottom edges — no inset on those three sides.
		// PartContent below carries the padding instead, for everything
		// that isn't the square.
		Part(PartRow,
			style.Anchor(),
			style.Row(style.Space2),
			style.KeepSize(),
			style.ControlBox(),
			style.Interactive(style.Page),
			style.Round(style.RadiusMd),
		).
		// Everything that ISN'T the leading square: label, doctor badge,
		// trigger, options. Grow() claims the row's remaining width after
		// PartLead's fixed square; Pad(Space2) is the only inset in this row
		// — it reads as the row's own gutter on top/end/bottom, and as the
		// gap between the square and the label on start, so there is no
		// double gutter to reconcile with PartRow (which has none).
		Part(PartContent,
			style.Row(style.Space2),
			style.Grow(),
			style.Pad(style.Space2),
		).
		// The leading date badge: no Surface of its own anymore — a filled
		// block (the first version of this row used As(Primary)) fought the
		// row's own selected-state color for every descendant's text (see
		// the git history on this file for the -webkit-text-fill-color
		// saga) and, once fixed, still didn't match the reference this row
		// follows: a legacy hour/day badge that sits on the SAME background
		// as the rest of the row, marked off by a plain rule, not a filled
		// chip. Divider(SideEnd) draws exactly that rule, independent of any
		// Surface — see except.go's Divider comment for why it is not an
		// As() option. MediaBox(AspectSquare)+ControlBox still size it off
		// --control-height, matching PartRow's own floor, so the square
		// still tracks the row's height; ControlBox is a minimum, not a
		// cap, so Pad(Space2) is free to make the box grow past it rather
		// than clipping three lines of text against the edge.
		Part(PartLead,
			style.StartContent(),
			style.MediaBox(style.AspectSquare),
			style.ControlBox(),
			style.Pad(style.Space2),
			style.Divider(style.SideEnd),
			style.KeepSize(),
		).
		// The actual three-line column, left-aligned (StartContent) as
		// PartLead's one child — matching the reference's flush-left
		// weekday/day/month, not centered under the (now gone) filled
		// square that used to visually anchor a centered layout.
		Part(PartLeadStack,
			style.Stack(style.SpaceNone),
			style.StartContent(),
		).
		// FontSize kept small on every line (TextXs top/bottom, TextSm main,
		// not TextLg/TextXl) — three lines at default 1.5 line-height need
		// to clear roughly --control-height on their own.
		Part(PartLeadTop,
			style.FontSize(style.TextXs),
		).
		Part(PartLeadMain,
			style.FontSize(style.TextSm),
			style.FontWeight(style.WeightBold),
		).
		Part(PartLeadBottom,
			style.FontSize(style.TextXs),
		).
		Part(PartLabel,
			style.FontWeight(style.WeightBold),
			style.Grow(),
		).
		// Same straddle treatment as targetlist.TargetList's PartBadge — see
		// that file's comment for why it is a negative-margin overlap, not a
		// transform.
		Part(PartBadge,
			style.As(style.Inset),
			style.Round(style.RadiusSm),
			style.FontSize(style.TextXs),
			style.ChipBox(),
			style.CenterContent(),
			style.KeepSize(),
			style.OnEdge(style.EdgeBottom, style.SideEnd, style.SpaceNone, style.Space3),
		).
		// DOM-trailing (see targethour.go's buildRow), not PushEnd: PartLabel's
		// Grow() already claims every pixel of the row's free space during
		// flex resolution, so a margin-auto trick on a DOM-leading trigger
		// would have nothing left to distribute. targethour has no inherited
		// mobile-sliver-reachability history to preserve (unlike
		// targetlist.TargetList's PartButton), so there is no Mobile carve-out
		// pulling it back to the leading edge either.
		Part(PartButton,
			style.As(style.Subtle),
			style.KeepSize(),
			style.Interactive(style.Subtle),
			style.Round(style.RadiusSm),
			style.Pad(style.Space1),
			style.CenterContent(),
		).
		Part(PartIcon,
			style.As(style.Subtle),
			style.IconBox(style.IconMd),
		).
		Part(PartOptions,
			style.Row(style.Space2),
			style.Width(style.Full),
			style.KeepSize(),
			style.RevealedBy(widget.Open),
		).
		Part(PartItemDanger,
			style.Interactive(style.Danger),
			style.Round(style.RadiusSm),
			style.Pad(style.Space2),
			style.Width(style.Content),
		).
		On(css.Mobile, PartItemDanger,
			style.As(style.Danger),
			style.Round(style.RadiusMd),
			style.Pad(style.Space2),
			style.Width(style.Content),
			style.Raise(style.Floating),
			style.CenterContent(),
		).
		OnlyOn(css.Mobile, PartItemIcon,
			style.CenterContent(),
			style.IconBox(style.IconLg),
		).
		On(css.Mobile, PartItemLabel,
			style.Hide(),
		).
		When(widget.Selected, PartRow,
			style.As(style.Accent),
		).
		Cue(widget.Hover, PartRow,
			style.As(style.AccentWash),
		).
		Cue(widget.Focus, PartRow,
			style.As(style.AccentWash),
		).
		Cue(widget.Press, PartRow,
			style.As(style.Accent),
		)
}
