//go:build !wasm

package button

import _ "embed"

//go:embed button.css
var css string

func (b *Button) RenderCSS() string {
	return css
}

func (b *Button) IconSvg() map[string]string {
	return nil
}
