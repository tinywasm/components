//go:build !wasm

package fieldset

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the fieldset widget visual contract using the style DSL.
func (f *Fieldset) RenderCSS() *css.Stylesheet {
	return style.For(f).
		// The field is not a card. It is a label chip sitting on an input box:
		// wrapping the pair in a second bordered panel nests a box inside a box
		// and reads as clutter once several fields stack up.
		// Pad, not a gap: the <form> that holds the fields carries no class, so
		// the only place to put air between consecutive fields is inside each
		// one.
		Root(
			style.Anchor(),
			style.Stack(style.SpaceNone),
			style.Pad(style.Space2),
			style.KeepSize(),
		).
		// A legend, not a caption: OnEdge centres the chip ON the input's top
		// border rather than stacking it above, and Space2 is the field's own
		// padding — the distance from the field's border to the input's.
		// No vertical Pad: the chip's box is its one TextXs line, and OnEdge
		// mounts it by half the declared --chip-height — the same token the
		// badge of a targetlist row straddles the line by. Both hang the same
		// 10px below the border, so the legend and the badge align by
		// construction, not by a padding value tuned against another repo's
		// skin. PadInline is exempt: it never touches the height, and text
		// flush against the chip's filled edge reads as a bug, not as density.
		// No ChipBox: a legend is read, not lined up. ChipBox pins every chip to
		// --chip-width and clips the overflow, which is right for a badge in a
		// column of badges (they align) and wrong here — one field per row means
		// there is no column to align with, so a short label like "id" padded a
		// fixed 112px reads as a mislabelled box, and a long one is silently cut
		// mid-word. Shrink-wrapped to its text, the chip says exactly its own
		// name at any length.
		//
		// Capitalize: the text arrives as a model's field name (id, name, ip),
		// lower-cased because that is what the struct tag says. Casing it here
		// keeps that decision in the skin instead of asking every model to
		// pre-format a display string.
		Part(widget.PartLabel,
			style.OnEdge(style.EdgeTop, style.SideStart, style.Space2, style.Space4),
			style.As(style.Primary),
			style.Round(style.RadiusSm),
			style.FontSize(style.TextXs),
			style.FontWeight(style.WeightBold),
			style.Raise(style.Raised),
			style.Capitalize(),
			style.StartContent(),
			style.PadInline(style.Space2),
		).
		// Red text growing leftward from the field's trailing edge — not a filled
		// bar across the whole field, and not riding the border either: OnEdge
		// straddled the input's top line and the message came out cut by it.
		// It sits inside the box with equal air above and to the right.
		// Space4, not Space2: an absolute box is laid out against the Root's
		// PADDING box, so the first Space2 only buys back the Root's own Pad and
		// leaves the text flush with the input's border. Space4 spends that Space2
		// and puts a real Space2 between the message and both edges.
		Part(widget.PartError,
			style.Docked(style.Parent, style.EdgeTop, style.SideEnd, style.Space4),
			style.Glyph(style.Danger),
			style.FontSize(style.TextXs),
		).
		// Space4, not Space2: the legend rides the input's top border and hangs
		// half its height inside the box, so the value needs room to clear it.
		// At the chip's 20px height the hang is 10px and Space4's 16px still
		// clears it.
		// Page, not Panel: an editable control is the one place on the screen
		// asking to be written in, so it gets the whitest surface the theme
		// has. Against the sunken card the form sits in, that white IS the
		// affordance — it reads as an empty writing area at a glance, without
		// a border having to announce it.
		Part(widget.PartInput,
			style.As(style.Page),
			style.Round(style.RadiusMd),
			style.Pad(style.Space4),
			style.ControlBox(),
		).
		Part(widget.PartRadioGroup,
			style.Row(style.Space3),
			style.Pad(style.Space1),
		).
		// Locked repaints the INPUT, not the field root. Painting the root was
		// what produced the "barrier": a state rule emits its border as an
		// outline (so the box cannot resize mid-hover), outlines paint last in
		// a stacking context — above positioned descendants — and the legend
		// straddles exactly that top border line, so the barrier drew a grey
		// rule straight across the chip. Nothing about "this field is
		// read-only" needs a second box around the pair; the control itself is
		// what is locked, so the control is what changes.
		//
		// Secondary is Panel's fill with no border: the control stops reading
		// as a writing area and settles back into the card, which is what a
		// value you are only reading should do. The text stays ColorOnSurface
		// at full contrast — a read-only value is still there to BE read, so
		// this deliberately does not reach for Inactive's muted text.
		// WhenWithin, not When: dom writes data-locked on the field wrapper
		// that owns the gate, so When(Locked, PartInput) would emit
		// `.tw-field__input[data-locked="true"]` and match nothing.
		// Round(RadiusMd) is not redundant: a surface with no explicit radius
		// falls back to its own default, and Secondary's is RadiusSm, so
		// locking silently shrank the control's corners from the 8px the base
		// rule gives it to 4px. Restating the base radius keeps the box the
		// same shape in both states — only its fill is supposed to change.
		WhenWithin(widget.Locked, "", widget.PartInput,
			style.As(style.Secondary),
			style.Round(style.RadiusMd),
		).
		// No filled state on the field itself: the red message in the corner is
		// the signal. Painting the whole box danger buried the value.
		When(widget.Invalid, widget.PartInput,
			style.Glyph(style.Danger),
		).
		Stylesheet()
}
