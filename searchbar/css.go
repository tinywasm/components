//go:build !wasm

package searchbar

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the search bar's visual contract using the style DSL.
func (s *SearchBar) RenderCSS() *css.Stylesheet {
	return style.For(s).
		// The bar is ONE control, not a card holding two loose pieces: the
		// magnifier is the bar's leading cap, the input its body, and a gap or a
		// card of its own between them saws the bar back into separate boxes.
		// The root carries the radius and clips, so cap and body stay square and
		// still read as one rounded bar.
		// ControlBox pins the bar to --control-height, the token every control
		// in the ecosystem answers to — that is what lets a host stack this bar
		// against its own buttons and have the heights agree by construction.
		Root(
			style.Row(style.SpaceNone),
			style.Round(style.RadiusMd),
			style.HideOverflow(),
			style.ControlBox(),
			style.KeepSize(),
		).
		// The magnifier is the bar's square cap: aspect-ratio, not padding, sets
		// the width — a padded box drifts off the control token (the old
		// Pad(Space2)+icon-box measured 40px against the host's 66), while the
		// square derives from the same --control-height as everything else.
		Part(PartIcon,
			style.As(style.Primary),
			style.MediaBox(style.AspectSquare),
			style.ControlBox(),
			style.KeepSize(),
		).
		// A bare <svg> with no box falls back to 300x150; IconBox pins it.
		Part(PartGlyph,
			style.IconBox(style.IconMd),
		).
		// The input is the body of the bar: it grows into whatever the cap
		// leaves and answers to the same control height, so cap and body can
		// never drift apart vertically — the mismatch that left a 25px field
		// floating in the middle of a 72px strip.
		Part(PartInput,
			style.As(style.Inset),
			style.Pad(style.Space2),
			style.Grow(),
			style.ControlBox(),
		).
		Stylesheet()
}
