//go:build !wasm

package targetdate

import (
	"github.com/tinywasm/icons/pencil"
	"github.com/tinywasm/icons/selectall"
	"github.com/tinywasm/icons/trash"
	"github.com/tinywasm/svg/sprite"
)

// IconSvg registers the shared trash/pencil action glyphs from tinywasm/icons,
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
