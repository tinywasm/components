//go:build !wasm

package card

import _ "embed"

//go:embed card.css
var css string

func (c *Card) RenderCSS() string {
	return css
}

func (c *Card) IconSvg() map[string]string {
	return nil
}
