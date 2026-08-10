//go:build !wasm

package targethour

import "github.com/tinywasm/svg/sprite"

// IconSvg registers the ⋮ (vertical dots) options glyph and the trash can
// (Eliminar) — identical paths to targetlist's own, under a th- prefix so
// the two sprites never collide once ssr fuses every IconSvg() into one.
func (t *TargetHour) IconSvg() *sprite.Sprite {
	return sprite.NewSprite(
		sprite.Define(iconDots, "0 0 16 16",
			sprite.Path("M8 2a1.6 1.6 0 1 0 0 3.2A1.6 1.6 0 0 0 8 2M8 6.4a1.6 1.6 0 1 0 0 3.2A1.6 1.6 0 0 0 8 6.4M8 10.8a1.6 1.6 0 1 0 0 3.2A1.6 1.6 0 0 0 8 10.8"),
		),
		sprite.Define(iconDelete, "0 0 448 512",
			sprite.Path("M135.2 17.7L128 32H32C14.3 32 0 46.3 0 64s14.3 32 32 32H416c17.7 0 32-14.3 32-32s-14.3-32-32-32H320l-7.2-14.3C307.4 6.8 296.3 0 284.2 0H163.8c-12.1 0-23.2 6.8-28.6 17.7zM416 128H32L53.2 467c1.6 25.3 22.6 45 47.9 45H346.9c25.3 0 46.3-19.7 47.9-45L416 128z"),
		),
	)
}
