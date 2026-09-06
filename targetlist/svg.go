//go:build !wasm

package targetlist

import (
	"webtyp.com/icons/pencil"
	"webtyp.com/icons/selectall"
	"webtyp.com/icons/trash"
	"webtyp.com/svg/sprite"
)

// IconSvg registers the glyphs the rows and the selection header use: the
// shared trash/pencil action glyphs from webtyp/icons (each row renders
// them through currentColor inside its check box, and the crud view's
// footer buttons ship the same two definitions from their own IconSvg() —
// same ids, so assetmin collapses them to one symbol each) plus selectall,
// the header's own select-all glyph (listselect.Header).
func (t *TargetList) IconSvg() *sprite.Sprite {
	return sprite.NewSprite(
		trash.Def(),
		pencil.Def(),
		selectall.Def(),
	)
}
