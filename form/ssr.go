//go:build !wasm

package form

import _ "embed"

//go:embed form.css
var css string

func (f *Form) RenderCSS() string {
	return css
}
