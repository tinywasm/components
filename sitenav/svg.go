//go:build !wasm

package sitenav

import "webtyp.com/svg/sprite"

// IconSvg registers the hamburger menu and close icon glyphs for mobile navigation.
func (sn *SiteNav) IconSvg() *sprite.Sprite {
	return sprite.NewSprite(
		sprite.Define(iconMenu, "0 0 24 24",
			sprite.Path("M3 18h18v-2H3v2zm0-5h18v-2H3v2zm0-7v2h18V6H3z"),
		),
		sprite.Define(iconClose, "0 0 24 24",
			sprite.Path("M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"),
		),
	)
}
