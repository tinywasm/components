//go:build !wasm

package targetdate

import (
	"webtyp.com/icons/pencil"
	"webtyp.com/icons/selectall"
	"webtyp.com/icons/trash"
	"webtyp.com/svg/sprite"
)

// IconSvg registers the shared trash/pencil action glyphs from webtyp/icons,
// the same two the rows render through currentColor inside the check box and
// the crud view's footer buttons ship from their own IconSvg(), plus
// selectall, the selection header's own select-all glyph (listselect.Header).
func (t *TargetDate) IconSvg() *sprite.Sprite {
	return sprite.NewSprite(
		trash.Def(),
		pencil.Def(),
		selectall.Def(),
	)
}
