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
func Apply(s *style.Sheet, check widget.Part) *style.Sheet {
	return s.
		Part(check,
			style.Hide(),
			style.IconBox(style.IconMd),
			style.KeepSize(),
			style.As(style.Inset),
			style.Round(style.RadiusSm),
			style.Animate(style.MotionFast),
		).
		WhenWithin(widget.Open, "", check,
			style.Show(),
		).
		WhenWithin(widget.Selected, "", check,
			style.As(style.Accent),
		)
}

// Ensure compile check for css import
var _ = css.Mobile
