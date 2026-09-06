package countbadge

import (
	. "webtyp.com/dom"
	. "webtyp.com/html"
	"webtyp.com/widget"
)

// NameCountBadge is the widget name for countbadge.
const NameCountBadge = widget.Name("countbadge")

const (
	// PartBadge is the bubble itself. It straddles the host's top-end
	// corner (see css.go); the host must be an Anchor — position: relative —
	// so the bubble resolves against the button, not the page.
	PartBadge = widget.Part("badge")
)

var (
	clsBadge = NameCountBadge.Class(PartBadge)
)

// CountBadge is a notification bubble: a count riding the top-end corner of
// the button that owns it. Out of the flow by construction (OnEdge), so the
// host keeps its box whether the bubble shows "1", "99" or nothing at all.
//
// Both signals are required:
//   - Count is the number to show, as text. It is a signal (not a string)
//     because the count changes without re-rendering the host.
//   - Visible decides whether the bubble paints. Showing "0" is noise — a
//     disabled commit button already says there is nothing to confirm — so
//     wire something that is false at zero (e.g. n > 0), not the mode.
type CountBadge struct {
	Element
	Count   *SignalString
	Visible *SignalBool
}

func (b *CountBadge) WidgetName() widget.Name { return NameCountBadge }
func (b *CountBadge) WidgetKind() widget.Kind { return widget.Region }

func (b *CountBadge) Render() *Element {
	return Span().Set(clsBadge.AsAttr()).
		BindText(b.Count).
		BindStateFunc(widget.Open, func() bool { return b.Visible.Get() })
}
