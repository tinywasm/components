//go:build !wasm

package selectsearch

import "webtyp.com/svg/sprite"

func (c *SelectSearch) IconSvg() *sprite.Sprite {
	return sprite.NewSprite(
		sprite.Define(iconArrowDown, "0 0 16 16",
			sprite.Path("M1.5 4.5l6.5 7 6.5-7H1.5z"),
		),
	)
}
