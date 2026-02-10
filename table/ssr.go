//go:build !wasm

package table

import _ "embed"

//go:embed table.css
var css string

func (t *Table) RenderCSS() string {
	return css
}

func (t *Table) IconSvg() map[string]string {
	return nil
}
