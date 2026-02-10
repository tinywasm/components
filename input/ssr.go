//go:build !wasm

package input

import (
	_ "embed"
)

//go:embed input.custom.css
var cssCustom string

//go:embed input.stlib.css
var cssStlib string

func (i *Input) RenderCSS() string {
	return cssStlib + "\n" + cssCustom
}

func (i *Input) IconSvg() map[string]string {
	return nil
}
