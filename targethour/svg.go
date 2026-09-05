//go:build !wasm

package targethour

import (
	"github.com/tinywasm/icons/pencil"
	"github.com/tinywasm/icons/selectall"
	"github.com/tinywasm/icons/trash"
	"github.com/tinywasm/svg/sprite"
)

// IconSvg registers the shared trash/pencil action glyphs (see
// targetlist/targetdate's IconSvg — same shape) plus selectall, the
// selection header's own select-all glyph (listselect.Header).
func (t *TargetHour) IconSvg() *sprite.Sprite {
	return sprite.NewSprite(trash.Def(), pencil.Def(), selectall.Def())
}
