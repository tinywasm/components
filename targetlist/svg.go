//go:build !wasm

package targetlist

import "github.com/tinywasm/svg/sprite"

// IconSvg registers the ⋮ (vertical dots) options glyph. Uses currentColor via the
// sprite default so CSS controls its color.
func (t *TargetList) IconSvg() *sprite.Sprite {
	return sprite.NewSprite(
		sprite.Define(iconDots, "0 0 16 16",
			sprite.Path("M8 2a1.6 1.6 0 1 0 0 3.2A1.6 1.6 0 0 0 8 2M8 6.4a1.6 1.6 0 1 0 0 3.2A1.6 1.6 0 0 0 8 6.4M8 10.8a1.6 1.6 0 1 0 0 3.2A1.6 1.6 0 0 0 8 10.8"),
		),
	)
}
