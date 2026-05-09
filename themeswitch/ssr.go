//go:build !wasm

package themeswitch

import _ "embed"

//go:embed themeswitch.css
var css string

func (t *ThemeSwitch) RenderCSS() string          { return css }
func (t *ThemeSwitch) IconSvg() map[string]string { return nil }
