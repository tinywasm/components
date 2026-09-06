//go:build !wasm

package listselect

import (
	"webtyp.com/css"
	"webtyp.com/widget"
	"webtyp.com/widget/style"
)

// ApplyRow adds the per-row selection-check skin to s (the host widget's
// sheet). The host calls this instead of hand-writing the block. row is the
// host's own row part (its class differs per widget), used only for the
// danger wash under a marked row.
//
// The check rides its row's top-end corner (OnEdge), out of the flow: the
// label never shifts when selection mode opens, and the box reads as a badge
// on the row rather than a column of its own.
//
// Three reveal facts, and each hangs off exactly the element that owns it:
//
//   - The BOX exists only inside the list root's Open state (selection mode).
//     In normal mode the root has no data-open, so there is no square at all —
//     not an empty one. Hide() is the sole display authority in the base rule;
//     the flex centring lives in the Open reveal, so nothing in the base
//     competes with it (a base rule mixing Hide() with CenterContent() emits
//     two `display` values and the revealed box comes back as a block, which
//     stacks its glyph off-centre — the exact defect this split fixes).
//
//   - WHICH GLYPH shows, and what colour the box wears, is the BOX's own
//     Selected/Invalid state — written by the row ONLY in selection mode (see
//     RowOf). Never the ROW's Selected/Invalid: those also mean "this is the
//     loaded record" in normal mode, and a glyph must never appear then.
//     Invalid → trash on a solid Danger box (white glyph via --color-on-danger);
//     Selected → pencil on a solid AccentInverse box (white via
//     --color-on-primary). A plain Accent/Inset box would tint the glyph
//     near-black through currentColor.
//
//   - The ROW carries the danger wash under the whole row while it is marked
//     for delete (the row binds Invalid too, alongside the box).
func ApplyRow(s *style.Sheet, row widget.Part) *style.Sheet {
	return s.
		Part(partCheck,
			style.Hide(),
			style.OnEdge(style.EdgeTop, style.SideEnd, style.SpaceNone, style.SpaceNone),
			style.IconBox(style.IconMd),
			style.KeepSize(),
			style.Round(style.RadiusSm),
			// The resting square in selection mode: marked nothing yet. A
			// checked box overrides this from its own state rule below
			// (higher specificity: .check[data-x] beats .check).
			style.As(style.Inset),
			style.Animate(style.MotionFast),
		).
		Part(partCheckTrash,
			style.Hide(),
			style.IconBox(style.IconSm),
		).
		Part(partCheckPencil,
			style.Hide(),
			style.IconBox(style.IconSm),
		).
		// Reveal the square only inside selection mode, as a centred flex box
		// so its single visible glyph sits dead centre. Row(SpaceNone) gives
		// Show() a flex display to restore; CenterContent aligns both axes.
		WhenWithin(widget.Open, "", partCheck,
			style.Show(),
			style.Row(style.SpaceNone),
			style.CenterContent(),
		).
		// Glyph reveal + box colour: the box's OWN state, set only in
		// selection mode.
		WhenWithin(widget.Invalid, partCheck, partCheckTrash,
			style.Show(),
		).
		WhenWithin(widget.Selected, partCheck, partCheckPencil,
			style.Show(),
		).
		When(widget.Invalid, partCheck,
			style.As(style.Danger),
		).
		When(widget.Selected, partCheck,
			style.As(style.AccentInverse),
		).
		// The whole row, washed red while marked for delete.
		When(widget.Invalid, row,
			style.As(style.DangerWash),
		)
}

// ApplyHeader adds the in-flow selection-header skin to s. The strip is the
// host root's first child, ALWAYS in flow: it never carries Hide(), so the
// list never shifts when selection mode opens or closes. Its reserved
// height is NOT a bespoke min-height and not the app's full control-row
// height (that reserves far more than an icon needs): partAllSpacer is a
// child that stays visible in every mode, sized like the select-all box
// (IconBox(IconMd)), so the strip's natural flex height IS that icon's
// footprint — nothing bigger, nothing declared twice.
//
// The count is ALSO always visible — "k / N" is useful outside selection
// mode too (N alone answers "how many records are there"), so it is never
// gated on Open. Only the select-all BOX is: it is an action control that
// means nothing without selection mode, so it alone stays hidden until
// then. The always-visible spacer mirrors the box's exact footprint on the
// leading edge, which is what keeps the count truly centred in BOTH modes
// instead of recentring itself the moment the box appears — Header() builds
// the markup in that [spacer][count][box] order; PushEnd sends the box past
// the count to the strip's trailing edge without needing a fixed-width
// count.
//
// Unlike the per-row check (ApplyRow), nothing here is OnEdge: the strip is
// in flow, a normal row above the <ul>, so nothing can overlap the first row.
func ApplyHeader(s *style.Sheet) *style.Sheet {
	return s.
		Part(partHeader,
			style.Row(style.Space2),
		).
		Part(partAllSpacer,
			style.IconBox(style.IconMd),
			style.KeepSize(),
		).
		Part(partAll,
			style.Hide(),
			style.PushEnd(),
			style.IconBox(style.IconMd),
			style.KeepSize(),
			style.Round(style.RadiusSm),
			// The resting square: nothing marked yet, or not in selection
			// mode at all. A checked box overrides this from its own state
			// rule below (higher specificity: .sel-all[data-x] beats
			// .sel-all).
			style.As(style.Inset),
			style.Animate(style.MotionFast),
		).
		// No Hide() here: it's the box's only child and rides whatever
		// visibility the box itself has — unlike the row's check glyphs
		// (ApplyRow), there is no second glyph to swap between.
		Part(partAllIcon,
			style.IconBox(style.IconSm),
		).
		// Always visible — see the doc comment above — so Row+CenterContent
		// live directly in the base rule. That's safe ONLY because this
		// part carries no Hide(): mixing Hide() with CenterContent() in one
		// rule is what emits two conflicting `display` values elsewhere in
		// this file (see ApplyRow's doc comment) — there is no such
		// conflict here since display:flex is the only display this rule
		// ever states.
		Part(partAllCount,
			style.Grow(),
			style.Row(style.SpaceNone),
			style.CenterContent(),
			style.FontSize(style.TextXs),
			style.FontWeight(style.WeightBold),
		).
		// Reveal the box only inside selection mode (see the doc comment
		// above for why the count needs no such reveal).
		WhenWithin(widget.Open, "", partAll,
			style.Show(),
			style.Row(style.SpaceNone),
			style.CenterContent(),
		).
		// The box's OWN state (set by Header from the Mode) still tones the
		// background — Danger red while marked for delete, Accent amber
		// while marked for a bulk edit — only the glyph on top stopped
		// swapping. Solid fills — white glyph, opaque box.
		When(widget.Invalid, partAll, style.As(style.Danger)).
		When(widget.Selected, partAll, style.As(style.Accent))
}

// Ensure compile check for css import
var _ = css.Mobile
