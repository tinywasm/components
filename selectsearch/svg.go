//go:build !wasm

package selectsearch

import "github.com/tinywasm/svg"

func (c *SelectSearch) IconSvg() *svg.Sprite {
	return svg.New().
		Add("ss-arrow-down", `<path fill="currentColor" d="M1.5 4.5l6.5 7 6.5-7H1.5z"/>`)
}
