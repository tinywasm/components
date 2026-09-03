//go:build !wasm

package targethour

import (
	"github.com/tinywasm/icons/pencil"
	"github.com/tinywasm/icons/trash"
	"github.com/tinywasm/svg/sprite"
)

func (t *TargetHour) IconSvg() *sprite.Sprite {
	return sprite.NewSprite(trash.Def(), pencil.Def())
}
