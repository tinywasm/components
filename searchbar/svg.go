//go:build !wasm

package searchbar

import (
	"github.com/tinywasm/svg/sprite"
)

// IconSvg registers the bar's magnifier. Method receiver (not a free function)
// so tinywasm/ssr detects a single receiver type for the package and emits
// RenderCSS + IconSvg together.
func (s *SearchBar) IconSvg() *sprite.Sprite {
	// A FILL shape (closed outline), never stroke lines: the sprite renders
	// every path with fill="currentColor" and no stroke, so a line-only path
	// (e.g. "M8 1v14") has zero area and is invisible. The lens hole renders via
	// the inner circle subpath winding opposite the outer.
	return sprite.NewSprite(
		sprite.Define(iconMagnifier, "0 0 512 512", sprite.Path("M416 208c0 45.9-14.9 88.3-40 122.7L502.6 457.4c12.5 12.5 12.5 32.8 0 45.3s-32.8 12.5-45.3 0L330.7 376c-34.4 25.2-76.8 40-122.7 40C93.1 416 0 322.9 0 208S93.1 0 208 0S416 93.1 416 208zM208 352a144 144 0 1 0 0-288 144 144 0 1 0 0 288z")),
	)
}
