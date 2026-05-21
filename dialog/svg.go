//go:build !wasm

package dialog

func (m *DialogWidget) IconSvg() map[string]string {
	return nil
}
