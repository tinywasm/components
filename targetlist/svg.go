//go:build !wasm

package targetlist

import (
	"github.com/tinywasm/icons/pencil"
	"github.com/tinywasm/icons/trash"
	"github.com/tinywasm/svg/sprite"
)

// IconSvg registers the glyphs the rows use: the shared trash/pencil action
// glyphs from tinywasm/icons. Each row renders them through currentColor
// inside the check box, and the crud view's footer buttons ship the same two
// definitions from their own IconSvg() — same ids, so assetmin collapses them
// to one symbol each.
func (t *TargetList) IconSvg() *sprite.Sprite {
	return sprite.NewSprite(
		trash.Def(),
		pencil.Def(),
	)
}
