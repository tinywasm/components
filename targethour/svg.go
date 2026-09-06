//go:build !wasm

package targethour

import (
	"webtyp.com/icons/pencil"
	"webtyp.com/icons/selectall"
	"webtyp.com/icons/trash"
	"webtyp.com/svg/sprite"
)

// IconSvg registers the shared trash/pencil action glyphs (see
// targetlist/targetdate's IconSvg — same shape) plus selectall, the
// selection header's own select-all glyph (listselect.Header).
func (t *TargetHour) IconSvg() *sprite.Sprite {
	return sprite.NewSprite(trash.Def(), pencil.Def(), selectall.Def())
}
