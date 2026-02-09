//go:build !wasm

package content

import _ "embed"

//go:embed content.css
var css string

func (b *Content) RenderCSS() string {
	return css
}
