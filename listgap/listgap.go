//go:build !wasm

// Package listgap owns one concern: the container part shared by every target*
// list component (targetlist, targetdate, …). It is a lego piece — the list
// components assemble it, they do not re-declare its shape or its spacing.
//
// Apply wires the desktop PartList contract in a single call:
//
//   - the vertical row rhythm (the "more air" knob between rows)
//   - the Scroll() region, which reserves the block edges for a FloatingChrome
//     host through var(--floating-top|bottom, 0px) — 0px when there is none
//   - ScrollGutter(inset): an ambient top/bottom gutter equal to the lateral
//     inset, so the list looks framed the same on all four sides instead of
//     sitting flush against the top while the sides carry a visible gutter.
//     It is ADDITIVE with the FloatingChrome reservation above (folded into
//     the same calc), never a plain padding-block that would replace it — a
//     `padding:`/PadEdge() shorthand here would land in the widgets layer and
//     clobber the seam's reservation outright, putting the last row back
//     under the host's floating button (measured: 72px on medicalhistory's
//     own crudview action). See widget/style's ScrollGutter doc for why a
//     plain override cannot do this safely.
//   - the lateral inset value, calibrated to crudview's master-detail indent
//     budget (16px total: crudview card 4 + list 4 + this 8). Retuning it means
//     re-checking crudview/css.go in webtyp/layout.
//
// MobileOpts carries the phone-breakpoint overrides. The caller spreads them
// into its OWN On(css.Mobile, …) — the conformance rail keeps the css import
// (and therefore the breakpoint device) inside css.go files, so listgap hands
// the values out as style.Options instead of naming the device itself.
//
// The four space steps are unexported on purpose: a caller cannot pass the
// wrong one to the wrong property, because the Stack()/PadInline() calls live
// here. That is also why the two list components stay visually interchangeable
// (crudview swaps one for the other) with no "keep in sync with the sibling"
// rule to remember.
//
// Leaf package: components/conformance_test.go imports targetlist, so a symbol
// here — rather than in the components root — keeps the test build acyclic.
package listgap

import (
	"webtyp.com/widget"
	"webtyp.com/widget/style"
)

// rowGap must clear PartBadge's own straddle: the badge hangs
// calc(-0.5 * var(--chip-height,1.25rem)) below its row's bottom edge — 10px
// at the default chip-height — so it visually eats into whatever sits below
// it. A gap smaller than that overshoot puts the badge on top of the NEXT
// row (measured: Space2/8px left it overlapping by 2px). Both steps below
// clear it with margin; re-check this if --chip-height is ever retuned.
const (
	rowGap       = style.Space4 // 16px — clears the 10px badge overshoot by 6px
	rowGapMobile = style.Space6 // 24px — clears it by 14px, plus mobile's own air

	inset       = style.Space1 // 4px  — lateral gutter; part of crudview's 16px budget
	insetMobile = style.Space2 // 8px
)

// Apply adds the desktop list-container contract for part `list` to s and
// returns s for chaining. Both target* lists call this instead of hand-writing
// the Part() block, so the container shape and its two spacing values cannot
// drift between them.
func Apply(s *style.Sheet, list widget.Part) *style.Sheet {
	return s.Part(list,
		style.Stack(rowGap),
		style.Scroll(),
		style.ScrollGutter(inset),
		style.PadInline(inset),
	)
}

// MobileOpts is the phone-breakpoint override for the list container: a larger
// row rhythm, the top/bottom gutter and the lateral inset both re-tied to the
// same (larger) mobile inset so the symmetry holds at this breakpoint too.
// Spread it into the caller's own On(css.Mobile, list, listgap.MobileOpts()...).
//
// Scroll() is repeated here on purpose: On() emits a part's flow/scroll shape
// as its own block in the media query, and ScrollGutter's calc rides that
// scroll emission — without Scroll() in the override the mobile block would
// keep the desktop gutter value.
func MobileOpts() []style.Option {
	return []style.Option{
		style.Stack(rowGapMobile),
		style.Scroll(),
		style.ScrollGutter(insetMobile),
		style.PadInline(insetMobile),
	}
}
