//go:build !wasm

package modal

import _ "embed"

//go:embed modal.css
var css string

func (m *Modal) RenderCSS() string {
	return css
}

func (m *Modal) IconSvg() map[string]string {
	return nil
}
