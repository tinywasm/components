//go:build !wasm

package targetdate

import "github.com/tinywasm/svg/sprite"

// IconSvg registers the checkmark glyph.
func (t *TargetDate) IconSvg() *sprite.Sprite {
	return sprite.NewSprite(
		sprite.Define(iconCheck, "0 0 16 16",
			sprite.Path("M13.78 4.22a.75.75 0 0 1 0 1.06l-6.25 6.25a.75.75 0 0 1-1.06 0L3.72 8.78a.75.75 0 0 1 1.06-1.06l2.22 2.22 5.72-5.72a.75.75 0 0 1 1.06 0z"),
		),
	)
}
