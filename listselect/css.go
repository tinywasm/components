//go:build !wasm

package listselect

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// Apply adds the selection check's skin to s and returns s for chaining. Both
// target* lists call this instead of hand-writing the block, so the check
// cannot drift between them.
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
//     buildRow). Never the ROW's Selected/Invalid: those also mean "this is
//     the loaded record" in normal mode, and a glyph must never appear then.
//     Invalid → trash on a solid Danger box (white glyph via --color-on-danger);
//     Selected → pencil on a solid AccentInverse box (white via
//     --color-on-primary). A plain Accent/Inset box would tint the glyph
//     near-black through currentColor.
//
//   - The ROW carries the danger wash under the whole row while it is marked
//     for delete (the row binds Invalid too, alongside the box).
func Apply(s *style.Sheet, check, trashIcon, pencilIcon, row widget.Part) *style.Sheet {
	return s.
		Part(check,
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
		Part(trashIcon,
			style.Hide(),
			style.IconBox(style.IconSm),
		).
		Part(pencilIcon,
			style.Hide(),
			style.IconBox(style.IconSm),
		).
		// Reveal the square only inside selection mode, as a centred flex box
		// so its single visible glyph sits dead centre. Row(SpaceNone) gives
		// Show() a flex display to restore; CenterContent aligns both axes.
		WhenWithin(widget.Open, "", check,
			style.Show(),
			style.Row(style.SpaceNone),
			style.CenterContent(),
		).
		// Glyph reveal + box colour: the box's OWN state, set only in
		// selection mode.
		WhenWithin(widget.Invalid, check, trashIcon,
			style.Show(),
		).
		WhenWithin(widget.Selected, check, pencilIcon,
			style.Show(),
		).
		When(widget.Invalid, check,
			style.As(style.Danger),
		).
		When(widget.Selected, check,
			style.As(style.AccentInverse),
		).
		// The whole row, washed red while marked for delete.
		When(widget.Invalid, row,
			style.As(style.DangerWash),
		)
}

// ApplyMaster adds the master (select-all) check's skin to s. targetlist,
// targetdate and targethour call this instead of hand-writing the block.
//
// Same shape as Apply's per-row check: hidden until selection mode opens
// (Open on the list root), then a centred flex box carrying one glyph beside a
// small "n / total" count label. Which glyph, and what colour the box wears,
// is the MODE: Invalid on the box (danger tone armed) → trash on a red wash;
// Selected on the box (edit) → pencil on an amber wash. The none / some / all
// distinction is read from the count text ("0 / 15", "7 / 15", "15 / 15"), not
// from a third fill — a red "all" box in edit mode would lie about the mode.
//
// MAINTAINER: a solid fill at "all" is a reasonable future polish — it needs a
// compound state the DSL does not have today (mode AND count), so it is left
// out of this first pass rather than faked mode-blind.
func ApplyMaster(s *style.Sheet, checkAll, trashIcon, pencilIcon, count widget.Part) *style.Sheet {
	return s.
		Part(checkAll,
			style.Hide(),
			style.IconBox(style.IconMd),
			style.KeepSize(),
			style.Round(style.RadiusSm),
			style.As(style.Inset),
			style.OnEdge(style.EdgeTop, style.SideEnd, style.SpaceNone, style.Space2),
			style.Animate(style.MotionFast),
		).
		Part(trashIcon, style.Hide(), style.IconBox(style.IconSm)).
		Part(pencilIcon, style.Hide(), style.IconBox(style.IconSm)).
		Part(count,
			style.Hide(),
			style.FontSize(style.TextXs),
			style.FontWeight(style.WeightBold),
		).
		WhenWithin(widget.Open, "", checkAll,
			style.Show(),
			style.Row(style.Space1),
			style.CenterContent(),
		).
		WhenWithin(widget.Open, "", count,
			style.Show(),
		).
		WhenWithin(widget.Invalid, checkAll, trashIcon, style.Show()).
		WhenWithin(widget.Selected, checkAll, pencilIcon, style.Show()).
		When(widget.Invalid, checkAll, style.As(style.DangerWash)). // delete-select mode
		When(widget.Selected, checkAll, style.As(style.AccentWash)) // edit-select mode
}

// Ensure compile check for css import
var _ = css.Mobile
