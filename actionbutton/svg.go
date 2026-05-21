//go:build !wasm

package actionbutton

func (b *ActionButton) IconSvg() map[string]string {
	return nil
}
