//go:build !wasm

package nav

import _ "embed"

//go:embed nav.css
var css string

func (n *Nav) RenderCSS() string {
	return css
}

func (n *Nav) IconSvg() map[string]string {
	return nil
}
