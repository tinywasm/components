//go:build !wasm

package listselect

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// Apply adds the check part's skin to s and returns s for chaining. Both
// target* lists call this instead of hand-writing the block, so the check
// cannot drift between them.
//
// The check rides its row's top-end corner (OnEdge), out of the flow: the
// label never shifts when selection mode opens, and the box reads as a
// badge on the row rather than a column of its own. checkIcon is the glyph
// inside the box, hidden until the row is actually marked — an always
// painted tick reads as pre-selected. row carries the danger wash while the
// host's danger tone is armed (see Mode.SetDanger); Selected and Invalid
// never coincide on one row (the row binds Selected only when the tone is
// off), so the two fills cannot race.
func Apply(s *style.Sheet, check, checkIcon, row widget.Part) *style.Sheet {
	return s.
		Part(check,
			style.OnEdge(style.EdgeTop, style.SideEnd, style.SpaceNone, style.SpaceNone),
			style.IconBox(style.IconMd),
			style.KeepSize(),
			style.As(style.Inset),
			style.Round(style.RadiusSm),
			style.CenterContent(),
			style.Animate(style.MotionFast),
		).
		Part(checkIcon,
			style.Hide(),
		).
		WhenWithin(widget.Open, "", check,
			style.Show(),
		).
		WhenWithin(widget.Selected, "", checkIcon,
			style.Show(),
		).
		WhenWithin(widget.Invalid, "", checkIcon,
			style.Show(),
		).
		WhenWithin(widget.Invalid, "", check,
			style.As(style.Danger),
		).
		When(widget.Invalid, row,
			style.As(style.DangerWash),
		)
}

// Ensure compile check for css import
var _ = css.Mobile
