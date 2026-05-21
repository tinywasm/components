//go:build !wasm

package themetoggle

func (t *ThemeToggle) IconSvg() map[string]string {
	return nil
}
