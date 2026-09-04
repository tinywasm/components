//go:build !wasm

package listselect

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
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
// host root's first child; it is hidden until the root's Open state (the host
// binds it from the Mode's On signal), then a centred flex row carrying the
// select-all box and the count.
//
// In normal mode the strip is display:none, so the host shows no header and
// no count until selection mode starts — identical to the pre-mode list.
//
// Unlike the per-row check (ApplyRow), nothing here is OnEdge: the strip is
// in flow, a normal row above the <ul>, so nothing can overlap the first row.
func ApplyHeader(s *style.Sheet) *style.Sheet {
	return s.
		Part(partHeader,
			style.Hide(),
		).
		Part(partAll,
			style.IconBox(style.IconMd),
			style.KeepSize(),
			style.Round(style.RadiusSm),
			style.As(style.Inset),
			style.Animate(style.MotionFast),
		).
		Part(partAllTrash,
			style.Hide(),
			style.IconBox(style.IconSm),
		).
		Part(partAllPencil,
			style.Hide(),
			style.IconBox(style.IconSm),
		).
		Part(partAllCount,
			style.FontSize(style.TextXs),
			style.FontWeight(style.WeightBold),
		).
		// Reveal the strip only inside selection mode, as a centred flex row.
		WhenWithin(widget.Open, "", partHeader,
			style.Show(),
			style.Row(style.Space2),
			style.CenterContent(),
		).
		// Glyph reveal: the box's OWN state (set by Header from the Mode),
		// never the row's.
		WhenWithin(widget.Invalid, partAll, partAllTrash, style.Show()).
		WhenWithin(widget.Selected, partAll, partAllPencil, style.Show()).
		// Solid fills — white glyph, opaque box (the old wash ruins read as
		// near-black glyphs on a pale fill).
		When(widget.Invalid, partAll, style.As(style.Danger)).
		When(widget.Selected, partAll, style.As(style.Accent))
}

// Ensure compile check for css import
var _ = css.Mobile
